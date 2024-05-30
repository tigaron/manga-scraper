package v1Model

import (
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
	"fourleaves.studio/manga-scraper/internal/database/redis"
)

type DBService struct {
	DB    *db.PrismaClient
	Redis *redis.RedisClient
}

func NewDBService(db *db.PrismaClient, redis *redis.RedisClient) *DBService {
	return &DBService{
		DB:    db,
		Redis: redis,
	}
}
