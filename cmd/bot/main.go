package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"

	"github.com/obumax/pet-password-generator/internal/generator"
	i18nutil "github.com/obumax/pet-password-generator/internal/i18n"
	"github.com/obumax/pet-password-generator/internal/session"
)

const (
	// LengthMin and LengthMax define the allowed password length range
	LengthMin = 4
	LengthMax = 35
)

var (
	// sessionStore holds the Redis-backed session store
	sessionStore session.Store
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found (or failed to read): %v", err)
	}
}

func main() {

	// Initialize i18n bundle for translations
	if err := i18nutil.InitBundle(); err != nil {
		log.Fatalf("i18n init failed: %v", err)
	}

	// Configure Redis-based session store
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL is not set")
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("invalid REDIS_URL: %v", err)
	}
	sessionStore = session.NewRedisStore(opt.Addr, opt.DB, opt.Password)
	session.InitStore(sessionStore)

	// Initialize Telegram bot
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("NewBotAPI error: %v", err)
	}
	log.Printf("Authorized on %s", bot.Self.UserName)

	// Start receiving updates
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 30
	updates := bot.GetUpdatesChan(ucfg)

	for upd := range updates {
		if upd.CallbackQuery != nil {
			handleCallback(bot, upd.CallbackQuery)
		} else if upd.Message != nil {
			handleMessage(bot, upd.Message)
		}
	}
}

// handleMessage processes incoming text messages and commands
func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	// Load or create session
	sess, err := sessionStore.Get(chatID)
	if err != nil && err != session.ErrNotFound {
		log.Printf("session.Get: %v", err)
		return
	}
	if sess == nil {
		sess = &session.Session{
			ChatID:   chatID,
			Language: "en",
			State:    "",
			Flags:    []string{},
		}
		if err := sessionStore.Set(chatID, sess); err != nil {
			log.Printf("session.Set init: %v", err)
		}
	}
	loc := i18nutil.Localizer(sess.Language)

	// /start resets FSM and prompts for language selection
	if msg.IsCommand() && msg.Command() == "start" {
		sess.State = ""
		_ = sessionStore.Set(chatID, sess)

		title := mustLocalize(loc, "start_choose_language")
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang:en"),
				tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang:ru"),
			),
		)
		send(bot, chatID, title, &kb)
		return
	}

	// If waiting for quick length or custom length, delegate
	switch sess.State {
	case "await_len_quick":
		onLenQuick(bot, msg.Text, sess)
		return
	case "await_len_custom":
		onLenCustom(bot, msg.Text, sess)
		return
	}

	// /lang triggers language re-selection
	if msg.IsCommand() && msg.Command() == "lang" {
		sess.State = ""
		_ = sessionStore.Set(chatID, sess)

		title := mustLocalize(loc, "start_choose_language")
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang:en"),
				tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang:ru"),
			),
		)
		send(bot, chatID, title, &kb)
		return
	}

	// Otherwise show the main menu
	showMainMenu(bot, sess)
}

