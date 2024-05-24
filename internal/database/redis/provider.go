package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
)

func (c *RedisClient) GetProviderV1(ctx context.Context, provider string) (v1Response.ProviderData, error) {
	cmd := c.client.Get(ctx, fmt.Sprintf("v1:provider:%s", provider))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return v1Response.ProviderData{}, err
	}

	b := bytes.NewReader(cmdb)

	var res v1Response.ProviderData

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return v1Response.ProviderData{}, err
	}

	return res, nil
}

func (c *RedisClient) SetProviderV1(ctx context.Context, provider string, p v1Response.ProviderData) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(p); err != nil {
		return err
	}

	return c.client.Set(ctx, fmt.Sprintf("v1:provider:%s", provider), b.Bytes(), 24*time.Hour).Err()
}

func (c *RedisClient) UnsetProviderV1(ctx context.Context, provider string) error {
	return c.client.Del(ctx, fmt.Sprintf("v1:provider:%s", provider)).Err()
}
