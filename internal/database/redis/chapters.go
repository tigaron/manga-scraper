package redis

import (
	"context"
	"fmt"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type ChapterStore interface {
	CreateInit(ctx context.Context, params internal.CreateInitChapterParams) (internal.Chapter, error)
	Find(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error)
	FindBC(ctx context.Context, params internal.FindChapterParams) (internal.ChapterBC, error)
	FindLatest(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error)
	Count(ctx context.Context, params internal.FindChapterParams) (int, error)
	FindAll(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error)
	FindPaginated(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error)
	UpdateInit(ctx context.Context, params internal.UpdateInitChapterParams) (internal.Chapter, error)
	Delete(ctx context.Context, params internal.FindChapterParams) error
}

type ChapterCache struct {
	client     *redis.Client
	store      ChapterStore
	expiration time.Duration
	logger     echo.Logger
}

func NewChapterCache(redisURL string, store ChapterStore, expiration time.Duration, logger echo.Logger) *ChapterCache {
	opts, _ := redis.ParseURL(redisURL)
	return &ChapterCache{
		client:     redis.NewClient(opts),
		store:      store,
		expiration: expiration,
		logger:     logger,
	}
}

func (c *ChapterCache) CreateInit(ctx context.Context, params internal.CreateInitChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterCache.CreateInit").Finish()

	chapter, err := c.store.CreateInit(ctx, params)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.CreateInit")
	}

	cacheKey := fmt.Sprintf("v1:chapters:%s:%s:%s", params.Provider, params.Series, params.Slug)

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.CreateInit",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   chapter,
	})

	err = c.setChapter(ctx, cacheKey, chapter)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "setChapter")
	}

	return chapter, nil
}

func (c *ChapterCache) Find(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterCache.Find").Finish()

	cacheKey := fmt.Sprintf("v1:chapters:%s:%s:%s", params.Provider, params.Series, params.Slug)

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.Find",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := c.getChapter(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.Find",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	chapter, err := c.store.Find(ctx, params)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.Find")
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.Find",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   chapter,
	})

	err = c.setChapter(ctx, cacheKey, chapter)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "setChapter")
	}

	return chapter, nil
}

func (c *ChapterCache) FindBC(ctx context.Context, params internal.FindChapterParams) (internal.ChapterBC, error) {
	defer newSentrySpan(ctx, "ChapterCache.FindBC").Finish()

	cacheKey := fmt.Sprintf("v1:chapters:%s:%s:%s:_bc", params.Provider, params.Series, params.Slug)

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindBC",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := c.getChapterBC(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindBC",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	chapter, err := c.store.FindBC(ctx, params)
	if err != nil {
		return internal.ChapterBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.FindBC")
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindBC",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   chapter,
	})

	err = c.setChapterBC(ctx, cacheKey, chapter)
	if err != nil {
		return internal.ChapterBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "setChapterBC")
	}

	return chapter, nil
}

func (c *ChapterCache) FindLatest(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterCache.FindLatest").Finish()

	cacheKey := fmt.Sprintf("v1:chapters:%s:%s:_latest", params.Provider, params.Series)

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindLatest",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := c.getChapter(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindLatest",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	chapter, err := c.store.FindLatest(ctx, params)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.FindLatest")
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindLatest",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   chapter,
	})

	err = c.setChapter(ctx, cacheKey, chapter)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "setChapter")
	}

	return chapter, nil
}

func (c *ChapterCache) Count(ctx context.Context, params internal.FindChapterParams) (int, error) {
	defer newSentrySpan(ctx, "ChapterCache.Count").Finish()

	cacheKey := fmt.Sprintf("v1:chapters:%s:%s:_count", params.Provider, params.Series)

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.Count",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := c.getChaptersCount(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.Count",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	count, err := c.store.Count(ctx, params)
	if err != nil {
		return 0, internal.WrapErrorf(err, internal.ErrUnknown, "store.Count")
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.Count",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   count,
	})

	err = c.setChaptersCount(ctx, cacheKey, count)
	if err != nil {
		return 0, internal.WrapErrorf(err, internal.ErrUnknown, "setChaptersCount")
	}

	return count, nil
}

func (c *ChapterCache) FindAll(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterCache.FindAll").Finish()

	cacheKey := fmt.Sprintf("v1:chapters:%s:%s:_list:%s:all", params.Provider, params.Series, params.Order)

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindAll",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := c.getManyChapters(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindAll",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	chapters, err := c.store.FindAll(ctx, params)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "store.FindAll")
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindAll",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   chapters,
	})

	err = c.setManyChapters(ctx, cacheKey, chapters)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "setManyChapters")
	}

	return chapters, nil
}

func (c *ChapterCache) FindPaginated(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterCache.FindPaginated").Finish()

	cacheKey := fmt.Sprintf("v1:chapters:%s:%s:_list:%s:page:%d:size:%d", params.Provider, params.Series, params.Order, params.Page, params.Size)

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindPaginated",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := c.getManyChapters(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindPaginated",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	chapters, err := c.store.FindPaginated(ctx, params)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "store.FindPaginated")
	}

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.FindPaginated",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   chapters,
	})

	err = c.setManyChapters(ctx, cacheKey, chapters)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "setManyChapters")
	}

	return chapters, nil
}

func (c *ChapterCache) UpdateInit(ctx context.Context, params internal.UpdateInitChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterCache.UpdateInit").Finish()

	chapter, err := c.store.UpdateInit(ctx, params)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.UpdateInit")
	}

	cacheKey := fmt.Sprintf("v1:chapters:%s:%s:%s", params.Provider, params.Series, params.Slug)

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.UpdateInit",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   chapter,
	})

	err = c.setChapter(ctx, cacheKey, chapter)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "setChapter")
	}

	_ = c.deleteManyChapters(ctx, fmt.Sprintf("v1:chapters:%s:%s:_list:*", params.Provider, params.Series))

	return chapter, nil
}

func (c *ChapterCache) Delete(ctx context.Context, params internal.FindChapterParams) error {
	defer newSentrySpan(ctx, "ChapterCache.Delete").Finish()

	cacheKey := fmt.Sprintf("v1:chapters:%s:%s:%s", params.Provider, params.Series, params.Slug)

	c.logger.Debugj(map[string]interface{}{
		"_source": "ChapterCache.Delete",
		"_msg":    "delete cache",
		"key":     cacheKey,
	})

	err := c.deleteChapter(ctx, cacheKey)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "deleteChapter")
	}

	_ = c.deleteManyChapters(ctx, fmt.Sprintf("v1:chapters:%s:%s:_list:*", params.Provider, params.Series))

	return nil
}
