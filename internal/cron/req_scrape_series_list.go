package cron

import (
	"context"
	"strings"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"go.uber.org/zap"
)

func (s *Cron) scrapeSeriesList() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	providers, err := s.provider.FindAll(ctx, internal.ASC)
	if err != nil {
		s.logger.Error("Failed to get providers", zap.Error(err))
		return
	}

	for i := range providers {
		params := internal.CreateScrapeRequestParams{
			Type:        internal.SeriesListRequestType,
			Status:      internal.PendingRequestStatus,
			BaseURL:     providers[i].BaseURL,
			RequestPath: strings.Replace(providers[i].ListURL, providers[i].BaseURL, "", 1),
			Provider:    providers[i].Slug,
		}

		_, err := s.scraper.Create(ctx, params)
		if err != nil {
			s.logger.Error("Failed to create scrape request", zap.Error(err))
		}
	}
}
