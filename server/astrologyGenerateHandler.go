package server

import (
	"net/http"

	"jenujari/go-sbc-webapp/html"
	"jenujari/go-sbc-webapp/lib"
)

type astrologyGenerateResultData struct {
	PanchangUpserted        int64
	PlanetPositionsUpserted int64
	Errors                  []string
}

func astrologyGenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	app, ok := requestApp(r)
	if !ok {
		http.Error(w, "Services not available", http.StatusInternalServerError)
		return
	}
	if app.DB == nil {
		http.Error(w, "Database connection is not available. Check db.url or DATABASE_URL.", http.StatusInternalServerError)
		return
	}
	if app.SweClient == nil {
		http.Error(w, "Swe gRPC Client is not available.", http.StatusInternalServerError)
		return
	}
	if app.Astrology == nil {
		http.Error(w, "Services not available", http.StatusInternalServerError)
		return
	}

	req, err := lib.ParseAstrologyRange(r.FormValue("from_date"), r.FormValue("to_date"), r.FormValue("time"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := app.Astrology.Generate(r.Context(), req)
	renderAstrologyGenerateResult(w, astrologyGenerateResultData{
		PanchangUpserted:        result.PanchangUpserted,
		PlanetPositionsUpserted: result.PlanetPositionsUpserted,
		Errors:                  result.Errors,
	})
}

func renderAstrologyGenerateResult(w http.ResponseWriter, data astrologyGenerateResultData) {
	w.WriteHeader(http.StatusOK)
	html.RenderPartial(w, "astrology_generate_result.html", data)
}
