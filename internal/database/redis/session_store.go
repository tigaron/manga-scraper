package redis

import (
	"context"

	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
)

func (r *RedisClient) NewSessionStore() (*redisstore.RedisStore, error) {
	store, err := redisstore.NewRedisStore(context.Background(), r.client)
	if err != nil {
		return nil, err
	}

	store.KeyPrefix("session:")
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: 2,
		// Secure:   true,
	})

	return store, nil
}
