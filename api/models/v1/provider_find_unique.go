package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindProviderUniqueV1(ctx context.Context, providerSlug string) (*db.ProviderModel, error) {
	provider, err := p.DB.Provider.FindUnique(
		db.Provider.Slug.Equals(providerSlug),
	).Exec(ctx)

	return provider, err
}
