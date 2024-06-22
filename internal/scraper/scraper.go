package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"

	"fourleaves.studio/manga-scraper/internal"
	"fourleaves.studio/manga-scraper/internal/scraper/agscomics"
	"fourleaves.studio/manga-scraper/internal/scraper/anigliscans"
	"fourleaves.studio/manga-scraper/internal/scraper/asura"
	"fourleaves.studio/manga-scraper/internal/scraper/flame"
	"fourleaves.studio/manga-scraper/internal/scraper/luminous"
	"fourleaves.studio/manga-scraper/internal/scraper/mangagalaxy"
	"fourleaves.studio/manga-scraper/internal/scraper/nightscans"
	"fourleaves.studio/manga-scraper/internal/scraper/surya"
)

type SeriesRepository interface {
	UpsertInit(ctx context.Context, params internal.CreateInitSeriesParams) (internal.Series, error)
	Find(ctx context.Context, params internal.FindSeriesParams) (internal.Series, error)
	UpdateInit(ctx context.Context, params internal.UpdateInitSeriesParams) (internal.Series, error)
	UpdateLatest(ctx context.Context, params internal.UpdateLatestSeriesParams) (internal.Series, error)
}

type SeriesSearchRepository interface {
	Index(ctx context.Context, series internal.Series) error
}

type ChapterRepository interface {
	UpsertInit(ctx context.Context, params internal.CreateInitChapterParams) (internal.Chapter, error)
	Find(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error)
	FindLatest(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error)
	Count(ctx context.Context, params internal.FindChapterParams) (int, error)
	UpdateInit(ctx context.Context, params internal.UpdateInitChapterParams) (internal.Chapter, error)
}

type ScrapeRequestRepository interface {
	Find(ctx context.Context, id string) (internal.ScrapeRequest, error)
	FindPendings(ctx context.Context, params internal.FindScrapeRequestParams) ([]internal.ScrapeRequest, error)
	Update(ctx context.Context, params internal.UpdateScrapeRequestParams) (internal.ScrapeRequest, error)
}

type Scraper struct {
	repo        ScrapeRequestRepository
	series      SeriesRepository
	search      SeriesSearchRepository
	chapter     ChapterRepository
	kafkaClient *kafka.Consumer
	logger      *zap.Logger
	browserURL  string
	doneC       chan struct{}
	closeC      chan struct{}
}

func NewScraper(
	repo ScrapeRequestRepository,
	series SeriesRepository,
	search SeriesSearchRepository,
	chapter ChapterRepository,
	kafkaClient *kafka.Consumer,
	logger *zap.Logger,
	browserURL string,
) *Scraper {
	return &Scraper{
		repo:        repo,
		series:      series,
		search:      search,
		chapter:     chapter,
		kafkaClient: kafkaClient,
		logger:      logger,
		browserURL:  browserURL,
		doneC:       make(chan struct{}),
		closeC:      make(chan struct{}),
	}
}

func (s *Scraper) StartServer() (<-chan error, error) {
	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()

		s.logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			_ = s.logger.Sync()
			_ = s.kafkaClient.Unsubscribe()

			stop()
			cancel()
			close(errC)
		}()

		if err := s.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		s.logger.Info("Shutdown completed")
	}()

	go func() {
		s.logger.Info("Listening and serving")

		if err := s.ListenAndServe(); err != nil {
			errC <- err
		}
	}()

	return errC, nil
}

