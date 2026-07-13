package theming

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"

	"jinya.de/goji/database"
)

func getThemeBasePath() string {
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

type ThemeSyncer struct{}

func NewThemeSyncer() *ThemeSyncer {
	return &ThemeSyncer{}
}

type themeConfig struct {
	DisplayName string `json:"displayName"`
	Description any    `json:"description"`
	HasApi      bool   `json:"hasApi"`
}

func (s *ThemeSyncer) SyncThemes() error {
	slog.Debug("Sync all themes")
	allThemes, err := database.FindAllThemes()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(getThemeBasePath())
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("Theme base path does not exist", "path", getThemeBasePath())
			return nil
		}
		return err
	}

	var themes []string
	for _, entry := range entries {
		if entry.IsDir() {
			themePath := filepath.Join(getThemeBasePath(), entry.Name())
			if _, err := os.Stat(filepath.Join(themePath, "theme.yaml")); err == nil {
				themes = append(themes, entry.Name())
			}
		}
	}

	slog.Debug("Iterate of themes and add to database")
	for _, name := range themes {
		themePath := filepath.Join(getThemeBasePath(), name)
		configFile, err := os.ReadFile(filepath.Join(themePath, "theme.yaml"))
		if err != nil {
			continue
		}

		var config themeConfig
		if err := json.Unmarshal(configFile, &config); err != nil {
			slog.Error("Failed to parse theme.json", "theme", name, "error", err)
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
				Configuration: make(map[string]any),
				ScssVariables: make(map[string]string),
				DisplayName:   config.DisplayName,
				Description:   description,
				HasApiTheme:   config.HasApi,
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
			dbTheme.HasApiTheme = config.HasApi
			err = dbTheme.Update()
			if err != nil {
				slog.Error("Failed to update theme in database", "theme", name, "error", err)
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
