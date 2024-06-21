package chapters

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

type ChapterService interface {
	CreateInit(ctx context.Context, params internal.CreateInitChapterParams) (internal.Chapter, error)
	Find(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error)
	FindBC(ctx context.Context, params internal.FindChapterParams) (internal.ChapterBC, error)
	FindLatest(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error)
	Count(ctx context.Context, params internal.FindChapterParams) (int, error)
	FindAll(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error)
	FindListWithRel(ctx context.Context, params internal.FindChapterParams) (internal.ChapterList, error)
	FindPaginated(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error)
	UpdateInit(ctx context.Context, params internal.UpdateInitChapterParams) (internal.Chapter, error)
	Delete(ctx context.Context, params internal.FindChapterParams) error
}

type ChapterHandler struct {
	svc ChapterService
}

func NewChapterHandler(svc ChapterService) *ChapterHandler {
	return &ChapterHandler{
		svc: svc,
	}
}

func (h *ChapterHandler) Register(g *echo.Group) {
	g.GET("/:provider_slug/:series_slug", h.FindPaginated)
	g.GET("/:provider_slug/:series_slug/_all", h.FindAll)
	g.GET("/:provider_slug/:series_slug/_list", h.FindListWithRel)
	g.GET("/:provider_slug/:series_slug/:chapter_slug", h.Find)
	g.GET("/:provider_slug/:series_slug/:chapter_slug/_bc", h.FindBC)
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
	Chapters []internal.Chapter `json:"chapters"`
}

func newSentrySpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	span.Name = "fourleaves.studio/manga-scraper/internal/rest/v1/chapters"

	return span
}
