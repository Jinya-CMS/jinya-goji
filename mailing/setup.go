package mailing

import (
	_ "embed"
	"log/slog"

	"github.com/ebitengine/purego"
	"jinya.de/goji/native"
)

//go:generate sh -c "cd templating && cargo build --release"
//go:embed templating/target/release/libtemplating.so
var libtemplating []byte

var renderMrmlTemplate func(tmpl string) string

func SetupMailing() {
	ptr, err := native.SetupNativeBridge(libtemplating, "libtemplating")
	if err != nil {
		panic(err)
	}

	slog.Info("Register function render_mrml_template")
	purego.RegisterLibFunc(&renderMrmlTemplate, *ptr, "render_mrml_template")
}
