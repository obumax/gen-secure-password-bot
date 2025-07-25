package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

var bundle *goi18n.Bundle

// InitBundle инициализирует Bundle и грузит все JSON-файлы / InitBundle initializes the Bundle and loads all JSON files

func InitBundle() error {
	bundle = goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files, err := fs.Glob(localeFS, "locales/*.json")
	if err != nil {
		return fmt.Errorf("i18n: glob locale files: %w", err)
	}
	for _, file := range files {
		data, err := localeFS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("i18n: read %s: %w", file, err)
		}
		if _, err := bundle.ParseMessageFileBytes(data, file); err != nil {
			return fmt.Errorf("i18n: parse %s: %w", file, err)
		}
	}
	return nil
}

// Localizer возвращает localizer для кода lang ("en", "ru" и т.д.) / Localizer returns a localizer for the language code (e.g., "en", "ru")

func Localizer(lang string) *goi18n.Localizer {
	if bundle == nil {
		log.Println("i18n: bundle not initialized, initializing")
		if err := InitBundle(); err != nil {
			log.Fatalf("i18n: InitBundle failed: %v", err)
		}
	}
	return goi18n.NewLocalizer(bundle, lang)
}
