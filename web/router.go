package web

import (
	"embed"
	"net/http"

	"github.com/gorilla/mux"
	"jinya.de/goji/web/engine"
)

//go:embed tmpl
var tmplFs embed.FS

func SetupWebRouter(router *mux.Router) {
	router.
		Path("/image.php").
		HandlerFunc(imagePhp)
	router.
		Path("/jinya-content/themes/{themeName}/asset/{assetName}").
		HandlerFunc(themeAsset)
	router.
		Path("/jinya-content/themes/cache/{themeName}/{file}").
		HandlerFunc(themeCache)
	router.
		Path("/jinya-content/{hash}").
		HandlerFunc(jinyaContent)
	router.
		PathPrefix("/designer").
		Handler(installedCheck(designerPage))
	router.
		Path("/installer").
		HandlerFunc(installer)
	router.
		Path("/installer/mailer").
		HandlerFunc(installerMailer)
	router.
		Path("/installer/admin").
		HandlerFunc(installerFirstAdmin)

	router.
		Path("/{year:\\d{4}}/{month:\\d{2}}/{day:\\d{2}}/{slug:\\S+}").
		HandlerFunc(engine.ForwardBlogRequest)
	router.NotFoundHandler = http.HandlerFunc(engine.ForwardStandardRequest)
}
