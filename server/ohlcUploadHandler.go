package server

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/html"
	"jenujari/go-sbc-webapp/lib"
	"jenujari/go-sbc-webapp/sqls"

	"github.com/jackc/pgx/v5/pgtype"
)

type ohlcUploadResultData struct {
	Inserted int64
	Skipped  int
	Errors   []string
}

func ohlcUploadPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ohlc-upload" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	services := ctx.Value("services").(map[string]any)
	webData := services["webData"].(lib.WebData)

	data := cloneWebData(webData)
	db, ok := services["db"].(*lib.DBService)
	if !ok || db == nil {
		data["Error"] = "Database connection is not available. Check db.url or DATABASE_URL."
	} else {
		tickers, err := db.Queries.ListTickers(ctx)
		if err != nil {
			config.GetLogger().Println("list tickers failed", err)
			data["Error"] = "Unable to load ticker list from database."
		} else {
			data["Tickers"] = tickers
		}
	}

	tpl, err := html.GetTpl().Clone()
	if err != nil {
		config.GetLogger().Println("template clone failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	tpl, err = tpl.ParseFS(html.GetViewsFs(), "layout.html", "ohlc_upload.html")
	if err != nil {
		config.GetLogger().Println("template not found", err)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	if err := tpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		config.GetLogger().Println("template execution failed", err)
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}

func cloneWebData(webData lib.WebData) lib.WebData {
	data := make(lib.WebData, len(webData)+2)
	for key, value := range webData {
		data[key] = value
	}
	return data
}

func ohlcUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	services := ctx.Value("services").(map[string]any)
	db, ok := services["db"].(*lib.DBService)
	if !ok || db == nil {
		renderOHLCUploadResult(w, ohlcUploadResultData{Errors: []string{"Database connection is not available."}}, http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(128 << 20); err != nil {
		renderOHLCUploadResult(w, ohlcUploadResultData{Errors: []string{"Invalid multipart form: " + err.Error()}}, http.StatusBadRequest)
		return
	}

	tickerID64, err := strconv.ParseInt(r.FormValue("ticker_id"), 10, 16)
	if err != nil || tickerID64 <= 0 {
		renderOHLCUploadResult(w, ohlcUploadResultData{Errors: []string{"Please select a valid ticker."}}, http.StatusBadRequest)
		return
	}
	tickerID := int16(tickerID64)

	file, header, err := r.FormFile("ohlc_file")
	if err != nil {
		renderOHLCUploadResult(w, ohlcUploadResultData{Errors: []string{"Please select a CSV file."}}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmpDir := filepath.Join(os.TempDir(), "go-sbc-webapp-ohlc")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		renderOHLCUploadResult(w, ohlcUploadResultData{Errors: []string{"Unable to create temp folder: " + err.Error()}}, http.StatusInternalServerError)
		return
	}
	tmpPath := filepath.Join(tmpDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename)))
	out, err := os.Create(tmpPath)
	if err != nil {
		renderOHLCUploadResult(w, ohlcUploadResultData{Errors: []string{"Unable to store upload: " + err.Error()}}, http.StatusInternalServerError)
		return
	}
	_, copyErr := io.Copy(out, file)
	closeErr := out.Close()
	if copyErr != nil || closeErr != nil {
		renderOHLCUploadResult(w, ohlcUploadResultData{Errors: []string{"Unable to save uploaded file."}}, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpPath)

	inserted, skipped, rowErrors, err := processOHLCFile(ctx, db, tmpPath, tickerID)
	if err != nil {
		rowErrors = append(rowErrors, err.Error())
		renderOHLCUploadResult(w, ohlcUploadResultData{Inserted: inserted, Skipped: skipped, Errors: rowErrors}, http.StatusInternalServerError)
		return
	}
	renderOHLCUploadResult(w, ohlcUploadResultData{Inserted: inserted, Skipped: skipped, Errors: rowErrors}, http.StatusOK)
}

func processOHLCFile(ctx context.Context, db *lib.DBService, tmpPath string, tickerID int16) (int64, int, []string, error) {
	f, err := os.Open(tmpPath)
	if err != nil {
		return 0, 0, nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true
	header, err := reader.Read()
	if err != nil {
		return 0, 0, nil, fmt.Errorf("unable to read CSV header: %w", err)
	}
	if err := validateOHLCHeader(header); err != nil {
		return 0, 0, nil, err
	}

	var inserted int64
	var skipped int
	var rowErrors []string
	line := 1
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		line++
		if err != nil {
			skipped++
			rowErrors = appendLimited(rowErrors, fmt.Sprintf("line %d: %v", line, err))
			continue
		}
		params, err := parseOHLCRecord(record, tickerID)
		if err != nil {
			skipped++
			rowErrors = appendLimited(rowErrors, fmt.Sprintf("line %d: %v", line, err))
			continue
		}
		n, err := db.Queries.UpsertOLHC(ctx, params)
		if err != nil {
			skipped++
			rowErrors = appendLimited(rowErrors, fmt.Sprintf("line %d: insert failed: %v", line, err))
			continue
		}
		inserted += n
	}
	return inserted, skipped, rowErrors, nil
}

func validateOHLCHeader(header []string) error {
	expected := []string{"date", "open", "high", "low", "close", "volume"}
	if len(header) < len(expected) {
		return fmt.Errorf("CSV must have columns: Date, Open, High, Low, Close, Volume")
	}
	for i, want := range expected {
		if strings.ToLower(strings.TrimSpace(header[i])) != want {
			return fmt.Errorf("unexpected CSV column %d %q; expected %q", i+1, header[i], want)
		}
	}
	return nil
}

func parseOHLCRecord(record []string, tickerID int16) (sqls.UpsertOLHCParams, error) {
	if len(record) < 6 {
		return sqls.UpsertOLHCParams{}, fmt.Errorf("expected 6 columns, got %d", len(record))
	}
	day, err := parseCSVTime(record[0])
	if err != nil {
		return sqls.UpsertOLHCParams{}, err
	}
	o, err := parseFloatPtr(record[1], "Open")
	if err != nil {
		return sqls.UpsertOLHCParams{}, err
	}
	h, err := parseFloatPtr(record[2], "High")
	if err != nil {
		return sqls.UpsertOLHCParams{}, err
	}
	l, err := parseFloatPtr(record[3], "Low")
	if err != nil {
		return sqls.UpsertOLHCParams{}, err
	}
	cl, err := parseFloatPtr(record[4], "Close")
	if err != nil {
		return sqls.UpsertOLHCParams{}, err
	}
	v, err := parseFloatPtr(record[5], "Volume")
	if err != nil {
		return sqls.UpsertOLHCParams{}, err
	}

	return sqls.UpsertOLHCParams{
		Day:      pgtype.Timestamptz{Time: day, Valid: true},
		TickerID: tickerID,
		O:        o,
		H:        h,
		L:        l,
		C:        cl,
		V:        v,
	}, nil
}

func parseCSVTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	layouts := []string{"1/2/2006 15:04:05", "1/2/2006 15:04", "2006-01-02 15:04:05", "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid Date %q", value)
}

func parseFloatPtr(value, label string) (*float64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid %s %q", label, value)
	}
	return &f, nil
}

func appendLimited(items []string, value string) []string {
	if len(items) >= 20 {
		return items
	}
	return append(items, value)
}

func renderOHLCUploadResult(w http.ResponseWriter, data ohlcUploadResultData, status int) {
	_ = status
	w.WriteHeader(http.StatusOK)
	tpl, err := html.GetTpl().Clone()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	tpl, err = tpl.ParseFS(html.GetViewsFs(), "ohlc_upload_result.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	if err := tpl.ExecuteTemplate(w, "ohlc_upload_result.html", data); err != nil {
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}
