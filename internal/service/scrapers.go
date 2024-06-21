package service

import (
	"context"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/labstack/echo/v4"
	"github.com/mercari/go-circuitbreaker"
)

type ScrapeRequestRepository interface {
	Create(ctx context.Context, params internal.CreateScrapeRequestParams) (internal.ScrapeRequest, error)
	Find(ctx context.Context, id string) (internal.ScrapeRequest, error)
	FindPendings(ctx context.Context, params internal.FindScrapeRequestParams) ([]internal.ScrapeRequest, error)
	Update(ctx context.Context, params internal.UpdateScrapeRequestParams) (internal.ScrapeRequest, error)
	Delete(ctx context.Context, id string) error
}

type ScrapeRequestMessageBroker interface {
	Created(ctx context.Context, params internal.ScrapeRequest) error
}

type ScraperService struct {
	repo      ScrapeRequestRepository
	msgBroker ScrapeRequestMessageBroker
	cb        *circuitbreaker.CircuitBreaker
}

func NewScraperService(repo ScrapeRequestRepository, msgBroker ScrapeRequestMessageBroker, logger echo.Logger) *ScraperService {
	return &ScraperService{
		repo:      repo,
		msgBroker: msgBroker,
		cb: circuitbreaker.New(
			circuitbreaker.WithOpenTimeout(time.Minute*2),
			circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncConsecutiveFailures(3)),
			circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
				logger.Infoj(map[string]interface{}{
					"_source": "ScraperService",
					"_msg":    "circuit breaker state change",
					"old":     string(oldState),
					"new":     string(newState),
				})
			}),
		),
	}
}

func (s *ScraperService) Create(ctx context.Context, params internal.CreateScrapeRequestParams) (receipt internal.ScrapeRequest, err error) {
	defer newSentrySpan(ctx, "Scraper.Create").Finish()

	if !s.cb.Ready() {
		return internal.ScrapeRequest{}, internal.WrapErrorf(nil, internal.ErrUnknown, "circuit breaker is open")
	}

	defer func() {
		err = s.cb.Done(ctx, err)
	}()

	if err := params.Validate(); err != nil {
		return internal.ScrapeRequest{}, internal.WrapErrorf(err, internal.ErrInvalidInput, "params.Validate")
	}

	receipt, err = s.repo.Create(ctx, params)
	if err != nil {
		return internal.ScrapeRequest{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.Create")
	}

	err = s.msgBroker.Created(ctx, receipt)
	if err != nil {
		return internal.ScrapeRequest{}, internal.WrapErrorf(err, internal.ErrUnknown, "msgBroker.Create")
	}

	return receipt, nil
}

func (s *ScraperService) Find(ctx context.Context, id string) (internal.ScrapeRequest, error) {
	defer newSentrySpan(ctx, "Scraper.Find").Finish()

	receipt, err := s.repo.Find(ctx, id)
	if err != nil {
		return internal.ScrapeRequest{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.Find")
	}

	return receipt, nil
}

func (s *ScraperService) FindPendings(ctx context.Context, params internal.FindScrapeRequestParams) ([]internal.ScrapeRequest, error) {
	defer newSentrySpan(ctx, "Scraper.FindPendings").Finish()

	receipts, err := s.repo.FindPendings(ctx, params)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "repo.FindPendings")
	}

	return receipts, nil
}

func (s *ScraperService) Update(ctx context.Context, params internal.UpdateScrapeRequestParams) (internal.ScrapeRequest, error) {
	defer newSentrySpan(ctx, "Scraper.Update").Finish()

	receipt, err := s.repo.Update(ctx, params)
	if err != nil {
		return internal.ScrapeRequest{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.Update")
	}

	return receipt, nil
}

func (s *ScraperService) Delete(ctx context.Context, id string) error {
	defer newSentrySpan(ctx, "Scraper.Delete").Finish()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "repo.Delete")
	}

	return nil
}
