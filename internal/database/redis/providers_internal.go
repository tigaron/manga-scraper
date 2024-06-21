package redis

import (
	"bytes"
	"context"
	"encoding/gob"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/getsentry/sentry-go"
)

func newSentrySpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	span.Name = "fourleaves.studio/manga-scraper/internal/database/redis"
	span.SetTag("db.system", "redis")

	return span
}

func (p *ProviderCache) setProvider(ctx context.Context, key string, value internal.Provider) error {
	defer newSentrySpan(ctx, "ProviderCache.setProvider").Finish()

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return err
	}

	return p.client.Set(ctx, key, b.Bytes(), p.expiration).Err()
}

func (p *ProviderCache) setProviderBC(ctx context.Context, key string, value internal.ProviderBC) error {
	defer newSentrySpan(ctx, "ProviderCache.setProviderBC").Finish()

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return err
	}

	return p.client.Set(ctx, key, b.Bytes(), p.expiration).Err()
}

func (p *ProviderCache) getProvider(ctx context.Context, key string) (internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderCache.getProvider").Finish()

	data, err := p.client.Get(ctx, key).Bytes()
	if err != nil {
		return internal.Provider{}, err
	}

	var provider internal.Provider
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&provider); err != nil {
		return internal.Provider{}, err
	}

	return provider, nil
}

func (p *ProviderCache) getProviderBC(ctx context.Context, key string) (internal.ProviderBC, error) {
	defer newSentrySpan(ctx, "ProviderCache.getProviderBC").Finish()

	data, err := p.client.Get(ctx, key).Bytes()
	if err != nil {
		return internal.ProviderBC{}, err
	}

	var provider internal.ProviderBC
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&provider); err != nil {
		return internal.ProviderBC{}, err
	}

	return provider, nil
}

func (p *ProviderCache) setManyProviders(ctx context.Context, key string, value []internal.Provider) error {
	defer newSentrySpan(ctx, "ProviderCache.setManyProviders").Finish()

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return err
	}

	return p.client.Set(ctx, key, b.Bytes(), p.expiration).Err()
}

func (p *ProviderCache) getManyProviders(ctx context.Context, key string) ([]internal.Provider, error) {
	defer newSentrySpan(ctx, "ProviderCache.getManyProviders").Finish()

	data, err := p.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var providers []internal.Provider
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&providers); err != nil {
		return nil, err
	}

	return providers, nil
}

func (p *ProviderCache) deleteProvider(ctx context.Context, key string) error {
	defer newSentrySpan(ctx, "ProviderCache.deleteProvider").Finish()

	return p.client.Del(ctx, key).Err()
}
