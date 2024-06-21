package redis

import (
	"bytes"
	"context"
	"encoding/gob"

	"fourleaves.studio/manga-scraper/internal"
)

func (s *SeriesCache) setSeries(ctx context.Context, key string, value internal.Series) error {
	defer newSentrySpan(ctx, "SeriesCache.setSeries").Finish()

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return err
	}

	return s.client.Set(ctx, key, b.Bytes(), s.expiration).Err()
}

func (s *SeriesCache) setSeriesBC(ctx context.Context, key string, value internal.SeriesBC) error {
	defer newSentrySpan(ctx, "SeriesCache.setSeriesBC").Finish()

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return err
	}

	return s.client.Set(ctx, key, b.Bytes(), s.expiration).Err()
}

func (s *SeriesCache) getSeries(ctx context.Context, key string) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesCache.getSeries").Finish()

	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return internal.Series{}, err
	}

	var series internal.Series
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&series); err != nil {
		return internal.Series{}, err
	}

	return series, nil
}

func (s *SeriesCache) getSeriesBC(ctx context.Context, key string) (internal.SeriesBC, error) {
	defer newSentrySpan(ctx, "SeriesCache.getSeriesBC").Finish()

	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return internal.SeriesBC{}, err
	}

	var series internal.SeriesBC
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&series); err != nil {
		return internal.SeriesBC{}, err
	}

	return series, nil
}

func (s *SeriesCache) setManySeries(ctx context.Context, key string, value []internal.Series) error {
	defer newSentrySpan(ctx, "SeriesCache.setManySeries").Finish()

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return err
	}

	return s.client.Set(ctx, key, b.Bytes(), s.expiration).Err()
}

func (s *SeriesCache) getManySeries(ctx context.Context, key string) ([]internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesCache.getManySeries").Finish()

	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var series []internal.Series
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&series); err != nil {
		return nil, err
	}

	return series, nil
}

func (s *SeriesCache) deleteSeries(ctx context.Context, key string) error {
	defer newSentrySpan(ctx, "SeriesCache.deleteSeries").Finish()

	return s.client.Del(ctx, key).Err()
}

func (s *SeriesCache) deleteManySeries(ctx context.Context, key string) error {
	defer newSentrySpan(ctx, "SeriesCache.deleteManySeries").Finish()

	keys, err := s.client.Keys(ctx, key).Result()
	if err != nil {
		return err
	}

	return s.client.Del(ctx, keys...).Err()
}
