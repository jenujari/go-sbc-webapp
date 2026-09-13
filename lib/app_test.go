package lib

import (
	"context"
	"testing"
)

func TestAppContextRoundTrip(t *testing.T) {
	want := &App{WebData: WebData{"appname": "webapp"}}

	got, ok := AppFromContext(WithApp(context.Background(), want))
	if !ok || got != want {
		t.Fatalf("expected same *App, ok=%v got=%v", ok, got)
	}
}

func TestAppFromContextMissing(t *testing.T) {
	if _, ok := AppFromContext(context.Background()); ok {
		t.Fatal("expected missing app")
	}
	if _, ok := AppFromContext(WithApp(context.Background(), nil)); ok {
		t.Fatal("expected nil app to be treated as missing")
	}
}
