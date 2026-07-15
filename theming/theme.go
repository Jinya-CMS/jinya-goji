package theming

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
	"jinya.de/goji/database"
	"jinya.de/goji/frontend"
)

type Theme struct {
	dbTheme       *database.Theme
	configuration map[string]any
}

func NewTheme(dbTheme *database.Theme) (*Theme, error) {
	t := &Theme{
		dbTheme: dbTheme,
	}

	err := t.parseThemeYaml()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Theme) parseThemeYaml() error {
	slog.Debug("Parse the theme.yaml file")
	themeYamlFile := filepath.Join(getThemeBasePath(), t.dbTheme.Name, "theme.yaml")
	if _, err := os.Stat(themeYamlFile); os.IsNotExist(err) {
		slog.Error("Could not find theme.yaml", "path", themeYamlFile)
		return fmt.Errorf("theme.yaml is not found, expecting it here: %s", themeYamlFile)
	}

	slog.Debug("Load theme.yaml")
	data, err := os.ReadFile(themeYamlFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &t.configuration)
}

func (t *Theme) GetConfigurationValues() map[string]any {
	if config, ok := t.configuration["configuration"].(map[string]any); ok {
		return config
	}
	return make(map[string]any)
}

func (t *Theme) GetConfigurationStructure() map[string]any {
	if structure, ok := t.configuration["configurationStructure"].(map[string]any); ok {
		return structure
	}
	return make(map[string]any)
}

func (t *Theme) GetPreviewImagePath() string {
	if preview, ok := t.configuration["previewImage"].(string); ok {
		basePath := getThemeBasePath()

		return path.Join(basePath, t.dbTheme.Name, preview)
	}
	return ""
}

func (t *Theme) GetStyleVariables() (map[string]string, error) {
	styles, ok := t.configuration["styles"].(map[string]any)
	if !ok {
		return nil, nil
	}

	variablesPath, ok := styles["variables"].(string)
	if !ok {
		return nil, nil
	}

	absVariablesPath := filepath.Join(getThemeBasePath(), t.dbTheme.Name, variablesPath)
	if _, err := os.Stat(absVariablesPath); os.IsNotExist(err) {
		return nil, nil
	}

	file, err := os.Open(absVariablesPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	variables := make(map[string]string)
	scanner := bufio.NewScanner(file)
	// Match lines like $variable: value !default;
	re := regexp.MustCompile(`^(\$.*):(.*)!default;\s*$`)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			key := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])
			variables[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return variables, nil
}

func (t *Theme) CompileAssetCache() error {
	slog.Debug("Compile theme asset cache", "name", t.dbTheme.Name)

	assets, ok := t.configuration["assets"].(map[string]any)
	if !ok {
		return nil
	}

	slog.Debug("Copy all assets to the database", "name", t.dbTheme.Name)
	for key, assetPathAny := range assets {
		assetPath, ok := assetPathAny.(string)
		if !ok {
			continue
		}

		absAssetPath := filepath.Join(getThemeBasePath(), t.dbTheme.Name, assetPath)
		assetFromDb, err := database.FindThemeAssetByThemeAndName(t.dbTheme.Id, key)

		if _, errStat := os.Stat(absAssetPath); os.IsNotExist(errStat) {
			if err == nil && assetFromDb != nil {
				_ = assetFromDb.Delete()
			}
			continue
		}

		slog.Debug("Store asset in database", "name", t.dbTheme.Name, "asset", assetPath)
		if err != nil || assetFromDb == nil {
			assetFromDb = &database.ThemeAsset{
				Name:    key,
				Path:    absAssetPath,
				ThemeId: t.dbTheme.Id,
			}
			_ = assetFromDb.Create()
		} else {
			assetFromDb.Path = absAssetPath
			_ = assetFromDb.Update()
		}
	}

	return frontend.CompileTheme(t.dbTheme.Id)
}

func DeleteTheme(name string) error {
	basePath := getThemeBasePath()
	themePath := path.Join(basePath, name)
	return os.RemoveAll(themePath)
}
