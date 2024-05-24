package v1Model

import (
	"context"

	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindSeriesManyPaginatedV1(ctx context.Context, provider string, req v1Binding.PaginatedRequest) ([]db.SeriesModel, error) {
	series, err := p.DB.Series.FindMany(
		db.Series.ProviderSlug.Equals(provider),
	).OrderBy(
		db.Series.Slug.Order(db.ASC),
	).Take(req.Size).
		Skip((req.Page - 1) * req.Size).
		Exec(ctx)

	return series, err
}
