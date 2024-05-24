package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindProviderManyV1(ctx context.Context) ([]db.ProviderModel, error) {
	providers, err := p.DB.Provider.FindMany().OrderBy(
		db.Provider.Slug.Order(db.ASC),
	).Exec(ctx)

	return providers, err
}
