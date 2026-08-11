package engine

import (
	_ "embed"
	"encoding/json"
	"errors"
	"log/slog"
	"runtime"

	"github.com/ebitengine/purego"
	"jinya.de/goji/database"
	"jinya.de/goji/native"
)

//go:generate sh -c "cargo build --release"
//go:embed target/release/libengine.so
var libengine []byte

var (
	nativeHandleRequest func(apiRequest bool, url string, output *uintptr) int
	nativeHandleError   func(apiRequest bool, url string, statusCode uint16, output *uintptr) int

	nativePreviewClassicPage func(themeId int, classicPageId int, classicPageData string, output *uintptr) int
	nativePreviewModernPage  func(themeId int, modernPageId int, modernPageSectionsData string, output *uintptr) int
	nativePreviewBlogPost    func(themeId int, blogPostId int, blogPostSectionsData string, output *uintptr) int

	freeEngineString func(ptr uintptr)
)

func SetupEngine() {
	ptr, err := native.SetupNativeBridge(libengine, "libengine")
	if err != nil {
		panic(err)
	}

	slog.Info("Register function handle_request")
	purego.RegisterLibFunc(&nativeHandleRequest, *ptr, "handle_request")
	purego.RegisterLibFunc(&nativeHandleError, *ptr, "handle_error")

	purego.RegisterLibFunc(&nativePreviewClassicPage, *ptr, "preview_classic_page")
	purego.RegisterLibFunc(&nativePreviewModernPage, *ptr, "preview_modern_page")
	purego.RegisterLibFunc(&nativePreviewBlogPost, *ptr, "preview_blog_post")

	purego.RegisterLibFunc(&freeEngineString, *ptr, "free_engine_string")
}

var (
	PreviewErrorFailedToRender = errors.New("failed to render template")
	PreviewErrorUnknownFailure = errors.New("failed to preview")
)

func PreviewClassicPage(themeId int, classicPage *database.ClassicPage) (string, error) {
	var output uintptr

	data := map[string]any{
		"id":            classicPage.Id,
		"creatorId":     classicPage.CreatorId,
		"updatedById":   classicPage.UpdatedById,
		"createdAt":     classicPage.CreatedAt,
		"lastUpdatedAt": classicPage.LastUpdatedAt,
		"content":       classicPage.Content,
		"title":         classicPage.Title,
	}

	pageData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	code := nativePreviewClassicPage(themeId, classicPage.Id, string(pageData), &output)
	runtime.KeepAlive(output)

	defer freeEngineString(output)

	if code == 1 {
		return "", PreviewErrorUnknownFailure
	}
	if code == -1 {
		return "", PreviewErrorFailedToRender
	}

	html := native.CStringToGoString(output)

	return html, nil
}

func PreviewModernPage(themeId, modernPageId int, modernPageSectionsJson string) (string, error) {
	var output uintptr

	code := nativePreviewModernPage(themeId, modernPageId, modernPageSectionsJson, &output)
	runtime.KeepAlive(output)

	defer freeEngineString(output)

	if code == 1 {
		return "", PreviewErrorUnknownFailure
	}
	if code == -1 {
		return "", PreviewErrorFailedToRender
	}

	html := native.CStringToGoString(output)

	return html, nil
}

func PreviewBlogPost(themeId, blogPostId int, blogPostSectionsJson string) (string, error) {
	var output uintptr

	code := nativePreviewBlogPost(themeId, blogPostId, blogPostSectionsJson, &output)
	runtime.KeepAlive(output)

	defer freeEngineString(output)

	if code == 1 {
		return "", PreviewErrorUnknownFailure
	}
	if code == -1 {
		return "", PreviewErrorFailedToRender
	}

	html := native.CStringToGoString(output)

	return html, nil
}

func handleRequest(apiRequest bool, url string) ([]byte, int) {
	var output uintptr

	code := nativeHandleRequest(apiRequest, url, &output)
	runtime.KeepAlive(output)

	defer freeEngineString(output)

	html := native.CStringToGoString(output)

	return []byte(html), code
}

func handleError(apiRequest bool, url string, statusCode int) ([]byte, int) {
	var output uintptr

	code := nativeHandleError(apiRequest, url, uint16(statusCode), &output)
	runtime.KeepAlive(output)

	defer freeEngineString(output)

	html := native.CStringToGoString(output)

	return []byte(html), code
}
