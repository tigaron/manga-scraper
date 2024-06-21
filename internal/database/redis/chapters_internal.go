package redis

import (
	"bytes"
	"context"
	"encoding/gob"

	"fourleaves.studio/manga-scraper/internal"
)

func (c *ChapterCache) setChapter(ctx context.Context, key string, value internal.Chapter) error {
	defer newSentrySpan(ctx, "ChapterCache.setChapter").Finish()

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return err
	}

	return c.client.Set(ctx, key, b.Bytes(), c.expiration).Err()
}

func (c *ChapterCache) setChapterBC(ctx context.Context, key string, value internal.ChapterBC) error {
	defer newSentrySpan(ctx, "ChapterCache.setChapterBC").Finish()

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return err
	}

	return c.client.Set(ctx, key, b.Bytes(), c.expiration).Err()
}

func (c *ChapterCache) getChapter(ctx context.Context, key string) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterCache.getChapter").Finish()

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return internal.Chapter{}, err
	}

	var chapter internal.Chapter
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&chapter); err != nil {
		return internal.Chapter{}, err
	}

	return chapter, nil
}

func (c *ChapterCache) getChapterBC(ctx context.Context, key string) (internal.ChapterBC, error) {
	defer newSentrySpan(ctx, "ChapterCache.getChapterBC").Finish()

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return internal.ChapterBC{}, err
	}

	var chapter internal.ChapterBC
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&chapter); err != nil {
		return internal.ChapterBC{}, err
	}

	return chapter, nil
}

func (c *ChapterCache) setManyChapters(ctx context.Context, key string, value []internal.Chapter) error {
	defer newSentrySpan(ctx, "ChapterCache.setManyChapters").Finish()

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return err
	}

	return c.client.Set(ctx, key, b.Bytes(), c.expiration).Err()
}

func (c *ChapterCache) getManyChapters(ctx context.Context, key string) ([]internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterCache.getManyChapters").Finish()

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var chapters []internal.Chapter
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&chapters); err != nil {
		return nil, err
	}

	return chapters, nil
}

func (c *ChapterCache) getChaptersCount(ctx context.Context, key string) (int, error) {
	defer newSentrySpan(ctx, "ChapterCache.getChaptersCount").Finish()

	data, err := c.client.Get(ctx, key).Int()
	if err != nil {
		return 0, err
	}

	return data, nil
}

func (c *ChapterCache) setChaptersCount(ctx context.Context, key string, value int) error {
	defer newSentrySpan(ctx, "ChapterCache.setChaptersCount").Finish()

	return c.client.Set(ctx, key, value, c.expiration).Err()
}

func (c *ChapterCache) deleteChapter(ctx context.Context, key string) error {
	defer newSentrySpan(ctx, "ChapterCache.deleteChapter").Finish()

	return c.client.Del(ctx, key).Err()
}

func (c *ChapterCache) deleteManyChapters(ctx context.Context, key string) error {
	defer newSentrySpan(ctx, "ChapterCache.deleteManyChapters").Finish()

	keys, err := c.client.Keys(ctx, key).Result()
	if err != nil {
		return err
	}

	return c.client.Del(ctx, keys...).Err()
}
