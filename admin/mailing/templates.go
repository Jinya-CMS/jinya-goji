package mailing

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	textTmpl "text/template"

	"jinya.de/goji/localization"
)

var (
	//go:embed templates
	templatesFs embed.FS
)

func renderPlain(r *http.Request, name string, data any) (string, error) {
	lang := r.Context().Value(localization.LanguageRequestContextKey).(string)
	return renderPlainByLanguage(lang, name, data)
}

func renderTemplate(r *http.Request, name string, data any) (string, error) {
	lang := r.Context().Value(localization.LanguageRequestContextKey).(string)
	return renderTemplateByLanguage(lang, name, data)
}

func renderPlainByLanguage(lang, name string, data any) (string, error) {
	tmpl, err := textTmpl.
		New("mail").
		Funcs(map[string]any{
			"localize": func(key string, parameters ...any) string {
				return localization.LocalizeByLanguage(lang, key, parameters)
			},
		}).
		ParseFS(templatesFs, fmt.Sprintf("templates/layout.gomjml"), fmt.Sprintf("templates/%s.gotext", name))

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func renderTemplateByLanguage(lang, name string, data any) (string, error) {
	slog.Debug("Render template", "name", name)
	t, err := template.
		New("email").
		Funcs(map[string]any{
			"localize": func(key string, parameters ...any) string {
				return localization.LocalizeByLanguage(lang, key, parameters)
			},
		}).
		ParseFS(templatesFs, fmt.Sprintf("templates/layout.gomjml"), fmt.Sprintf("templates/%s.gomjml", name))
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, data)
	return buffer.String(), err
}
