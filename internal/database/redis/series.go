package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
)

func (c *RedisClient) GetSeriesV1(ctx context.Context, provider string, series string) (v1Response.GetSeriesData, error) {
	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s:series:%s", provider, series))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return v1Response.GetSeriesData{}, err
	}

	b := bytes.NewReader(cmdb)

	var res v1Response.GetSeriesData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return v1Response.GetSeriesData{}, err
	}

	return res, nil
}

func (c *RedisClient) SetSeriesV1(ctx context.Context, provider string, series string, s v1Response.GetSeriesData) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(s); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s:series:%s", provider, series), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) UnsetSeriesV1(ctx context.Context, provider string, series string) error {
	return c.client.Del(ctx, fmt.Sprintf("v1:provider:%s:series:%s", provider, series)).Err()
}
