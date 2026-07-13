package web

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed tmpl
var tmplFs embed.FS

func indexPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("content").ParseFS(tmplFs, "tmpl/index.gohtml")
	if err == nil {
		t.Execute(w, nil)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
