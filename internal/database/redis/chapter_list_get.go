package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

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

func (c *RedisClient) GetChapterListV1(ctx context.Context, provider string, series string, page int, limit int) (v1Response.PaginatedChapterListData, error) {
	if c.environment == "development" {
		return v1Response.PaginatedChapterListData{}, fmt.Errorf("not available in development")
	}

	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:%d:%d", provider, series, page, limit))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return v1Response.PaginatedChapterListData{}, err
	}

	b := bytes.NewReader(cmdb)

	var res v1Response.PaginatedChapterListData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return v1Response.PaginatedChapterListData{}, err
	}

	return res, nil
}

func (c *RedisClient) GetChaptersListWithRelV1(
	ctx context.Context,
	provider, series string,
	order db.SortOrder,
) (
	result v1Response.ListChapterResult,
	err error,
) {
	if c.environment == "development" {
		err = fmt.Errorf("not available in development")
		return
	}

	cmd := c.client.Get(
		ctx,
		fmt.Sprintf(
			"v1:provider:%s:series:%s:chapter:list:%s:rel",
			provider,
			series,
			order,
		),
	)

	cmdb, err := cmd.Bytes()
	if err != nil {
		return
	}

	b := bytes.NewReader(cmdb)
	err = gob.NewDecoder(b).Decode(&result)
	return
}

func (c *RedisClient) GetChaptersListPaginatedV1(
	ctx context.Context,
	provider, series string,
	page, limit int,
	order db.SortOrder,
) (
	result v1Response.PaginatedChapterListData,
	err error,
) {
	if c.environment == "development" {
		err = fmt.Errorf("not available in development")
		return
	}

	cmd := c.client.Get(
		ctx,
		fmt.Sprintf(
			"v1:provider:%s:series:%s:chapter:list:%s:%d:%d",
			provider,
			series,
			order,
			page,
			limit,
		),
	)

	cmdb, err := cmd.Bytes()
	if err != nil {
		return
	}

	b := bytes.NewReader(cmdb)
	err = gob.NewDecoder(b).Decode(&result)
	return
}

func (c *RedisClient) GetChaptersListAllV1(
	ctx context.Context,
	provider, series string,
	order db.SortOrder,
) (
	result []v1Response.ChapterData,
	err error,
) {
	if c.environment == "development" {
		err = fmt.Errorf("not available in development")
		return
	}

	cmd := c.client.Get(
		ctx,
		fmt.Sprintf(
			"v1:provider:%s:series:%s:chapter:list:%s:all",
			provider,
			series,
			order,
		),
	)

	cmdb, err := cmd.Bytes()
	if err != nil {
		return
	}

	b := bytes.NewReader(cmdb)
	err = gob.NewDecoder(b).Decode(&result)
	return
}
