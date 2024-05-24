package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) UpdateScrapeRequestUniqueV1(ctx context.Context, id string, status string, duration float64, message string) (*db.ScrapeRequestModel, error) {
	receipt, err := p.DB.ScrapeRequest.FindUnique(
		db.ScrapeRequest.ID.Equals(id),
	).Update(
		db.ScrapeRequest.Status.Set(status),
		db.ScrapeRequest.TotalTime.Set(duration),
		db.ScrapeRequest.Message.Set(message),
	).Exec(ctx)

	return receipt, err
}
