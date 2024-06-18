package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

// One

func (c *RedisClient) CreateChapterUniqueV1(
	ctx context.Context,
	provider, series, chapter string,
	ch *db.ChapterModel,
) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(*ch); err != nil {
		return err
	}

	return c.client.Set(
		ctx,
		fmt.Sprintf(
			"v1:db:provider:%s:series:%s:chapter:%s",
			provider,
			series,
			chapter,
		),
		b.Bytes(),
		time.Hour,
	).Err()
}

func (c *RedisClient) CreateChapterUniqueWithRelV1(
	ctx context.Context,
	provider, series, chapter string,
	ch *db.ChapterModel,
) error {
	if c.environment == "development" {
		return nil
	}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(*ch); err != nil {
		return err
	}

	return c.client.Set(
		ctx,
		fmt.Sprintf(
			"v1:db:provider:%s:series:%s:chapter:%s:rel",
			provider,
			series,
			chapter,
		),
		b.Bytes(),
		time.Hour,
	).Err()
}

// Many

func (c *RedisClient) CreateChaptersListWithRelV1(
	ctx context.Context,
	provider, series string,
	order db.SortOrder,
	chapter *db.SeriesModel,
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
			"v1:db:provider:%s:series:%s:chapter:list:%s:rel",
			provider,
			series,
			order,
		), b.Bytes(),
		time.Hour,
	).Err()
}

func (c *RedisClient) CreateChaptersListAllV1(
	ctx context.Context,
	provider, series string,
	order db.SortOrder,
	chapter *db.SeriesModel,
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
			"v1:db:provider:%s:series:%s:chapter:list:%s:all",
			provider,
			series,
			order,
		), b.Bytes(),
		time.Hour,
	).Err()
}

func (c *RedisClient) CreateChaptersListPaginatedV1(
	ctx context.Context,
	provider, series string,
	page, limit int,
	order db.SortOrder,
	chapter *db.SeriesModel,
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
			"v1:db:provider:%s:series:%s:chapter:list:%s:%d:%d",
			provider,
			series,
			order,
			page,
			limit,
		),
		b.Bytes(),
		time.Hour,
	).Err()
}
