package html

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// RenderPage writes layout.html plus the given page template.
func RenderPage(w http.ResponseWriter, data any, page string) {
	if err := Execute(w, "layout.html", data, "layout.html", page); err != nil {
		log.Println("render page failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// RenderPartial writes a single fragment template (HTMX responses).
func RenderPartial(w http.ResponseWriter, name string, data any) {
	if err := Execute(w, name, data, name); err != nil {
		log.Println("render partial failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func Execute(w io.Writer, execName string, data any, files ...string) error {
	tpl, err := GetTpl().Clone()
	if err != nil {
		return fmt.Errorf("template clone: %w", err)
	}
	tpl, err = tpl.ParseFS(GetViewsFs(), files...)
	if err != nil {
		return fmt.Errorf("template parse %v: %w", files, err)
	}
	if err := tpl.ExecuteTemplate(w, execName, data); err != nil {
		return fmt.Errorf("template execute %s: %w", execName, err)
	}
	return nil
}
