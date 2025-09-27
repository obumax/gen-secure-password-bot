package telegram

import (
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"gen-secure-password-bot/internal/generator"
	"gen-secure-password-bot/internal/i18n"
	"gen-secure-password-bot/internal/session"
)

// Обработка приходящих сообщений / Handle incoming messages

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	sess, err := b.store.Get(chatID)
	if err != nil && err != session.ErrNotFound {
		log.Printf("session.Get: %v", err)
		return
	}
	if sess == nil {
		sess = &session.Session{
			ChatID:         chatID,
			Language:       "en",
			State:          "",
			Flags:          []string{},
			PasswordLength: generator.LengthMin,
		}
	}
	sess.LastActive = time.Now()
	if err := b.store.Set(chatID, sess); err != nil {
		log.Printf("store.Set: %v", err)
	}
	loc := i18n.Localizer(sess.Language)

	// /start

	if msg.IsCommand() && msg.Command() == "start" {
		sess.State = ""
		sess.LastActive = time.Now()
		if err := b.store.Set(chatID, sess); err != nil {
			log.Printf("store.Set: %v", err)
		}

		title := mustLocalize(loc, "start_choose_language")
		kb := languageKeyboard()
		send(b.api, chatID, title, &kb)
		return
	}

	// Если ожидается длина / If waiting for length input

	switch sess.State {
	case "await_len_quick":
		b.onLenQuick(msg.Text, sess)
		return
	case "await_len_custom":
		b.onLenCustom(msg.Text, sess)
		return
	}

	// /lang

	if msg.IsCommand() && msg.Command() == "lang" {
		sess.State = ""
		sess.LastActive = time.Now()
		if err := b.store.Set(chatID, sess); err != nil {
			log.Printf("store.Set: %v", err)
		}

		title := mustLocalize(loc, "start_choose_language")
		kb := languageKeyboard()
		send(b.api, chatID, title, &kb)
		return
	}

	// По умолчанию – главное меню / Default - main menu

	showMainMenu(b, sess)
}

// Обработка Inline-кнопок / Handle inline buttons

func (b *Bot) handleCallback(cq *tgbotapi.CallbackQuery) {
	chatID := cq.Message.Chat.ID
	data := cq.Data

	sess, err := b.store.Get(chatID)
	if err != nil && err != session.ErrNotFound {
		log.Printf("session.Get(cb): %v", err)
		return
	}
	if sess == nil {
		b.api.Request(tgbotapi.NewCallback(cq.ID, ""))
		return
	}
	sess.LastActive = time.Now()
	_ = b.store.Set(chatID, sess)
	loc := i18n.Localizer(sess.Language)

	// Выбор языка / Language selection

	if strings.HasPrefix(data, "lang:") {
		lang := strings.TrimPrefix(data, "lang:")
		sess.Language = lang
		sess.State = ""
		sess.LastActive = time.Now()
		if err := b.store.Set(chatID, sess); err != nil {
			log.Printf("store.Set: %v", err)
		}

		b.api.Request(tgbotapi.NewCallback(cq.ID, ""))
		sendLocalized(b.api, chatID, i18n.Localizer(lang), "start_greeting", nil)
		showMainMenu(b, sess)
		return
	}

	// Основные сценарии / Main scenarios

	switch {
	case data == "gen:quick":
		sess.State = "await_len_quick"
		sess.LastActive = time.Now()
		if err := b.store.Set(chatID, sess); err != nil {
			log.Printf("store.Set: %v", err)
		}

		b.api.Request(tgbotapi.NewCallback(cq.ID, ""))
		sendLocalized(b.api, chatID, loc, "prompt_length", map[string]interface{}{
			"Min": generator.LengthMin,
			"Max": generator.LengthMax,
		})

	case data == "gen:custom":
		sess.State = "await_len_custom"
		sess.LastActive = time.Now()
		if err := b.store.Set(chatID, sess); err != nil {
			log.Printf("store.Set: %v", err)
		}

		b.api.Request(tgbotapi.NewCallback(cq.ID, ""))
		sendLocalized(b.api, chatID, loc, "prompt_length", map[string]interface{}{
			"Min": generator.LengthMin,
			"Max": generator.LengthMax,
		})

	case strings.HasPrefix(data, "flag:"):
		b.api.Request(tgbotapi.NewCallback(cq.ID, ""))

		code := strings.TrimPrefix(data, "flag:")
		sess.Flags = toggle(sess.Flags, code)
		sess.LastActive = time.Now()
		if err := b.store.Set(chatID, sess); err != nil {
			log.Printf("store.Set: %v", err)
		}

		kb := flagsKeyboard(sess, loc)
		edit := tgbotapi.NewEditMessageText(chatID, cq.Message.MessageID,
			mustLocalize(loc, "prompt_flags"))
		edit.ParseMode = "MarkdownV2"
		edit.ReplyMarkup = &kb
		b.api.Send(edit)

	case data == "gen:run":
		b.api.Request(tgbotapi.NewCallback(cq.ID, ""))

		flagsStruct := toGenFlags(sess.Flags)
		pwd, err := generator.Generate(sess.PasswordLength, flagsStruct)
		if err != nil {
			text := i18n.MapError(loc, err, map[string]interface{}{
				"Min": generator.LengthMin,
				"Max": generator.LengthMax,
			})
			send(b.api, chatID, text, nil)
		} else {
			sendWithButtons(b.api, chatID, loc, pwd)
			sess.State = ""
			sess.LastActive = time.Now()
			if err := b.store.Set(chatID, sess); err != nil {
				log.Printf("store.Set: %v", err)
			}
		}

	case data == "gen:again":
		b.api.Request(tgbotapi.NewCallback(cq.ID, ""))

		flagsStruct := toGenFlags(sess.Flags)
		pwd, err := generator.Generate(sess.PasswordLength, flagsStruct)
		if err != nil {
			text := i18n.MapError(loc, err, map[string]interface{}{
				"Min": generator.LengthMin,
				"Max": generator.LengthMax,
			})
			send(b.api, chatID, text, nil)
		} else {
			sendWithButtons(b.api, chatID, loc, pwd)
			sess.LastActive = time.Now()
			if err := b.store.Set(chatID, sess); err != nil {
				log.Printf("store.Set: %v", err)
			}
		}

	case data == "gen:new":
		b.api.Request(tgbotapi.NewCallback(cq.ID, ""))
		sess.LastActive = time.Now()
		if err := b.store.Set(chatID, sess); err != nil {
			log.Printf("store.Set: %v", err)
		}
		showMainMenu(b, sess)
	}
}

