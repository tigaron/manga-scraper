package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) CreateSeriesDetailRowV1(ctx context.Context, id string, data SeriesDetail) (*db.SeriesDetailDataModel, error) {
	res, err := p.DB.SeriesDetailData.CreateOne(
		db.SeriesDetailData.ThumbnailURL.Set(data.ThumbnailURL),
		db.SeriesDetailData.Synopsis.Set(data.Synopsis),
		db.SeriesDetailData.Genres.Set(data.Genres),
		db.SeriesDetailData.Request.Link(db.ScrapeRequest.ID.Equals(id)),
	).Exec(ctx)

	return res, err
}
