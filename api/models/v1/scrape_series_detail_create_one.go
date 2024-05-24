package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) CreateSeriesDetailScrapeRequestV1(ctx context.Context, provider *db.ProviderModel, series *db.SeriesModel) (*db.ScrapeRequestModel, error) {
	receipt, err := p.DB.ScrapeRequest.CreateOne(
		db.ScrapeRequest.Type.Set(db.ScrapeRequestTypeSeriesDetail),
		db.ScrapeRequest.BaseURL.Set(provider.Scheme+provider.Host),
		db.ScrapeRequest.RequestPath.Set(series.SourcePath),
		db.ScrapeRequest.Provider.Set(provider.Slug),
		db.ScrapeRequest.Series.Set(series.Slug),
		db.ScrapeRequest.Status.Set("pending"),
	).Exec(ctx)

	return receipt, err
}
