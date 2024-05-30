package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) UpdateDetailChapterRowV1(ctx context.Context, provider, series, chapter string, data ChapterDetail) (*db.ChapterModel, error) {
	res, err := p.DB.Chapter.FindUnique(
		db.Chapter.ChapterUnique(
			db.Chapter.ProviderSlug.Equals(provider),
			db.Chapter.SeriesSlug.Equals(series),
			db.Chapter.Slug.Equals(chapter),
		),
	).Update(
		db.Chapter.FullTitle.Set(data.FullTitle),
		db.Chapter.SourcePath.Set(data.SourcePath),
		db.Chapter.ContentPaths.Set(data.ContentPaths),
		db.Chapter.NextSlug.Set(data.NextSlug),
		db.Chapter.NextPath.Set(data.NextPath),
		db.Chapter.PrevSlug.Set(data.PrevSlug),
		db.Chapter.PrevPath.Set(data.PrevPath),
	).Exec(ctx)

	_ = p.Redis.DeleteChapterUniqueV1(ctx, provider, series, chapter)

	return res, err
}
