package v1Model

import (
	"context"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) FindProviderUniqueV1(ctx context.Context, providerSlug string) (*db.ProviderModel, error) {
	cache, err := p.Redis.FindProviderUniqueV1(ctx, providerSlug)
	if err == nil {
		return cache, nil
	}

	provider, err := p.DB.Provider.FindUnique(
		db.Provider.Slug.Equals(providerSlug),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	_ = p.Redis.CreateProviderUniqueV1(ctx, provider)

	return provider, err
}
