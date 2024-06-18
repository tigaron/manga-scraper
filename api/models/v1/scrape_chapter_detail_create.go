package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) CreateChapterDetailRowV1(ctx context.Context, id string, data ChapterDetail) (*db.ChapterDetailDataModel, error) {
	res, err := p.DB.ChapterDetailData.CreateOne(
		db.ChapterDetailData.FullTitle.Set(data.FullTitle),
		db.ChapterDetailData.SourcePath.Set(data.SourcePath),
		db.ChapterDetailData.NextSlug.Set(data.NextSlug),
		db.ChapterDetailData.NextPath.Set(data.NextPath),
		db.ChapterDetailData.PrevSlug.Set(data.PrevSlug),
		db.ChapterDetailData.PrevPath.Set(data.PrevPath),
		db.ChapterDetailData.ContentPaths.Set(data.ContentPaths),
		db.ChapterDetailData.Request.Link(db.ScrapeRequest.ID.Equals(id)),
	).Exec(ctx)

	return res, err
}
