package server

import (
	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/html"
	"net/http"
	"time"
)

func planetShadbalaHandler(w http.ResponseWriter, r *http.Request) {
	app, ok := requestApp(r)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	webData := app.PageData()
	webData["currentTime"] = time.Now().Format("2006-01-02T15:04")

	html.RenderPage(w, webData, "planet_shadbala.html")
}

func planetShadbalaResultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	app, ok := requestApp(r)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	webData := app.PageData()
	shadbalaService := app.PlanetShadbala

	datetime := r.FormValue("datetime")
	if datetime == "" {
		http.Error(w, "datetime is required", http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse("2006-01-02T15:04", datetime)
	if err != nil {
		http.Error(w, "invalid datetime", http.StatusBadRequest)
		return
	}

	view, err := shadbalaService.GetPlanetShadbala(ctx, parsedDate.Format(time.RFC3339))
	if err != nil {
		config.GetLogger().Println("get planet shadbala failed", err)
		webData["shadbalaError"] = "Unable to fetch planetary strength details right now. Please try again."
	} else {
		webData["shadbala"] = view
	}
	webData["displayDate"] = parsedDate.Format("January 02, 2006 15:04")

	html.RenderPartial(w, "planet_shadbala_result.html", webData)
}
