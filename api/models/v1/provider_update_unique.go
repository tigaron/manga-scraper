package v1Model

import (
	"context"

	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) UpdateProviderUniqueV1(ctx context.Context, providerSlug string, req v1Binding.PutProviderRequest) (*db.ProviderModel, error) {
	provider, err := p.DB.Provider.FindUnique(
		db.Provider.Slug.Equals(providerSlug),
	).Update(
		db.Provider.Name.Set(req.Name),
		db.Provider.Scheme.Set(req.Scheme),
		db.Provider.Host.Set(req.Host),
		db.Provider.ListPath.Set(req.ListPath),
		db.Provider.IsActive.Set(*req.IsActive),
	).Exec(ctx)

	_ = p.Redis.DeleteProviderUniqueV1(ctx, providerSlug)

	return provider, err
}
