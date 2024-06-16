package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindLastChapterV1(ctx context.Context, providerSlug, seriesSlug string) (*db.ChapterModel, error) {
	chapter, err := p.DB.Chapter.FindFirst(
		db.Chapter.ProviderSlug.Equals(providerSlug),
		db.Chapter.SeriesSlug.Equals(seriesSlug),
		db.Chapter.NextSlug.Equals(""),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return chapter, err
}
