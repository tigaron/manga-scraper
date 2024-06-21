package service

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/getsentry/sentry-go"
)

type ProviderRepository interface {
	Create(ctx context.Context, params internal.ProviderParams) (internal.Provider, error)
	Find(ctx context.Context, slug string) (internal.Provider, error)
	FindBC(ctx context.Context, slug string) (internal.ProviderBC, error)
	FindAll(ctx context.Context, order internal.SortOrder) ([]internal.Provider, error)
	Update(ctx context.Context, params internal.ProviderParams) (internal.Provider, error)
	Delete(ctx context.Context, slug string) error
}

type ProviderService struct {
	repo ProviderRepository
}

func NewProviderService(repo ProviderRepository) *ProviderService {
	return &ProviderService{
		repo: repo,
	}
}

func (s *ProviderService) Create(ctx context.Context, params internal.ProviderParams) (internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderService.Create").Finish()

	if err := params.Validate(); err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrInvalidInput, "params.Validate")
	}

	provider, err := s.repo.Create(ctx, params)
	if err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.Create")
	}

	return provider, nil
}

func (s *ProviderService) Find(ctx context.Context, slug string) (internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderService.Find").Finish()

	provider, err := s.repo.Find(ctx, slug)
	if err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.Find")
	}

	return provider, nil
}

func (s *ProviderService) FindBC(ctx context.Context, slug string) (internal.ProviderBC, error) {
	defer newSentrySpan(ctx, "ProviderService.Find").Finish()

	provider, err := s.repo.FindBC(ctx, slug)
	if err != nil {
		return internal.ProviderBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.FindBC")
	}

	return provider, nil
}

func (s *ProviderService) FindAll(ctx context.Context, order internal.SortOrder) ([]internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderService.FindAll").Finish()

	providers, err := s.repo.FindAll(ctx, order)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "repo.FindAll")
	}

	return providers, nil
}

func (s *ProviderService) Update(ctx context.Context, params internal.ProviderParams) (internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderService.Update").Finish()

	if err := params.Validate(); err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrInvalidInput, "params.Validate")
	}

	provider, err := s.repo.Update(ctx, params)
	if err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.Update")
	}

	return provider, nil
}

func (s *ProviderService) Delete(ctx context.Context, slug string) error {
	defer newSentrySpan(ctx, "ProviderService.Delete").Finish()

	if err := s.repo.Delete(ctx, slug); err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "repo.Delete")
	}

	return nil
}

func newSentrySpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	span.Name = "fourleaves.studio/manga-scraper/internal/service"

	return span
}
