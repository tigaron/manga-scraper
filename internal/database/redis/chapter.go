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

func (c *RedisClient) GetChapterV1(ctx context.Context, provider string, series string, chapter string) (v1Response.ChapterData, error) {
	if c.environment == "development" {
		return v1Response.ChapterData{}, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter:%s", provider, series, chapter))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return v1Response.ChapterData{}, err
	}

	b := bytes.NewReader(cmdb)

	var res v1Response.ChapterData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return v1Response.ChapterData{}, err
	}

	return res, nil
}

func (c *RedisClient) SetChapterV1(ctx context.Context, provider string, series string, chapter string, ch v1Response.ChapterData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter:%s", provider, series, chapter), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) UnsetChapterV1(ctx context.Context, provider string, series string, chapter string) error {
	return c.client.Del(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter:%s", provider, series, chapter)).Err()
}

func (c *RedisClient) FindChapterUniqueV1(ctx context.Context, provider string, series string, chapter string) (*db.ChapterModel, error) {
	if c.environment == "development" {
		return nil, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:db:provider:%s:series:%s:chapter:%s", provider, series, chapter))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(cmdb)

	var res db.ChapterModel

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *RedisClient) CreateChapterUniqueV1(ctx context.Context, provider string, series string, chapter string, ch *db.ChapterModel) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(*ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:db:provider:%s:series:%s:chapter:%s", provider, series, chapter), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) DeleteChapterUniqueV1(ctx context.Context, provider string, series string, chapter string) error {
	if c.environment == "development" {
		return nil
	}

	return c.client.Del(ctx, fmt.Sprintf("v1:db:provider:%s:series:%s:chapter:%s", provider, series, chapter)).Err()
}
