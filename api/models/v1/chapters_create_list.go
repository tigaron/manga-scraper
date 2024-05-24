package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) CreateChapterListRowV1(ctx context.Context, id string, data ChapterList) (*db.ChapterListDataModel, error) {
	res, err := p.DB.ChapterListData.CreateOne(
		db.ChapterListData.ShortTitle.Set(data.ShortTitle),
		db.ChapterListData.Slug.Set(data.Slug),
		db.ChapterListData.Number.Set(data.Number),
		db.ChapterListData.Href.Set(data.Href),
		db.ChapterListData.Request.Link(db.ScrapeRequest.ID.Equals(id)),
	).Exec(ctx)

	return res, err
}
