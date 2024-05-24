package v1Model

import (
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

type DBService struct {
	DB *db.PrismaClient
}

func NewDBService(db *db.PrismaClient) *DBService {
	return &DBService{
		DB: db,
	}
}
