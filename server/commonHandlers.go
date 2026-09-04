package server

import (
	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/html"
	"net/http"
)

func staticHander() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.FS(html.GetAssetsFs())))
}

func indexhandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	app, ok := requestApp(r)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	webData := app.WebData
	sweClient := app.SweClient

	pingResp, err := sweClient.Ping(ctx)
	if err != nil {
		config.GetLogger().Println("ping failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	config.GetLogger().Println("ping response", pingResp)

	html.RenderPage(w, webData, "index.html")
}
