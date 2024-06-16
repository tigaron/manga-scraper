package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

func (c *RedisClient) GetProviderV1(ctx context.Context, provider string) (v1Response.ProviderData, error) {
	if c.environment == "development" {
		return v1Response.ProviderData{}, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s", provider))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return v1Response.ProviderData{}, err
	}

	b := bytes.NewReader(cmdb)

	var res v1Response.ProviderData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return v1Response.ProviderData{}, err
	}

	return res, nil
}

func (c *RedisClient) SetProviderV1(ctx context.Context, provider string, p v1Response.ProviderData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(p); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s", provider), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) UnsetProviderV1(ctx context.Context, provider string) error {
	return c.client.Del(ctx, fmt.Sprintf("v1:provider:%s", provider)).Err()
}

func (c *RedisClient) FindProviderUniqueV1(ctx context.Context, providerSlug string) (*db.ProviderModel, error) {
	if c.environment == "development" {
		return nil, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:db:provider:%s", providerSlug))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(cmdb)

	var res db.ProviderModel

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *RedisClient) CreateProviderUniqueV1(ctx context.Context, provider *db.ProviderModel) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(*provider); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:db:provider:%s", provider.Slug), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) DeleteProviderUniqueV1(ctx context.Context, providerSlug string) error {
	if c.environment == "development" {
		return nil
	}

	return c.client.Del(ctx, fmt.Sprintf("v1:db:provider:%s", providerSlug)).Err()
}
