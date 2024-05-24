package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) CreateChapterDetailScrapeRequestV1(ctx context.Context, provider *db.ProviderModel, series *db.SeriesModel, chapter *db.ChapterModel) (*db.ScrapeRequestModel, error) {
	sourcePath, ok := chapter.SourcePath()
	if !ok {
		sourcePath = chapter.SourceHref
	}

	receipt, err := p.DB.ScrapeRequest.CreateOne(
		db.ScrapeRequest.Type.Set(db.ScrapeRequestTypeChapterDetail),
		db.ScrapeRequest.BaseURL.Set(provider.Scheme+provider.Host),
		db.ScrapeRequest.RequestPath.Set(sourcePath),
		db.ScrapeRequest.Provider.Set(provider.Slug),
		db.ScrapeRequest.Series.Set(series.Slug),
		db.ScrapeRequest.Chapter.Set(chapter.Slug),
		db.ScrapeRequest.Status.Set("pending"),
	).Exec(ctx)

	return receipt, err
}
