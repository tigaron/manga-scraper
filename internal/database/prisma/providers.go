package prisma

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
)

type ProviderRepo struct {
	q *PrismaClient
}

func NewProviderRepo(prismaClient *PrismaClient) *ProviderRepo {
	return &ProviderRepo{
		q: prismaClient,
	}
}

func (p *ProviderModel) toProvider() internal.Provider {
	return internal.Provider{
		Slug:     p.Slug,
		Name:     p.Name,
		IsActive: p.IsActive,
		BaseURL:  p.Scheme + p.Host,
		ListURL:  p.Scheme + p.Host + p.ListPath,
	}
}

func (p *ProviderModel) toBC() internal.ProviderBC {
	return internal.ProviderBC{
		Provider: internal.Breadcrumb{
			Slug:  p.Slug,
			Title: p.Name,
		},
	}
}

func (p *ProviderRepo) Create(ctx context.Context, params internal.ProviderParams) (internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderRepo.Create").Finish()

	provider, err := p.q.Provider.CreateOne(
		Provider.Slug.Set(params.Slug),
		Provider.Name.Set(params.Name),
		Provider.Scheme.Set(params.Scheme),
		Provider.Host.Set(params.Host),
		Provider.ListPath.Set(params.ListPath),
		Provider.IsActive.Set(*params.IsActive),
	).Exec(ctx)
	if err != nil {
		if _, ok := IsErrUniqueConstraint(err); ok {
			return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUniqueConstraint, "provider already exists")
		}

		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrInvalidInput, "failed to create provider")
	}

	return provider.toProvider(), nil
}

func (p *ProviderRepo) Find(ctx context.Context, slug string) (internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderRepo.Find").Finish()

	provider, err := p.q.Provider.FindUnique(
		Provider.Slug.Equals(slug),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.Provider{}, internal.WrapErrorf(err, internal.ErrNotFound, "provider not found")
		}

		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find provider")
	}

	return provider.toProvider(), nil
}

func (p *ProviderRepo) FindBC(ctx context.Context, slug string) (internal.ProviderBC, error) {
	defer newSentrySpan(ctx, "ProviderRepo.FindBC").Finish()

	provider, err := p.q.Provider.FindUnique(
		Provider.Slug.Equals(slug),
	).Select(
		Provider.Slug.Field(),
		Provider.Name.Field(),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.ProviderBC{}, internal.WrapErrorf(err, internal.ErrNotFound, "provider not found")
		}

		return internal.ProviderBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find provider")
	}

	return provider.toBC(), nil
}

func (p *ProviderRepo) FindAll(ctx context.Context, order internal.SortOrder) ([]internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderRepo.FindAll").Finish()

	providers, err := p.q.Provider.FindMany().OrderBy(
		Provider.Slug.Order(newSortOrder(order)),
	).Exec(ctx)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find providers")
	}

	if len(providers) == 0 {
		return nil, internal.WrapErrorf(nil, internal.ErrNotFound, "no providers found")
	}

	var result []internal.Provider
	for _, provider := range providers {
		result = append(result, provider.toProvider())
	}

	return result, nil
}

func (p *ProviderRepo) Update(ctx context.Context, params internal.ProviderParams) (internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderRepo.Update").Finish()

	provider, err := p.q.Provider.FindUnique(
		Provider.Slug.Equals(params.Slug),
	).Update(
		Provider.Name.Set(params.Name),
		Provider.Scheme.Set(params.Scheme),
		Provider.Host.Set(params.Host),
		Provider.ListPath.Set(params.ListPath),
		Provider.IsActive.Set(*params.IsActive),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.Provider{}, internal.WrapErrorf(err, internal.ErrNotFound, "provider not found")
		}

		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to update provider")
	}

	return provider.toProvider(), nil
}

func (p *ProviderRepo) Delete(ctx context.Context, slug string) error {
	defer newSentrySpan(ctx, "ProviderRepo.Delete").Finish()

	_, err := p.q.Provider.FindUnique(
		Provider.Slug.Equals(slug),
	).Delete().Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.WrapErrorf(err, internal.ErrNotFound, "provider not found")
		}

		return internal.WrapErrorf(err, internal.ErrUnknown, "failed to delete provider")
	}

	return nil
}
