package cron

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"

	"fourleaves.studio/manga-scraper/internal"
)

type JobRepository interface {
	Create(ctx context.Context, params internal.CreateCronJobParams) (internal.CronJob, error)
	Upsert(ctx context.Context, params internal.CreateCronJobParams) (internal.CronJob, error)
	Find(ctx context.Context, id string) (internal.CronJob, error)
	FindAll(ctx context.Context) ([]internal.CronJob, error)
	CreateStatus(ctx context.Context, params internal.CreateCronJobStatusParams) (internal.CronJobStatus, error)
	UpdateStatus(ctx context.Context, params internal.UpdateCronJobStatusParams) (internal.CronJobStatus, error)
	Delete(ctx context.Context, id string) error
}

type ProviderRepository interface {
	Find(ctx context.Context, slug string) (internal.Provider, error)
	FindAll(ctx context.Context, order internal.SortOrder) ([]internal.Provider, error)
}

type SeriesRepository interface {
	Find(ctx context.Context, params internal.FindSeriesParams) (internal.Series, error)
	FindEmptyThumb(ctx context.Context, order internal.SortOrder) ([]internal.CreateScrapeRequestParams, error)
	FindOnGoing(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error)
	FindEmptyChapters(ctx context.Context, params internal.FindSeriesParams) ([]internal.CreateScrapeRequestParams, error)
}

type ScraperRepository interface {
	Create(ctx context.Context, params internal.CreateScrapeRequestParams) (internal.ScrapeRequest, error)
}

type SeriesSearchRepository interface {
	Index(ctx context.Context, series internal.Series) error
}

type Cron struct {
	provider    ProviderRepository
	series      SeriesRepository
	repo        JobRepository
	scraper     ScraperRepository
	search      SeriesSearchRepository
	cronMonitor *cronMonitor
	logger      *zap.Logger
	doneC       chan struct{}
}

func NewCron(
	provider ProviderRepository,
	series SeriesRepository,
	repo JobRepository,
	scraper ScraperRepository,
	search SeriesSearchRepository,
	logger *zap.Logger,
) *Cron {
	return &Cron{
		provider:    provider,
		series:      series,
		repo:        repo,
		scraper:     scraper,
		search:      search,
		cronMonitor: newCronMonitor(),
		logger:      logger,
		doneC:       make(chan struct{}),
	}
}

func (s *Cron) StartServer() (<-chan error, error) {
	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go s.handleShutdown(ctx, stop, errC)

	go s.serve(errC)

	return errC, nil
}

func (s *Cron) serve(errC chan<- error) {
	s.logger.Info("Listening and serving")

	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
		gocron.WithMonitor(s.cronMonitor),
	)
	if err != nil {
		errC <- internal.WrapErrorf(err, internal.ErrUnknown, "Failed to create scheduler")
	}

	if err := s.scheduleJobs(scheduler); err != nil {
		errC <- err
	}

	scheduler.Start()
	defer scheduler.Shutdown() // nolint:errcheck

	s.doneC <- struct{}{}
}

func (s *Cron) handleShutdown(ctx context.Context, stop context.CancelFunc, errC chan<- error) {
	<-ctx.Done()

	s.logger.Info("Shutdown signal received")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		_ = s.logger.Sync()
		stop()
		cancel()
		close(errC)
	}()

	if err := s.shutdown(ctxTimeout); err != nil {
		errC <- err
	}

	s.logger.Info("Shutdown completed")
}

func (s *Cron) shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")

	select {
	case <-ctx.Done():
		return internal.WrapErrorf(ctx.Err(), internal.ErrUnknown, "Context done")
	case <-s.doneC:
		return nil
	}
}
