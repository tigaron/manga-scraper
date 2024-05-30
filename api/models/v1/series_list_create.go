package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) CreateSeriesListRowV1(ctx context.Context, id string, data SeriesList) (*db.SeriesListDataModel, error) {
	res, err := p.DB.SeriesListData.CreateOne(
		db.SeriesListData.Title.Set(data.Title),
		db.SeriesListData.Slug.Set(data.Slug),
		db.SeriesListData.SourcePath.Set(data.SourcePath),
		db.SeriesListData.Request.Link(db.ScrapeRequest.ID.Equals(id)),
	).Exec(ctx)

	return res, err
}
