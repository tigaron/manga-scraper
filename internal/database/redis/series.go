package redis

import (
	"context"
	"fmt"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type SeriesStore interface {
	CreateInit(ctx context.Context, params internal.CreateInitSeriesParams) (internal.Series, error)
	Find(ctx context.Context, params internal.FindSeriesParams) (internal.Series, error)
	FindBC(ctx context.Context, params internal.FindSeriesParams) (internal.SeriesBC, error)
	FindAll(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error)
	FindPaginated(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error)
	UpdateInit(ctx context.Context, params internal.UpdateInitSeriesParams) (internal.Series, error)
	UpdateLatest(ctx context.Context, params internal.UpdateLatestSeriesParams) (internal.Series, error)
	Delete(ctx context.Context, params internal.FindSeriesParams) error
}

type SeriesCache struct {
	client     *redis.Client
	store      SeriesStore
	expiration time.Duration
	logger     echo.Logger
}

func NewSeriesCache(redisURL string, store SeriesStore, expiration time.Duration, logger echo.Logger) *SeriesCache {
	opts, _ := redis.ParseURL(redisURL)
	return &SeriesCache{
		client:     redis.NewClient(opts),
		store:      store,
		expiration: expiration,
		logger:     logger,
	}
}

func (s *SeriesCache) CreateInit(ctx context.Context, params internal.CreateInitSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesCache.CreateInit").Finish()

	series, err := s.store.CreateInit(ctx, params)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.CreateInit")
	}

	cacheKey := fmt.Sprintf("v1:series:%s", series.Slug)

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.CreateInit",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   series,
	})

	err = s.setSeries(ctx, cacheKey, series)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "setSeries")
	}

	return series, nil
}

func (s *SeriesCache) Find(ctx context.Context, params internal.FindSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesCache.Find").Finish()

	cacheKey := fmt.Sprintf("v1:series:%s", params.Slug)

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.Find",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	series, err := s.getSeries(ctx, cacheKey)
	if err == nil {
		return series, nil
	}

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.Find",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	series, err = s.store.Find(ctx, params)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.Find")
	}

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.Find",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   series,
	})

	err = s.setSeries(ctx, cacheKey, series)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "setSeries")
	}

	return series, nil
}

func (s *SeriesCache) FindBC(ctx context.Context, params internal.FindSeriesParams) (internal.SeriesBC, error) {
	defer newSentrySpan(ctx, "SeriesCache.FindBC").Finish()

	cacheKey := fmt.Sprintf("v1:series:%s:_bc", params.Slug)

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.FindBC",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	series, err := s.getSeriesBC(ctx, cacheKey)
	if err == nil {
		return series, nil
	}

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.FindBC",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	series, err = s.store.FindBC(ctx, params)
	if err != nil {
		return internal.SeriesBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.FindBC")
	}

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.FindBC",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   series,
	})

	err = s.setSeriesBC(ctx, cacheKey, series)
	if err != nil {
		return internal.SeriesBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "setSeriesBC")
	}

	return series, nil
}

func (s *SeriesCache) FindAll(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesCache.FindAll").Finish()

	cacheKey := fmt.Sprintf("v1:series:%s:_list:%s:all", params.Provider, params.Order)

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.FindAll",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := s.getManySeries(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.FindAll",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	series, err := s.store.FindAll(ctx, params)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "store.FindAll")
	}

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.FindAll",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   series,
	})

	err = s.setManySeries(ctx, cacheKey, series)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "setManySeries")
	}

	return series, nil
}

func (s *SeriesCache) FindPaginated(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesCache.FindPaginated").Finish()

	cacheKey := fmt.Sprintf("v1:series:%s:_list:%s:page:%d:size:%d", params.Provider, params.Order, params.Page, params.Size)

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.FindPaginated",
		"_msg":    "get cache",
		"key":     cacheKey,
	})

	data, err := s.getManySeries(ctx, cacheKey)
	if err == nil {
		return data, nil
	}

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.FindPaginated",
		"_msg":    "cache miss",
		"key":     cacheKey,
	})

	series, err := s.store.FindPaginated(ctx, params)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "store.FindPaginated")
	}

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.FindPaginated",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   series,
	})

	err = s.setManySeries(ctx, cacheKey, series)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "setManySeries")
	}

	return series, nil
}

func (s *SeriesCache) UpdateInit(ctx context.Context, params internal.UpdateInitSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesCache.UpdateInit").Finish()

	series, err := s.store.UpdateInit(ctx, params)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.UpdateInit")
	}

	cacheKey := fmt.Sprintf("v1:series:%s", series.Slug)

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.UpdateInit",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   series,
	})

	err = s.setSeries(ctx, cacheKey, series)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "setSeries")
	}

	_ = s.deleteManySeries(ctx, fmt.Sprintf("v1:series:%s:_list:*", series.Provider))

	return series, nil
}

func (s *SeriesCache) UpdateLatest(ctx context.Context, params internal.UpdateLatestSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesCache.UpdateLatest").Finish()

	series, err := s.store.UpdateLatest(ctx, params)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "store.UpdateLatest")
	}

	cacheKey := fmt.Sprintf("v1:series:%s", series.Slug)

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.UpdateLatest",
		"_msg":    "set cache",
		"key":     cacheKey,
		"value":   series,
	})

	err = s.setSeries(ctx, cacheKey, series)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "setSeries")
	}

	_ = s.deleteManySeries(ctx, fmt.Sprintf("v1:series:%s:_list:*", series.Provider))

	return series, nil
}

func (s *SeriesCache) Delete(ctx context.Context, params internal.FindSeriesParams) error {
	defer newSentrySpan(ctx, "SeriesCache.Delete").Finish()

	cacheKey := fmt.Sprintf("v1:series:%s", params.Slug)

	s.logger.Debugj(map[string]interface{}{
		"_source": "SeriesCache.Delete",
		"_msg":    "delete cache",
		"key":     cacheKey,
	})

	_ = s.deleteSeries(ctx, cacheKey)
	_ = s.deleteManySeries(ctx, fmt.Sprintf("v1:series:%s:_list:*", params.Provider))

	return s.store.Delete(ctx, params)
}
