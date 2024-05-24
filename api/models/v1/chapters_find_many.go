package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindChaptersManyV1(ctx context.Context, provider string, series string) ([]db.ChapterModel, error) {
	chapterList, err := p.DB.Chapter.FindMany(
		db.Chapter.ProviderSlug.Equals(provider),
		db.Chapter.SeriesSlug.Equals(series),
	).OrderBy(
		db.Chapter.Number.Order(db.ASC),
	).Exec(ctx)

	return chapterList, err
}
