package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
)

func (c *RedisClient) GetChapterBreadcrumbsV1(ctx context.Context, provider string, series string, chapter string) (v1Response.BreadcrumbsData, error) {
	if c.environment == "development" {
		return v1Response.BreadcrumbsData{}, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter:%s:bc", provider, series, chapter))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return v1Response.BreadcrumbsData{}, err
	}

	b := bytes.NewReader(cmdb)

	var res v1Response.BreadcrumbsData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return v1Response.BreadcrumbsData{}, err
	}

	return res, nil
}

func (c *RedisClient) SetChapterBreadcrumbsV1(ctx context.Context, provider string, series string, chapter string, bc v1Response.BreadcrumbsData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(bc); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter:%s:bc", provider, series, chapter), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) GetSeriesBreadcrumbsV1(ctx context.Context, provider string, series string) (v1Response.BreadcrumbsData, error) {
	if c.environment == "development" {
		return v1Response.BreadcrumbsData{}, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series:%s:bc", provider, series))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return v1Response.BreadcrumbsData{}, err
	}

	b := bytes.NewReader(cmdb)

	var res v1Response.BreadcrumbsData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return v1Response.BreadcrumbsData{}, err
	}

	return res, nil
}

func (c *RedisClient) SetSeriesBreadcrumbsV1(ctx context.Context, provider string, series string, bc v1Response.BreadcrumbsData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(bc); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:bc", provider, series), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) GetProviderBreadcrumbsV1(ctx context.Context, provider string) (v1Response.BreadcrumbsData, error) {
	if c.environment == "development" {
		return v1Response.BreadcrumbsData{}, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:bc", provider))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return v1Response.BreadcrumbsData{}, err
	}

	b := bytes.NewReader(cmdb)

	var res v1Response.BreadcrumbsData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return v1Response.BreadcrumbsData{}, err
	}

	return res, nil
}

func (c *RedisClient) SetProviderBreadcrumbsV1(ctx context.Context, provider string, bc v1Response.BreadcrumbsData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(bc); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:bc", provider), b.Bytes(), 24*time.Hour).Err()
}
