package i18n

import (
	"errors"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/obumax/gen-secure-password-bot/internal/generator"
)

const (
	ErrLengthID = "length_out_of_range"
	ErrNoCatID  = "no_category_selected"
)

// MapError переводит ошибки из генератора в локализованное сообщение / MapError translates generator errors to localized messages

func MapError(loc *goi18n.Localizer, err error, data map[string]interface{}) string {
	var msgID string
	switch {
	case errors.Is(err, generator.ErrLengthOutOfRange):
		msgID = ErrLengthID
	case errors.Is(err, generator.ErrNoCategorySelected):
		msgID = ErrNoCatID
	default:
		return err.Error()
	}
	s, lerr := loc.Localize(&goi18n.LocalizeConfig{
		MessageID:    msgID,
		TemplateData: data,
	})
	if lerr != nil {
		return err.Error()
	}
	return s
}
