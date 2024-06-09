package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindChaptersListV1(ctx context.Context, provider string, series string) ([]db.ChapterModel, error) {
	chapterList, err := p.DB.Chapter.FindMany(
		db.Chapter.ProviderSlug.Equals(provider),
		db.Chapter.SeriesSlug.Equals(series),
	).Select(
		db.Chapter.ID.Field(),
		db.Chapter.Slug.Field(),
		db.Chapter.Number.Field(),
		db.Chapter.ShortTitle.Field(),
		db.Chapter.ProviderSlug.Field(),
		db.Chapter.SeriesSlug.Field(),
	).OrderBy(
		db.Chapter.Number.Order(db.ASC),
	).Exec(ctx)

	return chapterList, err
}
