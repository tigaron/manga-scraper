package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindChaptersListWithRelV1(
	ctx context.Context,
	provider, series string,
	order db.SortOrder,
) (
	result *db.SeriesModel,
	err error,
) {
	result, err = p.Redis.FindChaptersListWithRelV1(ctx, provider, series, order)
	if err == nil {
		return
	}

	result, err = p.DB.Series.FindUnique(
		db.Series.SeriesUnique(
			db.Series.ProviderSlug.Equals(provider),
			db.Series.Slug.Equals(series),
		),
	).With(
		db.Series.Chapters.Fetch().Select(
			db.Chapter.Slug.Field(),
			db.Chapter.Number.Field(),
			db.Chapter.ShortTitle.Field(),
			db.Chapter.ProviderSlug.Field(),
			db.Chapter.SeriesSlug.Field(),
		).OrderBy(
			db.Chapter.Number.Order(order),
		),
		db.Series.Provider.Fetch(),
	).Exec(ctx)
	if err != nil {
		return
	}

	_ = p.Redis.CreateChaptersListWithRelV1(ctx, provider, series, order, result)

	return
}
