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

type ScraperHandler struct {
	svc ScraperService
}

func NewScraperHandler(svc ScraperService) *ScraperHandler {
	return &ScraperHandler{
		svc: svc,
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
	Type        string `json:"type" validate:"required,oneof=SERIES_LIST SERIES_DETAIL CHAPTER_LIST CHAPTER_DETAIL" example:"CHAPTER_DETAIL"`
	Status      string `json:"status" validate:"required,oneof=PENDING COMPLETED FAILED" example:"PENDING"`
	BaseURL     string `json:"baseURL" validate:"required" example:"https://asuratoon.com"`
	RequestPath string `json:"requestPath" validate:"required" example:"/manga/list-mode/"`
	Provider    string `json:"provider" validate:"required" example:"asura"`
	Series      string `json:"series,omitempty" example:"reincarnator"`
	Chapter     string `json:"chapter,omitempty" example:"reincarnator-chapter-1"`
} // @name CreateScrapeRequest

func newSentrySpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	span.Name = "fourleaves.studio/manga-scraper/internal/rest/v1/scrapers"

	return span
}
