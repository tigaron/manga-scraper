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

func (c *RedisClient) SetChapterV1(ctx context.Context, provider string, series string, chapter string, ch v1Response.ChapterData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter:%s", provider, series, chapter), b.Bytes(), time.Hour).Err()
}

func (c *RedisClient) SetChapterListV1(ctx context.Context, provider string, series string, page int, limit int, ch v1Response.PaginatedChapterListData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:%d:%d", provider, series, page, limit), b.Bytes(), time.Hour).Err()
}

func (c *RedisClient) SetChapterListAllV1(ctx context.Context, provider string, series string, ch []v1Response.ChapterData) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:all", provider, series), b.Bytes(), time.Hour).Err()
}

func (c *RedisClient) SetChapterListOnlyV1(ctx context.Context, provider string, series string, ch v1Response.ListChapterResult) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:only", provider, series), b.Bytes(), time.Hour).Err()
}

func (c *RedisClient) SetChaptersListWithRelV1(
	ctx context.Context,
	provider, series string,
	order db.SortOrder,
	chapter v1Response.ListChapterResult,
) (
	err error,
) {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(chapter); err != nil {
		return err
	}

	return c.client.Set(
		ctx,
		fmt.Sprintf(
			"v1:provider:%s:series:%s:chapter:list:%s:rel",
			provider,
			series,
			order,
		), b.Bytes(),
		time.Hour,
	).Err()
}
