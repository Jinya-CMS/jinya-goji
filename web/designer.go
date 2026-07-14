package web

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed tmpl
var tmplFs embed.FS

func designerPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("content").ParseFS(tmplFs, "tmpl/designer.gohtml")
	if err == nil {
		t.Execute(w, nil)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
