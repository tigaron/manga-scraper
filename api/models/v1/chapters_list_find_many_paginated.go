package v1Model

import (
	"context"

	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindChaptersListPaginatedV1(
	ctx context.Context,
	provider, series string,
	req v1Binding.PaginatedRequest,
	order db.SortOrder,
) (
	result *db.SeriesModel,
	err error,
) {
	result, err = p.Redis.FindChaptersListPaginatedV1(ctx, provider, series, req.Page, req.Size, order)
	if err == nil {
		return
	}

	result, err = p.DB.Series.FindUnique(
		db.Series.SeriesUnique(
			db.Series.ProviderSlug.Equals(provider),
			db.Series.Slug.Equals(series),
		),
	).With(
		db.Series.Chapters.Fetch().OrderBy(
			db.Chapter.Number.Order(order),
		).Take(req.Size).Skip((req.Page-1)*req.Size),
		db.Series.Provider.Fetch(),
	).Exec(ctx)
	if err != nil {
		return
	}

	_ = p.Redis.CreateChaptersListPaginatedV1(ctx, provider, series, req.Page, req.Size, order, result)

	return
}
