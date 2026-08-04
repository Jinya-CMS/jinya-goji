package web

import (
	"html/template"
	"net/http"
)

func designerPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("content").ParseFS(tmplFs, "tmpl/designer.gohtml")
	if err == nil {
		t.Execute(w, nil)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
