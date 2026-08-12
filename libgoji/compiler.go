package libgoji

import (
	"fmt"
	"log/slog"

	"github.com/ebitengine/purego"
)

var (
	ErrTokio    = fmt.Errorf("failed to start tokio")
	ErrCompiler = fmt.Errorf("failed to compile theme")
	ErrUnknown  = fmt.Errorf("unknown error")
)

var compileTheme func(themeId int) int

func bindCompiler(ptr uintptr) {
	slog.Info("Register function compileTheme")
	purego.RegisterLibFunc(&compileTheme, ptr, "compiler_compile_theme")
}

func CompileTheme(themeId int) error {
	result := compileTheme(themeId)
	switch result {
	case 0:
		return nil
	case 1:
		return ErrCompiler
	case 2:
		return ErrTokio
	default:
		return ErrUnknown
	}
}
