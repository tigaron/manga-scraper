package redis

import (
	"context"
	"fmt"
)

func (c *RedisClient) UnsetChapterV1(ctx context.Context, provider string, series string, chapter string) error {
	if c.environment == "development" {
		return nil
	}

	return c.client.Del(ctx, fmt.Sprintf("v1:provider:%s:series:%s:chapter:%s", provider, series, chapter)).Err()
}

func (c *RedisClient) UnsetChapterListV1(ctx context.Context, provider string, series string) error {
	if c.environment == "development" {
		return nil
	}

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
