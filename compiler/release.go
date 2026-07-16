package compiler

import "fmt"

var (
	ErrTokio    = fmt.Errorf("failed to start tokio")
	ErrCompiler = fmt.Errorf("failed to compile theme")
	ErrUnknown  = fmt.Errorf("unknown error")
)

func releaseCompileTheme(themeId int) error {
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
