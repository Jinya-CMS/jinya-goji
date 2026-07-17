package mailing

import (
	_ "embed"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/ebitengine/purego"
	"jinya.de/goji/native"
)

//go:generate sh -c "cd templating && cargo build --release"
//go:embed templating/target/release/libtemplating.so
var libtemplating []byte

const (
	errCodeMrmlInvalidData = 1
	errCodeMrmlParse       = 2
	errCodeMrmlRender      = 3
)

var (
	ErrMrmlInvalidData = fmt.Errorf("invalid data for the rendering call")
	ErrMrmlParse       = fmt.Errorf("parsing the mrml template failed")
	ErrMrmlRender      = fmt.Errorf("rendering the mrml tmplate failed")
)

var (
	nativeRenderMrmlTemplate func(tmpl string, output *uintptr) int
	freeMrmlString           func(ptr uintptr)
)

func SetupMailing() {
	ptr, err := native.SetupNativeBridge(libtemplating, "libtemplating")
	if err != nil {
		panic(err)
	}

	slog.Info("Register function render_mrml_template")
	purego.RegisterLibFunc(&nativeRenderMrmlTemplate, *ptr, "render_mrml_template")
	purego.RegisterLibFunc(&freeMrmlString, *ptr, "free_mrml_string")
}

func renderMrmlTemplate(tmpl string) (string, error) {
	var output uintptr

	code := nativeRenderMrmlTemplate(tmpl, &output)

	runtime.KeepAlive(output)

	switch code {
	case errCodeMrmlRender:
		return "", ErrMrmlRender
	case errCodeMrmlParse:
		return "", ErrMrmlParse
	case errCodeMrmlInvalidData:
		return "", ErrMrmlInvalidData
	}

	defer freeMrmlString(output)

	html := native.CStringToGoString(output)

	return html, nil
}
