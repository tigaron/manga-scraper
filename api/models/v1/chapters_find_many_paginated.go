package v1Model

import (
	"context"

	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindChaptersManyPaginatedV1(ctx context.Context, provider string, series string, req v1Binding.PaginatedRequest) ([]db.ChapterModel, error) {
	chapterList, err := p.DB.Chapter.FindMany(
		db.Chapter.ProviderSlug.Equals(provider),
		db.Chapter.SeriesSlug.Equals(series),
	).OrderBy(
		db.Chapter.Number.Order(db.ASC),
	).Take(req.Size).
		Skip((req.Page - 1) * req.Size).
		Exec(ctx)

	return chapterList, err
}