func (s *Scraper) ListenAndServe() error {
	commit := func(msg *kafka.Message) {
		if _, err := s.kafkaClient.CommitMessage(msg); err != nil {
			s.logger.Error("commit failed", zap.Error(err))
		}
	}

	go func() {
		run := true

		for run {
			select {
			case <-s.closeC:
				run = false
			default:
				msg, ok := s.kafkaClient.Poll(150).(*kafka.Message)
				if !ok {
					continue
				}

				var evt struct {
					Type  string
					Value internal.ScrapeRequest
				}

				if err := json.NewDecoder(bytes.NewReader(msg.Value)).Decode(&evt); err != nil {
					s.logger.Info("Ignoring message, invalid", zap.Error(err))
					commit(msg)

					continue
				}

				maxRetries := 3
				switch evt.Type {
				case string(internal.SeriesListRequestType):
					for attempt := 1; attempt <= maxRetries; attempt++ {
						if err := s.ScrapeSeriesList(context.Background(), evt.Value); err != nil {
							s.logger.Error("Error processing message, retrying", zap.Error(err), zap.Int("attempt", attempt))
							time.Sleep(time.Duration(attempt) * 5 * time.Second)
							continue
						}
						break
					}
				case string(internal.SeriesDetailRequestType):
					for attempt := 1; attempt <= maxRetries; attempt++ {
						if err := s.ScrapeSeriesDetail(context.Background(), evt.Value); err != nil {
							s.logger.Error("Error processing message, retrying", zap.Error(err), zap.Int("attempt", attempt))
							time.Sleep(time.Duration(attempt) * 5 * time.Second)
							continue
						}
						break
					}
				case string(internal.ChapterListRequestType):
					for attempt := 1; attempt <= maxRetries; attempt++ {
						if err := s.ScrapeChapterList(context.Background(), evt.Value); err != nil {
							s.logger.Error("Error processing message, retrying", zap.Error(err), zap.Int("attempt", attempt))
							time.Sleep(time.Duration(attempt) * 5 * time.Second)
							continue
						}
						break
					}
				case string(internal.ChapterDetailRequestType):
					for attempt := 1; attempt <= maxRetries; attempt++ {
						if err := s.ScrapeChapterDetail(context.Background(), evt.Value); err != nil {
							s.logger.Error("Error processing message, retrying", zap.Error(err), zap.Int("attempt", attempt))
							time.Sleep(time.Duration(attempt) * 5 * time.Second)
							continue
						}
						break
					}
				}

				s.logger.Info("Consumed", zap.String("type", evt.Type))
				commit(msg)
			}
		}

		s.logger.Info("No more messages to consume. Exiting.")

		s.doneC <- struct{}{}
	}()

	return nil
}

func (s *Scraper) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")

	close(s.closeC)

	for {
		select {
		case <-ctx.Done():
			return internal.WrapErrorf(ctx.Err(), internal.ErrUnknown, "context.Done")
		case <-s.doneC:
			return nil
		}
	}
}

func (s *Scraper) ScrapeSeriesList(ctx context.Context, event internal.ScrapeRequest) error {
	var result []internal.SeriesListResult
	var err error

	s.logger.Info(
		"received scrape request",
		zap.String("id", event.ID),
		zap.String("type", string(event.Type)),
		zap.String("provider", event.Provider),
		zap.String("baseURL", event.BaseURL),
		zap.String("requestPath", event.RequestPath),
	)

	requestURL := event.BaseURL + event.RequestPath

	startTime := time.Now()

	switch event.Provider {
	case "asura":
		result, err = asura.ScrapeSeriesList(ctx, s.browserURL, requestURL)
	case "surya":
		result, err = surya.ScrapeSeriesList(ctx, s.browserURL, requestURL)
	case "flame":
		result, err = flame.ScrapeSeriesList(ctx, s.browserURL, requestURL)
	case "luminous":
		result, err = luminous.ScrapeSeriesList(ctx, s.browserURL, requestURL)
	case "anigliscans":
		result, err = anigliscans.ScrapeSeriesList(ctx, s.browserURL, requestURL)
	case "agscomics":
		result, err = agscomics.ScrapeSeriesList(ctx, s.browserURL, requestURL)
	case "nightscans":
		result, err = nightscans.ScrapeSeriesList(ctx, s.browserURL, requestURL)
	case "mangagalaxy":
		result, err = mangagalaxy.ScrapeSeriesList(ctx, s.browserURL, requestURL)
	default:
		err = internal.NewErrorf(internal.ErrInvalidInput, "not implemented yet")
	}

	endTime := time.Since(startTime).Seconds()

	if err != nil {
		_, _ = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
			ID:        event.ID,
			Status:    internal.FailedRequestStatus,
			TotalTime: endTime,
			Error:     true,
			Message:   err.Error(),
		})

		return err
	}

	var wg sync.WaitGroup

	wg.Add(len(result))

	for i := range result {
		go func(i int) {
			defer wg.Done()

			series, err := s.series.UpsertInit(ctx, internal.CreateInitSeriesParams{
				Provider:   event.Provider,
				Slug:       result[i].Slug,
				Title:      result[i].Title,
				SourcePath: result[i].SourcePath,
			})
			if err != nil {
				s.logger.Error("failed to create series", zap.Error(err))
				return
			}

			err = s.search.Index(ctx, series)
			if err != nil {
				s.logger.Error("failed to index series", zap.Error(err))
				return
			}
		}(i)
	}

	wg.Wait()

	_, err = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
		ID:        event.ID,
		Status:    internal.CompletedRequestStatus,
		TotalTime: endTime,
		Error:     false,
		Message:   "Completed successfully",
	})

	return err
}

