package compiler

import (
	_ "embed"
	"log/slog"

	"github.com/ebitengine/purego"
	"jinya.de/goji/native"
)

//go:generate sh -c "cargo build --release"
//go:embed target/release/libcompiler.so
var libcompiler []byte

var compileTheme func(themeId int) int

func SetupCompiler() {
	ptr, err := native.SetupNativeBridge(libcompiler, "libcompiler")
	if err != nil {
		panic(err)
	}

	slog.Info("Register function compile_theme")
	purego.RegisterLibFunc(&compileTheme, *ptr, "compile_theme")
}
