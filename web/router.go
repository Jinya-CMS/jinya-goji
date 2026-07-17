package web

import (
	"net/http"

	"github.com/gorilla/mux"
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
	router.NotFoundHandler = http.HandlerFunc(designerPage)
}
