package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) CreateChapterListScrapeRequestV1(ctx context.Context, provider *db.ProviderModel, series *db.SeriesModel) (*db.ScrapeRequestModel, error) {
	receipt, err := p.DB.ScrapeRequest.CreateOne(
		db.ScrapeRequest.Type.Set(db.ScrapeRequestTypeChapterList),
		db.ScrapeRequest.BaseURL.Set(provider.Scheme+provider.Host),
		db.ScrapeRequest.RequestPath.Set(series.SourcePath),
		db.ScrapeRequest.Provider.Set(provider.Slug),
		db.ScrapeRequest.Series.Set(series.Slug),
		db.ScrapeRequest.Chapter.Set(""),
		db.ScrapeRequest.Status.Set("pending"),
		db.ScrapeRequest.Retries.Set(0),
		db.ScrapeRequest.TotalTime.Set(0),
		db.ScrapeRequest.Error.Set(false),
		db.ScrapeRequest.Message.Set(""),
	).Exec(ctx)

	return receipt, err
}
