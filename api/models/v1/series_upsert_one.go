package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) UpsertSeriesRowV1(ctx context.Context, provider string, data SeriesList) (*db.SeriesModel, error) {
	series, err := p.DB.Series.UpsertOne(
		db.Series.SeriesUnique(
			db.Series.ProviderSlug.Equals(provider),
			db.Series.Slug.Equals(data.Slug),
		),
	).Create(
		db.Series.Slug.Set(data.Slug),
		db.Series.Title.Set(data.Title),
		db.Series.SourcePath.Set(data.SourcePath),
		db.Series.Provider.Link(db.Provider.Slug.Equals(provider)),
	).Update(
		db.Series.Title.Set(data.Title),
		db.Series.SourcePath.Set(data.SourcePath),
	).Exec(ctx)

	return series, err
}
