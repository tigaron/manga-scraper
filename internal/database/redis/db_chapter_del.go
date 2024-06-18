package redis

import (
	"context"
	"fmt"
)

func (c *RedisClient) DeleteAllChapterCacheV1(
	ctx context.Context,
	provider, series, chapter string,
) (
	err error,
) {
	if c.environment == "development" {
		return
	}

	keys, err := c.client.Keys(
		ctx,
		fmt.Sprintf(
			"v1:db:provider:%s:series:%s:chapter:%s:*",
			provider,
			series,
			chapter,
		),
	).Result()
	if err != nil {
		return
	}

	for i := range keys {
		_ = c.client.Del(ctx, keys[i]).Err()
	}

	return
}
