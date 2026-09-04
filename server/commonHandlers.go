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

	tpl, err := html.GetTpl().Clone()
	if err != nil {
		config.GetLogger().Println("template clone failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tpl, err = tpl.ParseFS(html.GetViewsFs(), "layout.html", "index.html")

	if err != nil {
		config.GetLogger().Println("template not found", err)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "layout.html", webData)
	if err != nil {
		config.GetLogger().Println("template execution failed", err)
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}

}
