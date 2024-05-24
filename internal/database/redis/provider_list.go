package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
)

func (c *RedisClient) GetProviderListV1(ctx context.Context) ([]v1Response.GetProviderData, error) {
	cmd := c.client.Get(ctx, "v1:provider_list")

	cmdb, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(cmdb)

	var res []v1Response.GetProviderData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *RedisClient) SetProviderListV1(ctx context.Context, p []v1Response.GetProviderData) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(p); err != nil {
		return err
	}

	return c.client.Set(ctx, "v1:provider_list", b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) UnsetProviderListV1(ctx context.Context) error {
	return c.client.Del(ctx, "v1:provider_list").Err()
}
