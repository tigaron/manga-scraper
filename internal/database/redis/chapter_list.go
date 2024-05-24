package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
)

func (c *RedisClient) GetChapterListV1(ctx context.Context, provider string, series string, page int, limit int) ([]v1Response.ChapterData, error) {
	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:%d:%d", provider, series, page, limit))

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

func (c *RedisClient) SetChapterListV1(ctx context.Context, provider string, series string, page int, limit int, ch []v1Response.ChapterData) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:%d:%d", provider, series, page, limit), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) GetAllChapterListV1(ctx context.Context, provider string, series string) ([]v1Response.ChapterData, error) {
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

func (c *RedisClient) SetAllChapterListV1(ctx context.Context, provider string, series string, ch []v1Response.ChapterData) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:all", provider, series), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) UnsetChapterListV1(ctx context.Context, provider string, series string) error {
	return c.client.Del(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter_list:*", provider, series)).Err()
}
