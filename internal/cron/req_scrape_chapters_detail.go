package cron

import (
	"context"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"go.uber.org/zap"
)

func (s *Cron) scrapeChaptersDetail() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	providers, err := s.provider.FindAll(ctx, internal.ASC)
	if err != nil {
		s.logger.Error("Failed to get providers", zap.Error(err))
		return
	}

	for i := range providers {
		receipt, err := s.series.FindEmptyChapters(ctx, internal.FindSeriesParams{
			Provider: providers[i].Slug,
			Order:    internal.ASC,
		})
		if err != nil {
			s.logger.Error("Failed to get chapters", zap.Error(err))
			continue
		}

		for j := range receipt {
			_, err := s.scraper.Create(ctx, receipt[j])
			if err != nil {
				s.logger.Error("Failed to create scrape request", zap.Error(err))
			}
		}
	}
}
