package mailing

import (
	"embed"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ebitengine/purego"
)

//go:generate sh -c "cd templating && cargo build --release"
//go:embed templating/target/release/libtemplating.so
var native embed.FS

var (
	renderMrmlTemplate func(tmpl string) string
)

func registerFunctions() {
	slog.Info("Loading libtemplating.so from embed")
	data, err := native.ReadFile("templating/target/release/libtemplating.so")
	if err != nil {
		slog.Error("Failed to read binary from embed", "error", err)
		panic(err)
	}

	tmpPath := filepath.Join(os.TempDir(), "libtemplating.so")
	slog.Info("Write libtemplate.so to temp directory", "path", tmpPath)
	if err := os.WriteFile(tmpPath, data, 0o666); err != nil {
		slog.Error("Failed to write binary to temp directory", "error", err, "tmp-path", tmpPath)
		panic(err)
	}

	slog.Info("Perform Dlopen for libtemplating.so", "path", tmpPath)
	libtemplating, err := purego.Dlopen(tmpPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		slog.Error("Failed to open the templating library", "error", err)
		panic(err)
	}

	slog.Info("Register function render_mrml_template")
	purego.RegisterLibFunc(&renderMrmlTemplate, libtemplating, "render_mrml_template")
}
