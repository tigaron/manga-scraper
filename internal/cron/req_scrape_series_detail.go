package cron

import (
	"context"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"go.uber.org/zap"
)

func (s *Cron) scrapeSeriesDetail() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	series, err := s.series.FindEmptyThumb(ctx, internal.ASC)
	if err != nil {
		s.logger.Error("Failed to get series", zap.Error(err))
		return
	}

	for i := range series {
		_, err := s.scraper.Create(ctx, series[i])
		if err != nil {
			s.logger.Error("Failed to create scrape request", zap.Error(err))
		}
	}
}
