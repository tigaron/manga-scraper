package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
)

func (c *RedisClient) GetSeriesListV1(ctx context.Context, provider string, page int, size int) (v1Response.PaginatedSeriesData, error) {
	if c.environment == "development" {
		return v1Response.PaginatedSeriesData{}, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series_list:%d:%d", provider, page, size))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return v1Response.PaginatedSeriesData{}, err
	}

	b := bytes.NewReader(cmdb)

	var res v1Response.PaginatedSeriesData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return v1Response.PaginatedSeriesData{}, err
	}

	return res, nil
}

func (c *RedisClient) GetSeriesListAllV1(ctx context.Context, provider string) ([]v1Response.SeriesData, error) {
	if c.environment == "development" {
		return nil, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series_list:all", provider))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(cmdb)

	var res []v1Response.SeriesData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *RedisClient) SetSeriesListV1(ctx context.Context, provider string, page int, size int, s v1Response.PaginatedSeriesData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(s); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series_list:%d:%d", provider, page, size), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) SetSeriesListAllV1(ctx context.Context, provider string, s []v1Response.SeriesData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(s); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series_list:all", provider), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) UnsetSeriesListV1(ctx context.Context, provider string) error {
	if c.environment == "development" {
		return nil
	}

	keys, err := c.client.Keys(ctx, fmt.Sprintf("v1:provider:%s:series_list:*", provider)).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := c.client.Del(ctx, key).Err(); err != nil {
			return err
		}
	}

	return nil
}
