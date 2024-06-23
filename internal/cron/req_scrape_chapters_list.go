package cron

import (
	"context"
	"strings"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"go.uber.org/zap"
)

func (s *Cron) scrapeChaptersList() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	providers, err := s.provider.FindAll(ctx, internal.ASC)
	if err != nil {
		s.logger.Error("Failed to get providers", zap.Error(err))
		return
	}

	for i := range providers {
		series, err := s.series.FindOnGoing(ctx, internal.FindSeriesParams{
			Provider: providers[i].Slug,
			Order:    internal.ASC,
		})
		if err != nil {
			s.logger.Error("Failed to get series", zap.Error(err))
			continue
		}

		for j := range series {
			params := internal.CreateScrapeRequestParams{
				Type:        internal.ChapterListRequestType,
				Status:      internal.PendingRequestStatus,
				BaseURL:     providers[i].BaseURL,
				RequestPath: strings.Replace(series[j].SourceURL, providers[i].BaseURL, "", 1),
				Provider:    providers[i].Slug,
				Series:      series[j].Slug,
			}

			_, err := s.scraper.Create(ctx, params)
			if err != nil {
				s.logger.Error("Failed to create scrape request", zap.Error(err))
			}
		}
	}
}
