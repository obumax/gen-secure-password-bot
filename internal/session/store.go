package session

import "time"

// Session — ваши поля с JSON-тегами / Session is your fields with JSON tags

type Session struct {
	ChatID         int64     `json:"chat_id"`
	Language       string    `json:"language"`
	State          string    `json:"state"`
	PasswordLength int       `json:"password_length"`
	Flags          []string  `json:"flags"`
	LastActive     time.Time `json:"last_active"`
}

// Store — интерфейс хранилища сессий / Store is the session store interface

type Store interface {
	Get(chatID int64) (*Session, error)
	Set(chatID int64, sess *Session) error
	Delete(chatID int64) error
}

var defaultStore Store

// InitStore инициализирует глобальное хранилище / InitStore initializes the global store

func InitStore(s Store) {
	defaultStore = s
}

// Get, Set, Delete — обёртки для defaultStore / Get, Set, Delete wrappers for defaultStore

func Get(chatID int64) (*Session, error) {
	return defaultStore.Get(chatID)
}

func Set(chatID int64, sess *Session) error {
	return defaultStore.Set(chatID, sess)
}

func Delete(chatID int64) error {
	return defaultStore.Delete(chatID)
}
