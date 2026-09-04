package server

import (
	"errors"
	"fmt"
	"net/http"

	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/lib"

	rtc "github.com/jenujari/runtime-context"
)

var (
	server *http.Server
	router *http.ServeMux
)

func init() {
	cfg := config.GetConfig()

	server = &http.Server{
		Addr:              ":" + cfg.WebAppConfig.Port,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		MaxHeaderBytes:    0,
	}

	router = http.NewServeMux()

	router.Handle("/static/", staticHander())

	router.HandleFunc("/pos-table", planetPosHandler)
	router.HandleFunc("/positions", positionsHandler)
	router.HandleFunc("/tithy-table", tithyHandler)
	router.HandleFunc("/tithies", tithyTableHandler)
	router.HandleFunc("/planet-conjunction", conjunctionHandler)
	router.HandleFunc("/planet-conjunctions", conjunctionSearchHandler)
	router.HandleFunc("/planet-shadbala", planetShadbalaHandler)
	router.HandleFunc("/planet-shadbala/results", planetShadbalaResultsHandler)
	router.HandleFunc("/ohlc-upload", ohlcUploadPageHandler)
	router.HandleFunc("/ohlc-upload/import", ohlcUploadHandler)
	router.HandleFunc("/ohlc-upload/generate-astrology", astrologyGenerateHandler)
	router.HandleFunc("/", indexhandler)

	server.Handler = withApp(router)
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

func withApp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(lib.WithApp(r.Context(), lib.GetApp())))
	})
}

func requestApp(r *http.Request) (*lib.App, bool) {
	return lib.AppFromContext(r.Context())
}
