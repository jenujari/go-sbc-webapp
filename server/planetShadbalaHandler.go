package server

import (
	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/html"
	"jenujari/go-sbc-webapp/lib"
	"net/http"
	"time"
)

func planetShadbalaHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	services := ctx.Value("services").(map[string]any)
	webData := cloneWebData(services["webData"].(lib.WebData))
	webData["currentTime"] = time.Now().Format("2006-01-02T15:04")

	tpl, err := html.GetTpl().Clone()
	if err != nil {
		config.GetLogger().Println("template clone failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tpl, err = tpl.ParseFS(html.GetViewsFs(), "layout.html", "planet_shadbala.html")
	if err != nil {
		config.GetLogger().Println("template not found", err)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "layout.html", webData); err != nil {
		config.GetLogger().Println("template execution failed", err)
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}

func planetShadbalaResultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	services := ctx.Value("services").(map[string]any)
	webData := cloneWebData(services["webData"].(lib.WebData))
	shadbalaService := services["planetShadbalaService"].(lib.PlanetShadbalaService)

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

	tpl, err := html.GetTpl().Clone()
	if err != nil {
		config.GetLogger().Println("template clone failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tpl, err = tpl.ParseFS(html.GetViewsFs(), "planet_shadbala_result.html")
	if err != nil {
		config.GetLogger().Println("template not found", err)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "planet_shadbala_result.html", webData); err != nil {
		config.GetLogger().Println("template execution failed", err)
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}
