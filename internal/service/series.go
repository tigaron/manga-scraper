package service

import (
	"context"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/labstack/echo/v4"
	"github.com/mercari/go-circuitbreaker"
)

type SeriesRepository interface {
	CreateInit(ctx context.Context, params internal.CreateInitSeriesParams) (internal.Series, error)
	Find(ctx context.Context, params internal.FindSeriesParams) (internal.Series, error)
	FindBC(ctx context.Context, params internal.FindSeriesParams) (internal.SeriesBC, error)
	FindAll(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error)
	FindPaginated(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error)
	UpdateInit(ctx context.Context, params internal.UpdateInitSeriesParams) (internal.Series, error)
	UpdateLatest(ctx context.Context, params internal.UpdateLatestSeriesParams) (internal.Series, error)
	Delete(ctx context.Context, params internal.FindSeriesParams) error
}

type SeriesSearchRepository interface {
	Search(ctx context.Context, q string) ([]internal.Series, error)
	Index(ctx context.Context, series internal.Series) error
	Delete(ctx context.Context, provider, slug string) error
}

type SeriesService struct {
	repo   SeriesRepository
	search SeriesSearchRepository
	cb     *circuitbreaker.CircuitBreaker
}

func NewSeriesService(repo SeriesRepository, search SeriesSearchRepository, logger echo.Logger) *SeriesService {
	return &SeriesService{
		repo:   repo,
		search: search,
		cb: circuitbreaker.New(
			circuitbreaker.WithOpenTimeout(time.Minute*2),
			circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncConsecutiveFailures(3)),
			circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
				logger.Infoj(map[string]interface{}{
					"_source": "SeriesService",
					"_msg":    "circuit breaker state change",
					"old":     string(oldState),
					"new":     string(newState),
				})
			}),
		),
	}
}

func (s *SeriesService) CreateInit(ctx context.Context, params internal.CreateInitSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesService.CreateInit").Finish()

	if err := params.Validate(); err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrInvalidInput, "params.Validate")
	}

	series, err := s.repo.CreateInit(ctx, params)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.CreateInit")
	}

	err = s.search.Index(ctx, series)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "search.Index")
	}

	return series, nil
}

func (s *SeriesService) Search(ctx context.Context, q string) (results []internal.Series, err error) {
	defer newSentrySpan(ctx, "SeriesService.Search").Finish()

	if !s.cb.Ready() {
		return nil, internal.WrapErrorf(nil, internal.ErrUnknown, "circuit breaker is open")
	}

	defer func() {
		err = s.cb.Done(ctx, err)
	}()

	results, err = s.search.Search(ctx, q)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "search.Search")
	}

	return results, nil
}

func (s *SeriesService) Index(ctx context.Context, series []internal.Series) (err error) {
	defer newSentrySpan(ctx, "SeriesService.Index").Finish()

	if !s.cb.Ready() {
		return internal.WrapErrorf(nil, internal.ErrUnknown, "circuit breaker is open")
	}

	defer func() {
		err = s.cb.Done(ctx, err)
	}()

	var loopErr error

	for i := range series {
		err := s.search.Index(ctx, series[i])
		if err != nil {
			loopErr = internal.WrapErrorf(err, internal.ErrUnknown, "search.Index")
			break
		}
	}

	return loopErr
}

func (s *SeriesService) Find(ctx context.Context, params internal.FindSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesService.Find").Finish()

	series, err := s.repo.Find(ctx, params)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.Find")
	}

	return series, nil
}

func (s *SeriesService) FindBC(ctx context.Context, params internal.FindSeriesParams) (internal.SeriesBC, error) {
	defer newSentrySpan(ctx, "SeriesService.FindBC").Finish()

	series, err := s.repo.FindBC(ctx, params)
	if err != nil {
		return internal.SeriesBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.FindBC")
	}

	return series, nil
}

func (s *SeriesService) FindAll(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesService.FindAll").Finish()

	series, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "repo.FindAll")
	}

	return series, nil
}

func (s *SeriesService) FindPaginated(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesService.FindPaginated").Finish()

	series, err := s.repo.FindPaginated(ctx, params)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "repo.FindPaginated")
	}

	return series, nil
}

func (s *SeriesService) UpdateInit(ctx context.Context, params internal.UpdateInitSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesService.UpdateInit").Finish()

	if err := params.Validate(); err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrInvalidInput, "params.Validate")
	}

	series, err := s.repo.UpdateInit(ctx, params)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.UpdateInit")
	}

	err = s.search.Index(ctx, series)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "search.Index")
	}

	return series, nil
}

func (s *SeriesService) UpdateLatest(ctx context.Context, params internal.UpdateLatestSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesService.UpdateLatest").Finish()

	if err := params.Validate(); err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrInvalidInput, "params.Validate")
	}

	series, err := s.repo.UpdateLatest(ctx, params)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.UpdateLatest")
	}

	err = s.search.Index(ctx, series)
	if err != nil {
		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "search.Index")
	}

	return series, nil
}

func (s *SeriesService) Delete(ctx context.Context, params internal.FindSeriesParams) error {
	defer newSentrySpan(ctx, "SeriesService.Delete").Finish()

	err := s.repo.Delete(ctx, params)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "repo.Delete")
	}

	err = s.search.Delete(ctx, params.Provider, params.Slug)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "search.Delete")
	}

	return nil
}
