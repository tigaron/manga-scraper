package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) UpsertChaptersRowV1(ctx context.Context, provider string, series string, data ChapterList) (*db.ChapterModel, error) {
	chapter, err := p.DB.Chapter.UpsertOne(
		db.Chapter.ChapterUnique(
			db.Chapter.ProviderSlug.Equals(provider),
			db.Chapter.SeriesSlug.Equals(series),
			db.Chapter.Slug.Equals(data.Slug),
		),
	).Create(
		db.Chapter.Slug.Set(data.Slug),
		db.Chapter.Number.Set(data.Number),
		db.Chapter.ShortTitle.Set(data.ShortTitle),
		db.Chapter.SourceHref.Set(data.Href),
		db.Chapter.Provider.Link(db.Provider.Slug.Equals(provider)),
		db.Chapter.Series.Link(db.Series.SeriesUnique(
			db.Series.ProviderSlug.Equals(provider),
			db.Series.Slug.Equals(series),
		)),
	).Update(
		db.Chapter.Number.Set(data.Number),
		db.Chapter.ShortTitle.Set(data.ShortTitle),
		db.Chapter.SourceHref.Set(data.Href),
	).Exec(ctx)

	return chapter, err
}
