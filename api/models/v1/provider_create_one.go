package v1Model

import (
	"context"

	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (p *DBService) CreateProviderV1(ctx context.Context, req v1Binding.PostProviderRequest) (*db.ProviderModel, error) {
	provider, err := p.DB.Provider.CreateOne(
		db.Provider.Slug.Set(req.Slug),
		db.Provider.Name.Set(req.Name),
		db.Provider.Scheme.Set(req.Scheme),
		db.Provider.Host.Set(req.Host),
		db.Provider.ListPath.Set(req.ListPath),
		db.Provider.IsActive.Set(req.IsActive),
	).Exec(ctx)

	return provider, err
}
