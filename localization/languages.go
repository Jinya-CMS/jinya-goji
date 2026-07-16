package localization

import (
	"context"
	_ "embed"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"golang.org/x/text/language"
)

const LanguageRequestContextKey = "language"

//go:embed files/de.json
var german []byte

//go:embed files/en.json
var english []byte

var languageMap map[string]map[string]string

var serverLangs = []language.Tag{
	language.English,
	language.German,
}
var matcher = language.NewMatcher(serverLangs)

func SetupLocalization() {
	languageMap = make(map[string]map[string]string)
	en := make(map[string]string)
	de := make(map[string]string)

	err := json.Unmarshal(german, &de)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(english, &en)
	if err != nil {
		panic(err)
	}

	languageMap["en"] = en
	languageMap["de"] = de
}

func GetRequestLanguage(r *http.Request) string {
	acceptLanguage, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil {
		return "en"
	}
	tag, _, _ := matcher.Match(acceptLanguage...)
	base, _ := tag.Base()
	return base.String()
}

func SetLanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), LanguageRequestContextKey, GetRequestLanguage(r))))
	})
}

func Localize(r *http.Request, key string, replacementValues ...any) string {
	lang := r.Context().Value(LanguageRequestContextKey).(string)
	slog.Debug("Localizing", "language", lang, "key", key, "replacementValues", replacementValues)

	return LocalizeByLanguage(lang, key, replacementValues...)
}

func LocalizeByLanguage(lang, key string, replacementValues ...any) string {
	slog.Debug("Localizing", "language", lang, "key", key, "replacementValues", replacementValues)
	result, ok := languageMap[lang][key]
	if !ok {
		result, ok = languageMap["en"][key]
	}
	if !ok {
		return key
	}

	for i, value := range replacementValues {
		data, ok := value.([]interface{})
		if ok && len(data) >= 1 {
			result = strings.Replace(result, "%"+strconv.Itoa(i+1), cast.ToString(data[0]), -1)
		} else if !ok {
			data, ok := value.(interface{})
			if ok {
				result = strings.Replace(result, "%"+strconv.Itoa(i+1), cast.ToString(data), -1)
			}
		}
	}

	return result
}
