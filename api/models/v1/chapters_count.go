package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) CountChaptersV1(
	ctx context.Context,
	provider, series string,
) (
	int,
	error,
) {
	chapters, err := p.DB.Chapter.FindMany(
		db.Chapter.ProviderSlug.Equals(provider),
		db.Chapter.SeriesSlug.Equals(series),
	).Select(
		db.Chapter.Slug.Field(),
	).Exec(ctx)
	if err != nil {
		return 0, err
	}

	return len(chapters), nil
}
