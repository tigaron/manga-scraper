package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindChapterUniqueV1(ctx context.Context, providerSlug, seriesSlug, chapterSlug string) (*db.ChapterModel, error) {
	cache, err := p.Redis.FindChapterUniqueV1(ctx, providerSlug, seriesSlug, chapterSlug)
	if err == nil {
		return cache, nil
	}

	chapter, err := p.DB.Chapter.FindUnique(
		db.Chapter.ChapterUnique(
			db.Chapter.ProviderSlug.Equals(providerSlug),
			db.Chapter.SeriesSlug.Equals(seriesSlug),
			db.Chapter.Slug.Equals(chapterSlug),
		),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	_ = p.Redis.CreateChapterUniqueV1(ctx, providerSlug, seriesSlug, chapterSlug, chapter)

	return chapter, err
}
