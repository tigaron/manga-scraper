package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) UpdateChaptersCountV1(ctx context.Context, provider string, series string, count int) (*db.SeriesModel, error) {
	res, err := p.DB.Series.FindUnique(
		db.Series.SeriesUnique(
			db.Series.ProviderSlug.Equals(provider),
			db.Series.Slug.Equals(series),
		),
	).Update(
		db.Series.ChaptersCount.Set(count),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	_ = p.Redis.DeleteSeriesUniqueV1(ctx, provider, series)

	return res, nil
}

func (p *DBService) IncrementChaptersCountV1(ctx context.Context, provider string, series string, count int) (*db.SeriesModel, error) {
	res, err := p.DB.Series.FindUnique(
		db.Series.SeriesUnique(
			db.Series.ProviderSlug.Equals(provider),
			db.Series.Slug.Equals(series),
		),
	).Update(
		db.Series.ChaptersCount.Increment(count),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	_ = p.Redis.DeleteSeriesUniqueV1(ctx, provider, series)

	return res, nil
}
