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
	freeEngineString    func(ptr uintptr)
)

func SetupEngine() {
	ptr, err := native.SetupNativeBridge(libengine, "libengine")
	if err != nil {
		panic(err)
	}

	slog.Info("Register function handle_request")
	purego.RegisterLibFunc(&nativeHandleRequest, *ptr, "handle_request")
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
