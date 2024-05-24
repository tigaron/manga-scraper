package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindSeriesUniqueV1(ctx context.Context, providerSlug string, seriesSlug string) (*db.SeriesModel, error) {
	series, err := p.DB.Series.FindUnique(
		db.Series.SeriesUnique(
			db.Series.ProviderSlug.Equals(providerSlug),
			db.Series.Slug.Equals(seriesSlug),
		),
	).Exec(ctx)

	return series, err
}
