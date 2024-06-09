package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
)

func (c *RedisClient) GetChapterListV1(ctx context.Context, provider string, series string, page int, limit int) (v1Response.PaginatedChapterListData, error) {
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

func (c *RedisClient) SetChapterListV1(ctx context.Context, provider string, series string, page int, limit int, ch v1Response.PaginatedChapterListData) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:%d:%d", provider, series, page, limit), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) GetChapterListAllV1(ctx context.Context, provider string, series string) ([]v1Response.ChapterData, error) {
	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:all", provider, series))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(cmdb)

	var res []v1Response.ChapterData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *RedisClient) SetChapterListAllV1(ctx context.Context, provider string, series string, ch []v1Response.ChapterData) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:all", provider, series), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) GetChapterListOnlyV1(ctx context.Context, provider string, series string) ([]v1Response.ListChapterData, error) {
	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:only", provider, series))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(cmdb)

	var res []v1Response.ListChapterData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *RedisClient) SetChapterListOnlyV1(ctx context.Context, provider string, series string, ch []v1Response.ListChapterData) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:only", provider, series), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) UnsetChapterListV1(ctx context.Context, provider string, series string) error {
	keys, err := c.client.Keys(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:*", provider, series)).Result()
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