func (s *Scraper) ScrapeSeriesDetail(ctx context.Context, event internal.ScrapeRequest) error {
	var result internal.SeriesDetailResult
	var err error

	s.logger.Info(
		"received scrape request",
		zap.String("id", event.ID),
		zap.String("type", string(event.Type)),
		zap.String("provider", event.Provider),
		zap.String("series", event.Series),
		zap.String("baseURL", event.BaseURL),
		zap.String("requestPath", event.RequestPath),
	)

	requestURL := event.BaseURL + event.RequestPath

	startTime := time.Now()

	switch event.Provider {
	case "asura":
		result, err = asura.ScrapeSeriesDetail(ctx, s.browserURL, requestURL)
	case "surya":
		result, err = surya.ScrapeSeriesDetail(ctx, s.browserURL, requestURL)
	case "flame":
		result, err = flame.ScrapeSeriesDetail(ctx, s.browserURL, requestURL)
	case "luminous":
		result, err = luminous.ScrapeSeriesDetail(ctx, s.browserURL, requestURL)
	case "anigliscans":
		result, err = anigliscans.ScrapeSeriesDetail(ctx, s.browserURL, requestURL)
	case "agscomics":
		result, err = agscomics.ScrapeSeriesDetail(ctx, s.browserURL, requestURL)
	case "nightscans":
		result, err = nightscans.ScrapeSeriesDetail(ctx, s.browserURL, requestURL)
	case "mangagalaxy":
		result, err = mangagalaxy.ScrapeSeriesDetail(ctx, s.browserURL, requestURL)
	default:
		err = internal.NewErrorf(internal.ErrInvalidInput, "not implemented yet")
	}

	endTime := time.Since(startTime).Seconds()

	if err != nil {
		_, _ = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
			ID:        event.ID,
			Status:    internal.FailedRequestStatus,
			TotalTime: endTime,
			Error:     true,
			Message:   err.Error(),
		})

		return err
	}

	series, err := s.series.UpdateInit(ctx, internal.UpdateInitSeriesParams{
		Provider:     event.Provider,
		Slug:         event.Series,
		ThumbnailURL: result.ThumbnailURL,
		Synopsis:     result.Synopsis,
		Genres:       result.Genres,
	})

	if err != nil {
		_, _ = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
			ID:        event.ID,
			Status:    internal.FailedRequestStatus,
			TotalTime: endTime,
			Error:     true,
			Message:   err.Error(),
		})

		return err
	}

	if err := s.search.Index(ctx, series); err != nil {
		_, _ = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
			ID:        event.ID,
			Status:    internal.FailedRequestStatus,
			TotalTime: endTime,
			Error:     true,
			Message:   err.Error(),
		})

		return err
	}

	_, err = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
		ID:        event.ID,
		Status:    internal.CompletedRequestStatus,
		TotalTime: endTime,
		Error:     false,
		Message:   "Completed successfully",
	})

	return err
}

