package compiler

import "os"

func CompileTheme(themeId int) error {
	if os.Getenv("DEBUG_THEME_COMPILER") == "true" {
		return debugCompileTheme(themeId)
	}

	return releaseCompileTheme(themeId)
}