// handleCallback processes inline button callbacks
func handleCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	chatID := cq.Message.Chat.ID
	data := cq.Data

	sess, err := sessionStore.Get(chatID)
	if err != nil && err != session.ErrNotFound {
		log.Printf("session.Get(cb): %v", err)
		return
	}
	if sess == nil {
		bot.Request(tgbotapi.NewCallback(cq.ID, ""))
		return
	}
	loc := i18nutil.Localizer(sess.Language)

	// Language selection callback
	if strings.HasPrefix(data, "lang:") {
		lang := strings.TrimPrefix(data, "lang:")
		sess.Language = lang
		sess.State = ""
		_ = sessionStore.Set(chatID, sess)

		bot.Request(tgbotapi.NewCallback(cq.ID, ""))
		sendLocalized(bot, chatID, i18nutil.Localizer(lang), "start_greeting", nil)
		sendLocalized(bot, chatID, i18nutil.Localizer(lang), "start_commands", nil)
		showMainMenu(bot, sess)
		return
	}

	// Quick generate callback
	if data == "gen:quick" {
		sess.State = "await_len_quick"
		_ = sessionStore.Set(chatID, sess)

		bot.Request(tgbotapi.NewCallback(cq.ID, ""))
		sendLocalized(bot, chatID, loc, "prompt_length", map[string]interface{}{
			"Min": LengthMin, "Max": LengthMax,
		})
		return
	}

	// Custom generate callback
	if data == "gen:custom" {
		sess.State = "await_len_custom"
		_ = sessionStore.Set(chatID, sess)

		bot.Request(tgbotapi.NewCallback(cq.ID, ""))
		sendLocalized(bot, chatID, loc, "prompt_length", map[string]interface{}{
			"Min": LengthMin, "Max": LengthMax,
		})
		return
	}

	// Toggle flag callback (U/L/D/S/X)
	if strings.HasPrefix(data, "flag:") {
		bot.Request(tgbotapi.NewCallback(cq.ID, ""))

		code := strings.TrimPrefix(data, "flag:")
		sess.Flags = toggle(sess.Flags, code)
		_ = sessionStore.Set(chatID, sess)

		kb := flagsKeyboard(sess.Flags, loc)
		edit := tgbotapi.NewEditMessageText(chatID, cq.Message.MessageID,
			mustLocalize(loc, "prompt_flags"))
		edit.ReplyMarkup = &kb
		bot.Send(edit)
		return
	}

	// Run custom generation callback
	if data == "gen:run" {
		bot.Request(tgbotapi.NewCallback(cq.ID, ""))

		flagsStruct := toGenFlags(sess.Flags)
		pwd, err := generator.Generate(sess.PasswordLength, flagsStruct)
		if err != nil {
			text := i18nutil.MapError(loc, err, map[string]interface{}{
				"Min": LengthMin, "Max": LengthMax,
			})
			send(bot, chatID, text, nil)
			return
		}
		sendWithButtons(bot, chatID, loc, pwd)

		sess.State = ""
		_ = sessionStore.Set(chatID, sess)
		return
	}

	// "Again" callback
	if data == "gen:again" {
		bot.Request(tgbotapi.NewCallback(cq.ID, ""))

		flagsStruct := toGenFlags(sess.Flags)
		pwd, err := generator.Generate(sess.PasswordLength, flagsStruct)
		if err != nil {
			text := i18nutil.MapError(loc, err, map[string]interface{}{
				"Min": LengthMin, "Max": LengthMax,
			})
			send(bot, chatID, text, nil)
			return
		}
		sendWithButtons(bot, chatID, loc, pwd)
		return
	}

	// "New generator" callback
	if data == "gen:new" {
		bot.Request(tgbotapi.NewCallback(cq.ID, ""))
		showMainMenu(bot, sess)
		return
	}
}

// showMainMenu displays the two main actions: quick/custom generate
func showMainMenu(bot *tgbotapi.BotAPI, sess *session.Session) {
	loc := i18nutil.Localizer(sess.Language)
	sess.State = ""
	_ = sessionStore.Set(sess.ChatID, sess)

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				mustLocalize(loc, "btn_quick"), "gen:quick"),
			tgbotapi.NewInlineKeyboardButtonData(
				mustLocalize(loc, "btn_custom"), "gen:custom"),
		),
	)
	send(bot, sess.ChatID, mustLocalize(loc, "start_commands"), &kb)
}

