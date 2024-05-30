package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) UpdateDetailSeriesRowV1(ctx context.Context, provider string, series string, data SeriesDetail) (*db.SeriesModel, error) {
	res, err := p.DB.Series.FindUnique(
		db.Series.SeriesUnique(
			db.Series.ProviderSlug.Equals(provider),
			db.Series.Slug.Equals(series),
		),
	).Update(
		db.Series.ThumbnailURL.Set(data.ThumbnailURL),
		db.Series.Synopsis.Set(data.Synopsis),
		db.Series.Genres.Set(data.Genres),
	).Exec(ctx)

	_ = p.Redis.DeleteSeriesUniqueV1(ctx, provider, series)

	return res, err
}
