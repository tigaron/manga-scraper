package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
)

func (c *RedisClient) GetChapterV1(ctx context.Context, provider string, series string, chapter string) (v1Response.ChapterData, error) {
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
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ch); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter:%s", provider, series, chapter), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) UnsetChapterV1(ctx context.Context, provider string, series string, chapter string) error {
	return c.client.Del(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter:%s", provider, series, chapter)).Err()
}
