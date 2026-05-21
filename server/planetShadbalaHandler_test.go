package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"jenujari/go-sbc-webapp/html"
	"jenujari/go-sbc-webapp/lib"
)

func TestPlanetShadbalaResultsHandlerHTMXFragment(t *testing.T) {
	form := url.Values{"datetime": {"2026-05-21T12:30"}}
	req := httptest.NewRequest(http.MethodPost, "/planet-shadbala/results", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	req = req.WithContext(context.WithValue(req.Context(), "services", map[string]any{
		"webData": lib.WebData{"appname": "webapp"},
		"planetShadbalaService": fakePlanetShadbalaService{view: lib.PlanetShadbalaView{Planets: []lib.PlanetShadbalaRecord{{
			ID:           "sun",
			Name:         "Sun",
			Symbol:       "☉",
			UdayBala:     90,
			UchchaBala:   80,
			VakraBala:    0,
			NavamshaBala: 70,
			KshetraBala:  60,
			Total:        60,
			Sign:         "Leo",
			Nakshatra:    "Magha",
			Speed:        "0.9800° / day",
			State:        "Direct",
			Retrograde:   false,
		}}},
		},
	}))

	rr := httptest.NewRecorder()
	planetShadbalaResultsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Detailed Strength Analysis") || !strings.Contains(body, "Sun") || !strings.Contains(body, "Total Power") {
		t.Fatalf("expected shadbala fragment, got: %s", body)
	}
}

func TestPlanetShadbalaResultTemplateNilSafe(t *testing.T) {
	tpl, err := html.GetTpl().Clone()
	if err != nil {
		t.Fatal(err)
	}
	tpl, err = tpl.ParseFS(html.GetViewsFs(), "planet_shadbala_result.html")
	if err != nil {
		t.Fatal(err)
	}

	var out strings.Builder
	err = tpl.ExecuteTemplate(&out, "planet_shadbala_result.html", lib.WebData{
		"displayDate": "May 21, 2026 12:30",
		"shadbala": lib.PlanetShadbalaView{Planets: []lib.PlanetShadbalaRecord{{
			ID:    "moon",
			Name:  "Moon",
			State: "Direct",
		}}},
	})
	if err != nil {
		t.Fatalf("expected nil-safe template execution, got %v", err)
	}
	if !strings.Contains(out.String(), "Moon") {
		t.Fatalf("expected rendered planet, got %s", out.String())
	}
}

type fakePlanetShadbalaService struct {
	view lib.PlanetShadbalaView
	err  error
}

func (f fakePlanetShadbalaService) GetPlanetShadbala(ctx context.Context, timestamp string) (lib.PlanetShadbalaView, error) {
	return f.view, f.err
}
