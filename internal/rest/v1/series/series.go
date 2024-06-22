package series

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
	"fourleaves.studio/manga-scraper/internal/rest/middlewares"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

type SeriesService interface {
	CreateInit(ctx context.Context, params internal.CreateInitSeriesParams) (internal.Series, error)
	Search(ctx context.Context, q string) ([]internal.Series, error)
	Index(ctx context.Context, series []internal.Series) error
	Find(ctx context.Context, params internal.FindSeriesParams) (internal.Series, error)
	FindBC(ctx context.Context, params internal.FindSeriesParams) (internal.SeriesBC, error)
	FindAll(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error)
	FindPaginated(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error)
	UpdateInit(ctx context.Context, params internal.UpdateInitSeriesParams) (internal.Series, error)
	UpdateLatest(ctx context.Context, params internal.UpdateLatestSeriesParams) (internal.Series, error)
	Delete(ctx context.Context, params internal.FindSeriesParams) error
}

type SeriesHandler struct {
	svc SeriesService
}

func NewSeriesHandler(svc SeriesService) *SeriesHandler {
	return &SeriesHandler{
		svc: svc,
	}
}

func (h *SeriesHandler) Register(g *echo.Group, mid *middlewares.Middleware) {
	g.GET("", h.Search)
	g.PUT("/:provider_slug", h.Index, mid.IsAdmin)
	g.GET("/:provider_slug", h.FindPaginated)
	g.GET("/:provider_slug/_all", h.FindAll)
	g.GET("/:provider_slug/:series_slug", h.Find)
	g.GET("/:provider_slug/:series_slug/_bc", h.FindBC)
	// g.GET("/:provider_slug/:series_slug/_bc", h.GetSeriesBreadcrumbs)
	// g.PUT("/:provider_slug/:series_slug/_chc", h.PutSeriesChaptersCount /* , middlewares.IsAdmin(s.config.AdminSub) */)
	// g.PUT("/:provider_slug/:series_slug/_lch", h.PutSeriesLastChapter /* , middlewares.IsAdmin(s.config.AdminSub) */)
}

type PaginatedRequest struct {
	Sort string `query:"sort" validate:"omitempty,oneof=asc desc" example:"asc"`
	Page int    `query:"page" validate:"required,gt=0" example:"1"`
	Size int    `query:"size" validate:"required,gt=0,lte=100" example:"10"`
}

type PaginationData struct {
	PrevPage int `json:"prevPage,omitempty"`
	NextPage int `json:"nextPage,omitempty"`
	Total    int `json:"total,omitempty"`
}

type PaginatedResponse struct {
	PaginationData
	Series []internal.Series `json:"series"`
}

func newSentrySpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	span.Name = "fourleaves.studio/manga-scraper/internal/rest/v1/series"

	return span
}
