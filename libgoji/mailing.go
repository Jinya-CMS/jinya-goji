package libgoji

import (
	_ "embed"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/ebitengine/purego"
)

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
)

func bindMailing(ptr uintptr) {
	slog.Info("Register function render_mrml_template")
	purego.RegisterLibFunc(&nativeRenderMrmlTemplate, ptr, "mailing_render_mrml_template")
}

func RenderMrmlTemplate(tmpl string) (string, error) {
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

	defer freeString(output)

	html := cStringToGoString(output)

	return html, nil
}
