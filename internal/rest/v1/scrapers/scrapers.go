package scrapers

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
	"fourleaves.studio/manga-scraper/internal/rest/middlewares"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

type ScraperService interface {
	Create(ctx context.Context, params internal.CreateScrapeRequestParams) (internal.ScrapeRequest, error)
	Find(ctx context.Context, id string) (internal.ScrapeRequest, error)
	FindPendings(ctx context.Context, params internal.FindScrapeRequestParams) ([]internal.ScrapeRequest, error)
	Update(ctx context.Context, params internal.UpdateScrapeRequestParams) (internal.ScrapeRequest, error)
	Delete(ctx context.Context, id string) error
}

type ProviderService interface {
	Find(ctx context.Context, slug string) (internal.Provider, error)
}

type SeriesService interface {
	Find(ctx context.Context, params internal.FindSeriesParams) (internal.Series, error)
}

type ChapterService interface {
	Find(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error)
}

type ScraperHandler struct {
	svc      ScraperService
	provider ProviderService
	series   SeriesService
	chapter  ChapterService
}

func NewScraperHandler(
	svc ScraperService,
	provider ProviderService,
	series SeriesService,
	chapter ChapterService,
) *ScraperHandler {
	return &ScraperHandler{
		svc:      svc,
		provider: provider,
		series:   series,
		chapter:  chapter,
	}
}

func (h *ScraperHandler) Register(g *echo.Group, mid *middlewares.Middleware) {
	g.POST("", h.Create, mid.IsAdmin)
	// g.GET("", h.FindPendings)
	g.GET("/:id", h.Find, mid.IsAdmin)
	// g.PUT("/:id", h.Update)
	// g.DELETE("/:id", h.Delete)
}

type CreateScrapeRequest struct {
	Type     string `json:"type" validate:"required,oneof=SERIES_LIST SERIES_DETAIL CHAPTER_LIST CHAPTER_DETAIL" example:"CHAPTER_DETAIL"`
	Provider string `json:"provider" validate:"required" example:"asura"`
	Series   string `json:"series,omitempty" example:"reincarnator"`
	Chapter  string `json:"chapter,omitempty" example:"reincarnator-chapter-1"`
} // @name CreateScrapeRequest

func newSentrySpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	span.Name = "fourleaves.studio/manga-scraper/internal/rest/v1/scrapers"

	return span
}
