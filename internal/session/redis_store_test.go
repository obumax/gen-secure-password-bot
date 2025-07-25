package session

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisStoreGetSetDelete(t *testing.T) {
	client, mock := redismock.NewClientMock()
	store := &RedisStore{
		client: client,
		prefix: "session:",
		ttl:    0,
	}
	chatID := int64(12345)
	key := store.key(chatID)

	// GET → not found
	mock.ExpectGet(key).RedisNil()
	_, err := store.Get(chatID)
	assert.Equal(t, ErrNotFound, err)

	// SET + GET
	sess := &Session{ChatID: chatID, Language: "ru", State: "init", LastActive: time.Now()}
	buf, _ := json.Marshal(sess)
	mock.ExpectSet(key, buf, store.ttl).SetVal("OK")
	mock.ExpectGet(key).SetVal(string(buf))

	assert.NoError(t, store.Set(chatID, sess))
	got, err := store.Get(chatID)
	assert.NoError(t, err)

	gotBuf, _ := json.Marshal(got)
	assert.Equal(t, string(buf), string(gotBuf))

	// DELETE
	mock.ExpectDel(key).SetVal(1)
	assert.NoError(t, store.Delete(chatID))

	assert.NoError(t, mock.ExpectationsWereMet())
}
