package redis

import (
	"context"
	"fmt"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type ProviderStore interface {
	Create(ctx context.Context, params internal.ProviderParams) (internal.Provider, error)
	Find(ctx context.Context, slug string) (internal.Provider, error)
	FindBC(ctx context.Context, slug string) (internal.ProviderBC, error)
	FindAll(ctx context.Context, order internal.SortOrder) ([]internal.Provider, error)
	Update(ctx context.Context, params internal.ProviderParams) (internal.Provider, error)
	Delete(ctx context.Context, slug string) error
}

type ProviderCache struct {
	client     *redis.Client
	store      ProviderStore
	expiration time.Duration
	logger     echo.Logger
}

func NewProviderCache(redisURL string, store ProviderStore, expiration time.Duration, logger echo.Logger) *ProviderCache {
	opts, _ := redis.ParseURL(redisURL)
	return &ProviderCache{
		client:     redis.NewClient(opts),
		store:      store,
		expiration: expiration,
		logger:     logger,
	}
}

func (p *ProviderCache) Create(ctx context.Context, params internal.ProviderParams) (internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderCache.Create").Finish()

	provider, err := p.store.Create(ctx, params)
	if err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.Create")
	}

	cacheKey := fmt.Sprintf("v1:provider:%s", provider.Slug)

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.Create",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   provider,
	})

	err = p.setProvider(ctx, cacheKey, provider)
	if err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "setProvider")
	}

	return provider, nil
}

func (p *ProviderCache) Find(ctx context.Context, slug string) (internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderCache.Find").Finish()

	cacheKey := fmt.Sprintf("v1:provider:%s", slug)

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.Find",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := p.getProvider(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.Find",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	provider, err := p.store.Find(ctx, slug)
	if err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.Find")
	}

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.Find",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   provider,
	})

	err = p.setProvider(ctx, cacheKey, provider)
	if err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "setProvider")
	}

	return provider, nil
}

func (p *ProviderCache) FindBC(ctx context.Context, slug string) (internal.ProviderBC, error) {
	defer newSentrySpan(ctx, "ProviderCache.FindBC").Finish()

	cacheKey := fmt.Sprintf("v1:provider:%s:_bc", slug)

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.FindBC",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := p.getProviderBC(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.FindBC",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	provider, err := p.store.FindBC(ctx, slug)
	if err != nil {
		return internal.ProviderBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.FindBC")
	}

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.FindBC",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   provider,
	})

	err = p.setProviderBC(ctx, cacheKey, provider)
	if err != nil {
		return internal.ProviderBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "setProviderBC")
	}

	return provider, nil
}

func (p *ProviderCache) FindAll(ctx context.Context, order internal.SortOrder) ([]internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderCache.FindAll").Finish()

	cacheKey := "v1:providers:_list"

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.FindAll",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := p.getManyProviders(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.FindAll",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	providers, err := p.store.FindAll(ctx, order)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "store.FindAll")
	}

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.FindAll",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   providers,
	})

	err = p.setManyProviders(ctx, cacheKey, providers)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "setManyProviders")
	}

	return providers, nil
}

func (p *ProviderCache) Update(ctx context.Context, params internal.ProviderParams) (internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderCache.Update").Finish()

	provider, err := p.store.Update(ctx, params)
	if err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.Update")
	}

	cacheKey := fmt.Sprintf("v1:provider:%s", provider.Slug)

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.Update",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   provider,
	})

	err = p.setProvider(ctx, cacheKey, provider)
	if err != nil {
		return internal.Provider{}, internal.WrapErrorf(err, internal.ErrUnknown, "setProvider")
	}

	_ = p.deleteProvider(ctx, "v1:providers:_list")

	return provider, nil
}

func (p *ProviderCache) Delete(ctx context.Context, slug string) error {
	defer newSentrySpan(ctx, "ProviderCache.Delete").Finish()

	cacheKey := fmt.Sprintf("v1:provider:%s", slug)

	p.logger.Debugj(map[string]interface{}{
		"_source": "ProviderCache.Delete",
		"_msg":    "delete cache",
		"key":     cacheKey,
	})

	_ = p.deleteProvider(ctx, cacheKey)
	_ = p.deleteProvider(ctx, "v1:providers:_list")

	return p.store.Delete(ctx, slug)
}
