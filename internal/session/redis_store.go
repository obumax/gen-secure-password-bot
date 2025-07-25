package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrNotFound = fmt.Errorf("session not found")

type RedisStore struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

// NewRedisStore создает новый RedisStore с заданными параметрами / NewRedisStore creates a new RedisStore with the given parameters

func NewRedisStore(opts *redis.Options) *RedisStore {
	client := redis.NewClient(opts)
	return &RedisStore{
		client: client,
		prefix: "session:",
		ttl:    0,
	}
}

// key генерирует ключ для Redis по chatID / key generates a Redis key based on chatID

func (r *RedisStore) key(chatID int64) string {
	return fmt.Sprintf("%s%d", r.prefix, chatID)
}

// Get получает сессию по chatID / Get retrieves a session by chatID

func (r *RedisStore) Get(chatID int64) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	data, err := r.client.Get(ctx, r.key(chatID)).Result()
	if err == redis.Nil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var sess Session
	if err := json.Unmarshal([]byte(data), &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

// Set сохраняет сессию в Redis по chatID / Set saves a session to Redis by chatID

func (r *RedisStore) Set(chatID int64, sess *Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	buf, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(chatID), buf, r.ttl).Err()
}

// Delete удаляет сессию по chatID / Delete removes a session by chatID

func (r *RedisStore) Delete(chatID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := r.client.Del(ctx, r.key(chatID)).Result()
	return err
}
