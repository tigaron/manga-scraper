package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindSeriesUniqueV1(ctx context.Context, providerSlug string, seriesSlug string) (*db.SeriesModel, error) {
	cache, err := p.Redis.FindSeriesUniqueV1(ctx, providerSlug, seriesSlug)
	if err == nil {
		return cache, nil
	}

	series, err := p.DB.Series.FindUnique(
		db.Series.SeriesUnique(
			db.Series.ProviderSlug.Equals(providerSlug),
			db.Series.Slug.Equals(seriesSlug),
		),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	_ = p.Redis.CreateSeriesUniqueV1(ctx, series)

	return series, err
}
