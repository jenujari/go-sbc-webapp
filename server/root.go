package server

import (
	"errors"
	"fmt"
	"net/http"

	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/html"

	rtc "github.com/jenujari/runtime-context"
)

var (
	server *http.Server
	router *http.ServeMux
)

func init() {

	server = &http.Server{
		Addr:              ":3000",
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		MaxHeaderBytes:    0,
	}

	router = http.NewServeMux()

	router.Handle("/static/", staticHander())

	router.HandleFunc("/", indexhandler)

	server.Handler = router
	config.GetLogger().Println("server initialization complete.")
}

func RunServer() {
	pc := rtc.GetMainProcess()

	go func(cmdx *rtc.ProcessContext) {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			cmdx.FatalErrorChan <- fmt.Errorf("ListenAndServe(): %v", err)
		}
	}(pc)

	<-pc.CTX.Done()
	config.GetLogger().Println("shutting down server...")
	if err := server.Shutdown(pc.CTX); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	config.GetLogger().Println("server shutdown complete...")
}

func GetServer() *http.Server {
	return server
}

func staticHander() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.FS(html.GetAssetsFs())))
}

func indexhandler(w http.ResponseWriter, r *http.Request) {

	tpl, err := html.GetTpl().Clone()
	if err != nil {
		config.GetLogger().Println("template clone failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tpl, err = tpl.ParseFS(html.GetViewsFs(), "layout.html", "index.html")

	// fmt.Println(tpl.DefinedTemplates())

	if err != nil {
		config.GetLogger().Println("template not found", err)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "layout.html", nil)
	if err != nil {
		config.GetLogger().Println("template execution failed", err)
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}

}