func (s *Scraper) ScrapeChapterList(ctx context.Context, event internal.ScrapeRequest) error {
	var result []internal.ChapterListResult
	var err error

	s.logger.Info(
		"received scrape request",
		zap.String("id", event.ID),
		zap.String("type", string(event.Type)),
		zap.String("provider", event.Provider),
		zap.String("series", event.Series),
		zap.String("baseURL", event.BaseURL),
		zap.String("requestPath", event.RequestPath),
	)

	requestURL := event.BaseURL + event.RequestPath

	startTime := time.Now()

	switch event.Provider {
	case "asura":
		result, err = asura.ScrapeChapterList(ctx, s.browserURL, requestURL)
	case "surya":
		result, err = surya.ScrapeChapterList(ctx, s.browserURL, requestURL)
	case "flame":
		result, err = flame.ScrapeChapterList(ctx, s.browserURL, requestURL)
	case "luminous":
		result, err = luminous.ScrapeChapterList(ctx, s.browserURL, requestURL)
	case "anigliscans":
		result, err = anigliscans.ScrapeChapterList(ctx, s.browserURL, requestURL)
	case "agscomics":
		result, err = agscomics.ScrapeChapterList(ctx, s.browserURL, requestURL)
	case "nightscans":
		result, err = nightscans.ScrapeChapterList(ctx, s.browserURL, requestURL)
	case "mangagalaxy":
		result, err = mangagalaxy.ScrapeChapterList(ctx, s.browserURL, requestURL)
	default:
		err = internal.NewErrorf(internal.ErrInvalidInput, "not implemented yet")
	}

	endTime := time.Since(startTime).Seconds()

	if err != nil {
		_, _ = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
			ID:        event.ID,
			Status:    internal.FailedRequestStatus,
			TotalTime: endTime,
			Error:     true,
			Message:   err.Error(),
		})

		return err
	}

	var wg sync.WaitGroup

	wg.Add(len(result))

	for i := range result {
		go func(i int) {
			defer wg.Done()

			_, err := s.chapter.UpsertInit(ctx, internal.CreateInitChapterParams{
				Provider:   event.Provider,
				Series:     event.Series,
				Slug:       result[i].Slug,
				Number:     result[i].Number,
				ShortTitle: result[i].ShortTitle,
				SourceHref: result[i].Href,
			})
			if err != nil {
				s.logger.Error("failed to create chapter", zap.Error(err))
				return
			}
		}(i)
	}

	wg.Wait()

	// TODO: handle series count and latest chapter

	_, err = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
		ID:        event.ID,
		Status:    internal.CompletedRequestStatus,
		TotalTime: endTime,
		Error:     false,
		Message:   "Completed successfully",
	})

	return err
}

func (s *Scraper) ScrapeChapterDetail(ctx context.Context, event internal.ScrapeRequest) error {
	var result internal.ChapterDetailResult
	var err error

	s.logger.Info(
		"received scrape request",
		zap.String("id", event.ID),
		zap.String("type", string(event.Type)),
		zap.String("provider", event.Provider),
		zap.String("series", event.Series),
		zap.String("chapter", event.Chapter),
		zap.String("baseURL", event.BaseURL),
		zap.String("requestPath", event.RequestPath),
	)

	requestURL := event.BaseURL + event.RequestPath

	startTime := time.Now()

	switch event.Provider {
	case "asura":
		result, err = asura.ScrapeChapterDetail(ctx, s.browserURL, requestURL)
	case "surya":
		result, err = surya.ScrapeChapterDetail(ctx, s.browserURL, requestURL)
	case "flame":
		result, err = flame.ScrapeChapterDetail(ctx, s.browserURL, requestURL)
	case "luminous":
		result, err = luminous.ScrapeChapterDetail(ctx, s.browserURL, requestURL)
	case "anigliscans":
		result, err = anigliscans.ScrapeChapterDetail(ctx, s.browserURL, requestURL)
	case "agscomics":
		result, err = agscomics.ScrapeChapterDetail(ctx, s.browserURL, requestURL)
	case "nightscans":
		result, err = nightscans.ScrapeChapterDetail(ctx, s.browserURL, requestURL)
	case "mangagalaxy":
		result, err = mangagalaxy.ScrapeChapterDetail(ctx, s.browserURL, requestURL)
	default:
		err = internal.NewErrorf(internal.ErrInvalidInput, "not implemented yet")
	}

	endTime := time.Since(startTime).Seconds()

	if err != nil {
		_, _ = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
			ID:        event.ID,
			Status:    internal.FailedRequestStatus,
			TotalTime: endTime,
			Error:     true,
			Message:   err.Error(),
		})

		return err
	}

	_, err = s.chapter.UpdateInit(ctx, internal.UpdateInitChapterParams{
		Provider:     event.Provider,
		Series:       event.Series,
		Slug:         event.Chapter,
		FullTitle:    result.FullTitle,
		SourcePath:   result.SourcePath,
		ContentPaths: result.ContentPaths,
		NextSlug:     result.NextSlug,
		NextPath:     result.NextPath,
		PrevSlug:     result.PrevSlug,
		PrevPath:     result.PrevPath,
	})

	if err != nil {
		_, _ = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
			ID:        event.ID,
			Status:    internal.FailedRequestStatus,
			TotalTime: endTime,
			Error:     true,
			Message:   err.Error(),
		})

		return err
	}

	_, err = s.repo.Update(ctx, internal.UpdateScrapeRequestParams{
		ID:        event.ID,
		Status:    internal.CompletedRequestStatus,
		TotalTime: endTime,
		Error:     false,
		Message:   "Completed successfully",
	})

	return err
}
