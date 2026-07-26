package engine

import (
	_ "embed"
	"log/slog"
	"runtime"

	"github.com/ebitengine/purego"
	"jinya.de/goji/native"
)

//go:generate sh -c "cargo build --release"
//go:embed target/release/libengine.so
var libengine []byte

var (
	nativeHandleRequest func(apiRequest bool, url string, output *uintptr) int
	nativeHandleError   func(apiRequest bool, url string, statusCode uint16, output *uintptr) int
	freeEngineString    func(ptr uintptr)
)

func SetupEngine() {
	ptr, err := native.SetupNativeBridge(libengine, "libengine")
	if err != nil {
		panic(err)
	}

	slog.Info("Register function handle_request")
	purego.RegisterLibFunc(&nativeHandleRequest, *ptr, "handle_request")
	purego.RegisterLibFunc(&nativeHandleError, *ptr, "handle_error")
	purego.RegisterLibFunc(&freeEngineString, *ptr, "free_engine_string")
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
