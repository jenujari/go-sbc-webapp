package server

import (
	"net/http"
	"strconv"

	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/html"
	"jenujari/go-sbc-webapp/lib"
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
	app, ok := requestApp(r)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	data := app.PageData()
	db := app.DB
	if db == nil {
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

	html.RenderPage(w, data, "ohlc_upload.html")
}

func ohlcUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	app, ok := requestApp(r)
	if !ok || app.OHLC == nil {
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

	file, _, err := r.FormFile("ohlc_file")
	if err != nil {
		renderOHLCUploadResult(w, ohlcUploadResultData{Errors: []string{"Please select a CSV file."}}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	result, err := app.OHLC.Import(r.Context(), file, int16(tickerID64))
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		renderOHLCUploadResult(w, ohlcResult(result), http.StatusInternalServerError)
		return
	}
	renderOHLCUploadResult(w, ohlcResult(result), http.StatusOK)
}

func ohlcResult(result lib.OHLCImportResult) ohlcUploadResultData {
	return ohlcUploadResultData{Inserted: result.Inserted, Skipped: result.Skipped, Errors: result.Errors}
}

func renderOHLCUploadResult(w http.ResponseWriter, data ohlcUploadResultData, status int) {
	_ = status
	w.WriteHeader(http.StatusOK)
	html.RenderPartial(w, "ohlc_upload_result.html", data)
}

func appendLimited(items []string, value string) []string {
	if len(items) >= 20 {
		return items
	}
	return append(items, value)
}
