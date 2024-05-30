package v1Handler

import (
	v1Model "fourleaves.studio/manga-scraper/api/models/v1"
	"fourleaves.studio/manga-scraper/internal/authenticator"
	"fourleaves.studio/manga-scraper/internal/config"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
	"fourleaves.studio/manga-scraper/internal/database/redis"
)

type Handler struct {
	config *config.Config
	prisma *v1Model.DBService
	redis  *redis.RedisClient
	auth   *authenticator.Authenticator
}

func NewV1Handler(config *config.Config, db *db.PrismaClient, redis *redis.RedisClient, auth *authenticator.Authenticator) *Handler {
	return &Handler{
		config: config,
		prisma: v1Model.NewDBService(db, redis),
		redis:  redis,
		auth:   auth,
	}
}
