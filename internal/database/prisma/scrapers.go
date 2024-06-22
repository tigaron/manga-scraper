package prisma

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
)

type ScraperRepo struct {
	q *PrismaClient
}

func NewScraperRepo(prismaClient *PrismaClient) *ScraperRepo {
	return &ScraperRepo{
		q: prismaClient,
	}
}

func (s *ScrapeRequestModel) toScrapeRequest() internal.ScrapeRequest {
	return internal.ScrapeRequest{
		ID:          s.ID,
		Type:        internal.ScrapeRequestType(s.Type),
		BaseURL:     s.BaseURL,
		RequestPath: s.RequestPath,
		Provider:    s.Provider,
		Series:      s.Series,
		Chapter:     s.Chapter,
		Status:      internal.ScrapeRequestStatus(s.Status),
		Retries:     s.Retries,
		TotalTime:   s.TotalTime,
		Error:       s.Error,
		Message:     s.Message,
	}
}

func (r *ScraperRepo) Create(ctx context.Context, params internal.CreateScrapeRequestParams) (internal.ScrapeRequest, error) {
	defer newSentrySpan(ctx, "ScraperRepo.Create").Finish()

	receipt, err := r.q.ScrapeRequest.CreateOne(
		ScrapeRequest.Type.Set(ScrapeRequestType(params.Type)),
		ScrapeRequest.BaseURL.Set(params.BaseURL),
		ScrapeRequest.RequestPath.Set(params.RequestPath),
		ScrapeRequest.Provider.Set(params.Provider),
		ScrapeRequest.Series.Set(params.Series),
		ScrapeRequest.Chapter.Set(params.Chapter),
		ScrapeRequest.Status.Set(string(params.Status)),
		ScrapeRequest.Retries.Set(0),
		ScrapeRequest.TotalTime.Set(0),
		ScrapeRequest.Error.Set(false),
		ScrapeRequest.Message.Set(""),
	).Exec(ctx)
	if err != nil {
		return internal.ScrapeRequest{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to create scrape request")
	}

	return receipt.toScrapeRequest(), nil
}

func (r *ScraperRepo) Find(ctx context.Context, id string) (internal.ScrapeRequest, error) {
	defer newSentrySpan(ctx, "ScraperRepo.Find").Finish()

	receipt, err := r.q.ScrapeRequest.FindUnique(
		ScrapeRequest.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return internal.ScrapeRequest{}, internal.WrapErrorf(err, internal.ErrNotFound, "scrape request not found")
	}

	return receipt.toScrapeRequest(), nil
}

func (r *ScraperRepo) FindPendings(ctx context.Context, params internal.FindScrapeRequestParams) ([]internal.ScrapeRequest, error) {
	defer newSentrySpan(ctx, "ScraperRepo.FindPending").Finish()

	receipts, err := r.q.ScrapeRequest.FindMany(
		ScrapeRequest.Status.Equals(string(params.Status)),
	).Take(params.Size).Skip(params.Page * params.Size).OrderBy(
		ScrapeRequest.CreatedAt.Order(newSortOrder(params.Order)),
	).Exec(ctx)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find pending scrape requests")
	}

	var result []internal.ScrapeRequest
	for _, receipt := range receipts {
		result = append(result, receipt.toScrapeRequest())
	}

	return result, nil
}

func (r *ScraperRepo) Update(ctx context.Context, params internal.UpdateScrapeRequestParams) (internal.ScrapeRequest, error) {
	defer newSentrySpan(ctx, "ScraperRepo.Update").Finish()

	receipt, err := r.q.ScrapeRequest.FindUnique(
		ScrapeRequest.ID.Equals(params.ID),
	).Update(
		ScrapeRequest.Status.Set(string(params.Status)),
		ScrapeRequest.Retries.Increment(1),
		ScrapeRequest.TotalTime.Set(params.TotalTime),
		ScrapeRequest.Error.Set(params.Error),
		ScrapeRequest.Message.Set(params.Message),
	).Exec(ctx)
	if err != nil {
		return internal.ScrapeRequest{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to update scrape request")
	}

	return receipt.toScrapeRequest(), nil
}

func (r *ScraperRepo) Delete(ctx context.Context, id string) error {
	defer newSentrySpan(ctx, "ScraperRepo.Delete").Finish()

	_, err := r.q.ScrapeRequest.FindUnique(
		ScrapeRequest.ID.Equals(id),
	).Delete().Exec(ctx)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "failed to delete scrape request")
	}

	return nil
}
