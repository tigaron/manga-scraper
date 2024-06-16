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

func (c *RedisClient) GetSeriesV1(ctx context.Context, provider string, series string) (v1Response.SeriesData, error) {
	if c.environment == "development" {
		return v1Response.SeriesData{}, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series:%s", provider, series))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return v1Response.SeriesData{}, err
	}

	b := bytes.NewReader(cmdb)

	var res v1Response.SeriesData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return v1Response.SeriesData{}, err
	}

	return res, nil
}

func (c *RedisClient) SetSeriesV1(ctx context.Context, provider string, series string, s v1Response.SeriesData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(s); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s", provider, series), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) UnsetSeriesV1(ctx context.Context, provider string, series string) error {
	if c.environment == "development" {
		return nil
	}

	return c.client.Del(ctx, fmt.Sprintf("v1:provider:%s:series:%s", provider, series)).Err()
}

func (c *RedisClient) FindSeriesUniqueV1(ctx context.Context, provider string, series string) (*db.SeriesModel, error) {
	if c.environment == "development" {
		return nil, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:db:provider:%s:series:%s", provider, series))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(cmdb)

	var res db.SeriesModel

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *RedisClient) CreateSeriesUniqueV1(ctx context.Context, s *db.SeriesModel) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(*s); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:db:provider:%s:series:%s", s.ProviderSlug, s.Slug), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) DeleteSeriesUniqueV1(ctx context.Context, provider string, series string) error {
	if c.environment == "development" {
		return nil
	}

	return c.client.Del(ctx, fmt.Sprintf("v1:db:provider:%s:series:%s", provider, series)).Err()
}
