package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindLastChapterV1(
	ctx context.Context,
	provider, series string,
) (
	result *db.ChapterModel,
	err error,
) {
	return p.DB.Chapter.FindFirst(
		db.Chapter.And(
			db.Chapter.ProviderSlug.Equals(provider),
			db.Chapter.SeriesSlug.Equals(series),
			db.Chapter.NextSlug.Equals(""),
			db.Chapter.NextPath.Equals(""),
		),
	).Exec(ctx)
}
