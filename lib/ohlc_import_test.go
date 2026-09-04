package lib

import (
	"context"
	"strings"
	"testing"

	"jenujari/go-sbc-webapp/sqls"
)

type stubOHLCStore struct {
	n     int64
	err   error
	calls int
}

func (s *stubOHLCStore) UpsertOLHC(ctx context.Context, arg sqls.UpsertOLHCParams) (int64, error) {
	s.calls++
	if s.err != nil {
		return 0, s.err
	}
	return s.n, nil
}

func TestOHLCImporterImportHappyPath(t *testing.T) {
	store := &stubOHLCStore{n: 1}
	svc := &OHLCImporter{Store: store}
	csv := "Date,Open,High,Low,Close,Volume\n2026-05-21,1,2,0.5,1.5,100\n"

	result, err := svc.Import(context.Background(), strings.NewReader(csv), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Inserted != 1 || result.Skipped != 0 || store.calls != 1 {
		t.Fatalf("unexpected result: %+v calls=%d", result, store.calls)
	}
}

func TestOHLCImporterRejectsBadHeader(t *testing.T) {
	svc := &OHLCImporter{Store: &stubOHLCStore{n: 1}}
	_, err := svc.Import(context.Background(), strings.NewReader("foo,bar\n"), 7)
	if err == nil {
		t.Fatal("expected header error")
	}
}

func TestOHLCImporterSkipsBadRows(t *testing.T) {
	store := &stubOHLCStore{n: 1}
	svc := &OHLCImporter{Store: store}
	csv := "Date,Open,High,Low,Close,Volume\nbad-date,1,2,3,4,5\n2026-05-21,1,2,3,4,5\n"

	result, err := svc.Import(context.Background(), strings.NewReader(csv), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Inserted != 1 || result.Skipped != 1 || len(result.Errors) != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestNewOHLCImporterNilDB(t *testing.T) {
	if NewOHLCImporter(nil) != nil {
		t.Fatal("expected nil importer")
	}
}
