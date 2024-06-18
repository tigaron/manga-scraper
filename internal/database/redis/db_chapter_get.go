package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

// One

func (c *RedisClient) FindChapterUniqueV1(
	ctx context.Context,
	provider, series, chapter string,
) (
	result *db.ChapterModel,
	err error,
) {
	if c.environment == "development" {
		err = fmt.Errorf("not available in development")
		return
	}

	cmd := c.client.Get(
		ctx,
		fmt.Sprintf(
			"v1:db:provider:%s:series:%s:chapter:%s",
			provider,
			series,
			chapter,
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

func (c *RedisClient) FindChapterUniqueWithRelV1(
	ctx context.Context,
	provider, series, chapter string,
) (
	result *db.ChapterModel,
	err error,
) {
	if c.environment == "development" {
		err = fmt.Errorf("not available in development")
		return
	}

	cmd := c.client.Get(
		ctx,
		fmt.Sprintf(
			"v1:db:provider:%s:series:%s:chapter:%s:rel",
			provider,
			series,
			chapter,
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

// Many

func (c *RedisClient) FindChaptersListWithRelV1(
	ctx context.Context,
	provider, series string,
	order db.SortOrder,
) (
	result *db.SeriesModel,
	err error,
) {
	if c.environment == "development" {
		err = fmt.Errorf("not available in development")
		return
	}

	cmd := c.client.Get(
		ctx,
		fmt.Sprintf(
			"v1:db:provider:%s:series:%s:chapter:list:%s:rel",
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

func (c *RedisClient) FindChaptersListAllV1(
	ctx context.Context,
	provider, series string,
	order db.SortOrder,
) (
	result *db.SeriesModel,
	err error,
) {
	if c.environment == "development" {
		err = fmt.Errorf("not available in development")
		return
	}

	cmd := c.client.Get(
		ctx,
		fmt.Sprintf(
			"v1:db:provider:%s:series:%s:chapter:list:%s:all",
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

func (c *RedisClient) FindChaptersListPaginatedV1(
	ctx context.Context,
	provider, series string,
	page, limit int,
	order db.SortOrder,
) (
	result *db.SeriesModel,
	err error,
) {
	if c.environment == "development" {
		err = fmt.Errorf("not available in development")
		return
	}

	cmd := c.client.Get(
		ctx,
		fmt.Sprintf(
			"v1:db:provider:%s:series:%s:chapter:list:%s:%d:%d",
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