// onLenQuick handles length input for quick generation
func onLenQuick(bot *tgbotapi.BotAPI, text string, sess *session.Session) {
	chatID := sess.ChatID
	loc := i18nutil.Localizer(sess.Language)

	n, err := strconv.Atoi(text)
	if err != nil || n < LengthMin || n > LengthMax {
		sendLocalized(bot, chatID, loc, "length_out_of_range", map[string]interface{}{
			"Min": LengthMin, "Max": LengthMax,
		})
		return
	}

	sess.PasswordLength = n
	sess.Flags = []string{"U", "L", "D", "S"} // all categories
	_ = sessionStore.Set(chatID, sess)

	pwd, _ := generator.Generate(n, generator.FlagsSet{
		Upper: true, Lower: true, Digits: true, SpecSymbols: true,
	})
	sendWithButtons(bot, chatID, loc, pwd)

	sess.State = ""
	_ = sessionStore.Set(chatID, sess)
}

// onLenCustom handles length input for custom generation
func onLenCustom(bot *tgbotapi.BotAPI, text string, sess *session.Session) {
	chatID := sess.ChatID
	loc := i18nutil.Localizer(sess.Language)

	n, err := strconv.Atoi(text)
	if err != nil || n < LengthMin || n > LengthMax {
		sendLocalized(bot, chatID, loc, "length_out_of_range", map[string]interface{}{
			"Min": LengthMin, "Max": LengthMax,
		})
		return
	}

	sess.PasswordLength = n
	// initialize flags list including exclude-similar
	sess.Flags = []string{"U", "L", "D", "S", "X"}
	sess.State = "await_flags"
	_ = sessionStore.Set(chatID, sess)

	kb := flagsKeyboard(sess.Flags, loc)
	send(bot, chatID, mustLocalize(loc, "prompt_flags"), &kb)
}

// send sends a message with optional inline keyboard
func send(bot *tgbotapi.BotAPI, chatID int64, text string, kb *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	if kb != nil {
		msg.ReplyMarkup = kb
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("send msg err: %v", err)
	}
}

// sendLocalized localizes the message ID with data and sends it
func sendLocalized(bot *tgbotapi.BotAPI, chatID int64, loc *goi18n.Localizer,
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
	send(bot, chatID, text, nil)
}

// sendWithButtons sends the generated password in monospaced block plus action buttons
func sendWithButtons(bot *tgbotapi.BotAPI, chatID int64, loc *goi18n.Localizer, pwd string) {
	text := fmt.Sprintf("```%s```", pwd)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "MarkdownV2"
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(mustLocalize(loc, "btn_again"), "gen:again"),
			tgbotapi.NewInlineKeyboardButtonData(mustLocalize(loc, "btn_new"), "gen:new"),
		),
	)
	msg.ReplyMarkup = &kb
	if _, err := bot.Send(msg); err != nil {
		log.Printf("send pwd err: %v", err)
	}
}

// flagsKeyboard builds inline toggles for categories U/L/D/S/X plus Generate button
func flagsKeyboard(flags []string, loc *goi18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	isOn := func(code string) bool {
		for _, f := range flags {
			if f == code {
				return true
			}
		}
		return false
	}
	mark := func(on bool) string {
		if on {
			return "✔"
		}
		return "❌"
	}
	// rows of toggle buttons
	row1 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", mark(isOn("U")), mustLocalize(loc, "flag_upper")), "flag:U"),
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", mark(isOn("L")), mustLocalize(loc, "flag_lower")), "flag:L"),
	}
	row2 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", mark(isOn("D")), mustLocalize(loc, "flag_digits")), "flag:D"),
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", mark(isOn("S")), mustLocalize(loc, "flag_symbols")), "flag:S"),
	}
	row3 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", mark(isOn("X")), mustLocalize(loc, "flag_exclude")), "flag:X"),
	}
	row4 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(mustLocalize(loc, "btn_generate"), "gen:run"),
	}
	return tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3, row4)
}

// toggle flips presence of code in flags slice
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

// toGenFlags converts slice of codes into generator.FlagsSet
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

// mustLocalize returns localized text or the message ID on error
func mustLocalize(loc *goi18n.Localizer, id string) string {
	s, err := loc.Localize(&goi18n.LocalizeConfig{MessageID: id})
	if err != nil {
		log.Printf("localize %s: %v", id, err)
		return id
	}
	return s
}
