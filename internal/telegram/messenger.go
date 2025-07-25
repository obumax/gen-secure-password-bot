package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/obumax/gen-secure-password-bot/internal/generator"
	"github.com/obumax/gen-secure-password-bot/internal/session"
)

// Базовая отправка сообщения / Basic message sending

func send(api *tgbotapi.BotAPI, chatID int64, text string, kb *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	if kb != nil {
		msg.ReplyMarkup = kb
	}
	if _, err := api.Send(msg); err != nil {
		log.Printf("send msg err: %v", err)
	}
}

// Локализация и отправка / Localization and sending

func sendLocalized(api *tgbotapi.BotAPI, chatID int64, loc *goi18n.Localizer,
	id string, data map[string]interface{},
) {
	text, err := loc.Localize(&goi18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if err != nil {
		text = fmt.Sprintf("[%s]", id)
		log.Printf("loc err %q: %v", id, err)
	}
	send(api, chatID, text, nil)
}

// Отправка с кнопками "Еще" и "Новый" / Sending with buttons "Again" and "New"

func sendWithButtons(api *tgbotapi.BotAPI, chatID int64, loc *goi18n.Localizer, pwd string) {
	text := fmt.Sprintf("`%s`", pwd)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "MarkdownV2"
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				mustLocalize(loc, "button.more"),
				"gen:again",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				mustLocalize(loc, "button.new"),
				"gen:new",
			),
		),
	)
	msg.ReplyMarkup = &kb
	if _, err := api.Send(msg); err != nil {
		log.Printf("send pwd err: %v", err)
	}
}

// Кнопки выбора языка / Buttons for language selection

func languageKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang:en"),
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang:ru"),
		),
	)
}

// Кнопки переключения флагов / Buttons for toggling flags

func flagsKeyboard(sess *session.Session, loc *goi18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	isOn := func(code string) bool {
		for _, f := range sess.Flags {
			if f == code {
				return true
			}
		}
		return false
	}
	mark := func(on bool) string {
		if on {
			return "✅"
		}
		return "❌"
	}
	row1 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", mark(isOn("U")), mustLocalize(loc, "button.uppercase")),
			"flag:U",
		),
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", mark(isOn("L")), mustLocalize(loc, "button.lowercase")),
			"flag:L",
		),
	}
	row2 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", mark(isOn("D")), mustLocalize(loc, "button.digits")),
			"flag:D",
		),
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", mark(isOn("S")), mustLocalize(loc, "button.symbols")),
			"flag:S",
		),
	}
	row3 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", mark(isOn("X")), mustLocalize(loc, "button.exclude_similar")),
			"flag:X",
		),
	}
	row4 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			mustLocalize(loc, "button.generate"),
			"gen:run",
		),
	}
	return tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3, row4)
}

// Вспомогательные функции / Helpers

func toggle(flags []string, code string) []string {
	out := make([]string, 0, len(flags))
	found := false
	for _, f := range flags {
		if f == code {
			found = true
		} else {
			out = append(out, f)
		}
	}
	if !found {
		out = append(out, code)
	}
	return out
}

func toGenFlags(codes []string) generator.FlagsSet {
	f := generator.FlagsSet{}
	for _, c := range codes {
		switch c {
		case "U":
			f.Upper = true
		case "L":
			f.Lower = true
		case "D":
			f.Digits = true
		case "S":
			f.SpecSymbols = true
		case "X":
			f.ExcludeSimilar = true
		}
	}
	return f
}

func mustLocalize(loc *goi18n.Localizer, id string) string {
	s, err := loc.Localize(&goi18n.LocalizeConfig{MessageID: id})
	if err != nil {
		log.Printf("localize %s: %v", id, err)
		return id
	}
	return s
}