// showMainMenu выводит основное меню / Show the main menu

func showMainMenu(b *Bot, sess *session.Session) {
	loc := i18n.Localizer(sess.Language)
	sess.State = ""
	sess.LastActive = time.Now()
	if err := b.store.Set(sess.ChatID, sess); err != nil {
		log.Printf("store.Set: %v", err)
	}

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				mustLocalize(loc, "button.generate_strong"),
				"gen:quick",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				mustLocalize(loc, "button.configure_generate"),
				"gen:custom",
			),
		),
	)
	send(b.api, sess.ChatID, mustLocalize(loc, "start_commands"), &kb)
}

func (b *Bot) onLenQuick(text string, sess *session.Session) {
	chatID := sess.ChatID
	loc := i18n.Localizer(sess.Language)

	n, err := strconv.Atoi(text)
	if err != nil || n < generator.LengthMin || n > generator.LengthMax {
		sendLocalized(b.api, chatID, loc, "length_out_of_range", map[string]interface{}{
			"Min": generator.LengthMin,
			"Max": generator.LengthMax,
		})
		return
	}

	sess.PasswordLength = n
	sess.Flags = []string{"U", "L", "D", "S"}
	sess.LastActive = time.Now()
	if err := b.store.Set(chatID, sess); err != nil {
		log.Printf("store.Set: %v", err)
	}

	pwd, _ := generator.Generate(n, generator.FlagsSet{
		Upper:       true,
		Lower:       true,
		Digits:      true,
		SpecSymbols: true,
	})
	sendWithButtons(b.api, chatID, loc, pwd)

	sess.State = ""
	sess.LastActive = time.Now()
	if err := b.store.Set(chatID, sess); err != nil {
		log.Printf("store.Set: %v", err)
	}
}

func (b *Bot) onLenCustom(text string, sess *session.Session) {
	chatID := sess.ChatID
	loc := i18n.Localizer(sess.Language)

	n, err := strconv.Atoi(text)
	if err != nil || n < generator.LengthMin || n > generator.LengthMax {
		sendLocalized(b.api, chatID, loc, "length_out_of_range", map[string]interface{}{
			"Min": generator.LengthMin,
			"Max": generator.LengthMax,
		})
		return
	}

	sess.PasswordLength = n
	sess.Flags = []string{"U", "L", "D", "S", "X"}
	sess.State = "await_flags"
	sess.LastActive = time.Now()
	if err := b.store.Set(chatID, sess); err != nil {
		log.Printf("store.Set: %v", err)
	}

	kb := flagsKeyboard(sess, loc)
	send(b.api, chatID, mustLocalize(loc, "prompt_flags"), &kb)
}
