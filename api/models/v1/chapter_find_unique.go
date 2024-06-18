package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindChapterUniqueV1(
	ctx context.Context,
	provider, series, chapter string,
) (
	result *db.ChapterModel,
	err error,
) {
	result, err = p.Redis.FindChapterUniqueV1(ctx, provider, series, chapter)
	if err == nil {
		return
	}

	result, err = p.DB.Chapter.FindUnique(
		db.Chapter.ChapterUnique(
			db.Chapter.ProviderSlug.Equals(provider),
			db.Chapter.SeriesSlug.Equals(series),
			db.Chapter.Slug.Equals(chapter),
		),
	).Exec(ctx)
	if err != nil {
		return
	}

	_ = p.Redis.CreateChapterUniqueV1(ctx, provider, series, chapter, result)

	return
}

func (p *DBService) FindChapterUniqueWithRelV1(
	ctx context.Context,
	provider, series, chapter string,
) (
	result *db.ChapterModel,
	err error,
) {
	result, err = p.Redis.FindChapterUniqueWithRelV1(ctx, provider, series, chapter)
	if err == nil {
		return
	}

	result, err = p.DB.Chapter.FindUnique(
		db.Chapter.ChapterUnique(
			db.Chapter.ProviderSlug.Equals(provider),
			db.Chapter.SeriesSlug.Equals(series),
			db.Chapter.Slug.Equals(chapter),
		),
	).With(
		db.Chapter.Series.Fetch(),
		db.Chapter.Provider.Fetch(),
	).Exec(ctx)
	if err != nil {
		return
	}

	_ = p.Redis.CreateChapterUniqueWithRelV1(ctx, provider, series, chapter, result)

	return
}
