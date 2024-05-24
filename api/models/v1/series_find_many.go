package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindSeriesManyV1(ctx context.Context, provider string) ([]db.SeriesModel, error) {
	series, err := p.DB.Series.FindMany(
		db.Series.ProviderSlug.Equals(provider),
	).OrderBy(
		db.Series.Slug.Order(db.ASC),
	).Exec(ctx)

	return series, err
}
