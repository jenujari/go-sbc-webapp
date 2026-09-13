package html

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderPageUsesLayoutAndPage(t *testing.T) {
	rr := httptest.NewRecorder()
	RenderPage(rr, map[string]any{"appname": "webapp"}, "index.html")

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if !strings.Contains(body, "webapp") || !strings.Contains(body, "Welcome back") {
		t.Fatalf("expected layout+index content, got: %s", body)
	}
}

func TestExecuteUnknownTemplate(t *testing.T) {
	var out strings.Builder
	if err := Execute(&out, "missing.html", nil, "missing.html"); err == nil {
		t.Fatal("expected parse error for missing template")
	}
}
