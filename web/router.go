package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"jinya.de/goji/web/engine"
)

func SetupWebRouter(router *mux.Router) {
	router.
		Path("/image.php").
		HandlerFunc(imagePhp)
	router.
		Path("/jinya-content/themes/{themeName}/asset/{assetName}").
		HandlerFunc(themeAsset)
	router.
		Path("/jinya-content/{hash}").
		HandlerFunc(jinyaContent)
	router.
		PathPrefix("/designer").
		HandlerFunc(designerPage)

	router.
		Path("/{year:\\d{4}}/{month:\\d{2}}/{day:\\d{2}}/{slug:\\S+}").
		HandlerFunc(engine.ForwardBlogRequest)
	router.NotFoundHandler = http.HandlerFunc(engine.ForwardStandardRequest)
}
