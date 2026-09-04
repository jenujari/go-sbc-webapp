package lib

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"jenujari/go-sbc-webapp/sqls"

	"github.com/jackc/pgx/v5/pgtype"
)

type OHLCImportResult struct {
	Inserted int64
	Skipped  int
	Errors   []string
}

type OHLCStore interface {
	UpsertOLHC(ctx context.Context, arg sqls.UpsertOLHCParams) (int64, error)
}

type OHLCImporter struct {
	Store OHLCStore
}

func NewOHLCImporter(db *DBService) *OHLCImporter {
	if db == nil || db.Queries == nil {
		return nil
	}
	return &OHLCImporter{Store: db.Queries}
}

func (s *OHLCImporter) Import(ctx context.Context, r io.Reader, tickerID int16) (OHLCImportResult, error) {
	if s == nil || s.Store == nil {
		return OHLCImportResult{}, fmt.Errorf("ohlc importer is not configured")
	}

	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	header, err := reader.Read()
	if err != nil {
		return OHLCImportResult{}, fmt.Errorf("unable to read CSV header: %w", err)
	}
	if err := validateOHLCHeader(header); err != nil {
		return OHLCImportResult{}, err
	}

	var result OHLCImportResult
	line := 1
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		line++
		if err != nil {
			result.Skipped++
			result.Errors = appendLimited(result.Errors, fmt.Sprintf("line %d: %v", line, err))
			continue
		}
		params, err := parseOHLCRecord(record, tickerID)
		if err != nil {
			result.Skipped++
			result.Errors = appendLimited(result.Errors, fmt.Sprintf("line %d: %v", line, err))
			continue
		}
		n, err := s.Store.UpsertOLHC(ctx, params)
		if err != nil {
			result.Skipped++
			result.Errors = appendLimited(result.Errors, fmt.Sprintf("line %d: insert failed: %v", line, err))
			continue
		}
		result.Inserted += n
	}
	return result, nil
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
