package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) CreateSeriesListScrapeRequestV1(ctx context.Context, provider *db.ProviderModel) (*db.ScrapeRequestModel, error) {
	receipt, err := p.DB.ScrapeRequest.CreateOne(
		db.ScrapeRequest.Type.Set(db.ScrapeRequestTypeSeriesList),
		db.ScrapeRequest.BaseURL.Set(provider.Scheme+provider.Host),
		db.ScrapeRequest.RequestPath.Set(provider.ListPath),
		db.ScrapeRequest.Provider.Set(provider.Slug),
		db.ScrapeRequest.Status.Set("pending"),
	).Exec(ctx)

	return receipt, err
}
