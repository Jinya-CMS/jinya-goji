package theming

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
	"jinya.de/goji/database"
)

func GetThemeBasePath() string {
	themePath := os.Getenv("THEME_PATH")
	if themePath != "" {
		return themePath
	}

	wd, err := os.Getwd()
	if err != nil {
		return "themes"
	}

	return path.Join(wd, themePath)
}

func GetThemePath(name string) string {
	return path.Join(GetThemeBasePath(), name)
}

type themeConfig struct {
	DisplayName string `yaml:"displayName"`
	Description any    `yaml:"description"`
	Type        string `yaml:"type"`
}

func SyncThemes() error {
	slog.Debug("Sync all themes")
	allThemes, err := database.FindAllThemes()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(GetThemeBasePath())
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("Theme base path does not exist", "path", GetThemeBasePath())
			return nil
		}
		return err
	}

	var themes []string
	for _, entry := range entries {
		if entry.IsDir() {
			themePath := filepath.Join(GetThemeBasePath(), entry.Name())
			if _, err := os.Stat(filepath.Join(themePath, "theme.yaml")); err == nil {
				themes = append(themes, entry.Name())
			}
		}
	}

	slog.Debug("Iterate of themes and add to database")
	for _, name := range themes {
		themePath := GetThemePath(name)
		configFile, err := os.ReadFile(filepath.Join(themePath, "theme.yaml"))
		if err != nil {
			continue
		}

		var config themeConfig
		if err := yaml.Unmarshal(configFile, &config); err != nil {
			slog.Error("Failed to parse theme.yaml", "theme", name, "error", err)
			continue
		}

		description := make(map[string]string)
		if d, ok := config.Description.(string); ok {
			description["en"] = d
		} else if d, ok := config.Description.(map[string]any); ok {
			for k, v := range d {
				if s, ok := v.(string); ok {
					description[k] = s
				}
			}
		}

		dbTheme, err := database.FindThemeByName(name)
		if err != nil {
			slog.Debug("Create theme", "theme", name)
			dbTheme = &database.Theme{
				Name:          name,
				Configuration: map[string]map[string]any{},
				ScssVariables: map[string]string{},
				DisplayName:   config.DisplayName,
				Description:   description,
				Type:          config.Type,
				Path:          themePath,
			}
			if dbTheme.DisplayName == "" {
				dbTheme.DisplayName = name
			}
			err = dbTheme.Create()
			if err != nil {
				slog.Error("Failed to create theme in database", "theme", name, "error", err)
			}
		} else {
			slog.Debug("Update theme", "theme", name)
			dbTheme.DisplayName = config.DisplayName
			if dbTheme.DisplayName == "" {
				dbTheme.DisplayName = name
			}
			dbTheme.Description = description
			dbTheme.Type = config.Type
			dbTheme.Path = themePath
			err = dbTheme.Update()
			if err != nil {
				slog.Error("Failed to update theme in database", "theme", name, "error", err)
			}
		}

		theme, err := NewTheme(dbTheme)
		if err != nil {
			slog.Warn("Invalid theme yaml file", "theme", name, "error", err)
		} else {
			err = theme.CompileAssetCache()
			if err != nil {
				slog.Warn("Failed to compile asset cache", "theme", name, "error", err)
			}
		}
	}

	var nonExistingThemes []database.Theme
	for _, dbTheme := range allThemes {
		if !slices.Contains(themes, dbTheme.Name) {
			nonExistingThemes = append(nonExistingThemes, dbTheme)
		}
	}

	activeTheme, err := database.GetActiveTheme()
	if err != nil {
		slog.Warn("Failed to get active theme", "error", err)
	}

	defaultTheme, _ := database.FindThemeByName("jinya-default-theme")

	for _, nonExistingTheme := range nonExistingThemes {
		if activeTheme == nil || nonExistingTheme.Name == activeTheme.Name {
			if defaultTheme != nil {
				_ = defaultTheme.MakeActiveTheme()
			}
		}
		_ = nonExistingTheme.Delete()
	}

	if activeTheme == nil && defaultTheme != nil {
		_ = defaultTheme.MakeActiveTheme()
	}

	return nil
}
