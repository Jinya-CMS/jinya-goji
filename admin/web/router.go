package web

import (
	"net/http"
	"strings"

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
	router.
		MatcherFunc(func(r *http.Request, match *mux.RouteMatch) bool {
			return !strings.HasPrefix(r.URL.Path, "/api") && !strings.HasPrefix(r.URL.Path, "/static")
		}).
		HandlerFunc(indexPage)
}
