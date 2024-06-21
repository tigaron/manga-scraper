package providersHandler

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
	"fourleaves.studio/manga-scraper/internal/rest/middlewares"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

type ProviderService interface {
	Create(ctx context.Context, params internal.ProviderParams) (internal.Provider, error)
	Find(ctx context.Context, slug string) (internal.Provider, error)
	FindBC(ctx context.Context, slug string) (internal.ProviderBC, error)
	FindAll(ctx context.Context, order internal.SortOrder) ([]internal.Provider, error)
	Update(ctx context.Context, params internal.ProviderParams) (internal.Provider, error)
	Delete(ctx context.Context, slug string) error
}

type ProviderHandler struct {
	svc ProviderService
}

func NewProviderHandler(svc ProviderService) *ProviderHandler {
	return &ProviderHandler{
		svc: svc,
	}
}

func (h *ProviderHandler) Register(g *echo.Group, mid *middlewares.Middleware) {
	g.POST("", h.Create, mid.IsAdmin)
	g.GET("", h.FindAll)
	g.GET("/:provider_slug", h.Find)
	g.PUT("/:provider_slug", h.Update, mid.IsAdmin)
	g.GET("/:provider_slug/_bc", h.FindBC)
	// g.DELETE("/:provider_slug", h.Delete)
}

type CreateProviderRequest struct {
	Slug     string `json:"slug" validate:"required" example:"asura"`
	Name     string `json:"name" validate:"required" example:"Asura Scans"`
	Scheme   string `json:"scheme" validate:"required" example:"https://"`
	Host     string `json:"host" validate:"required" example:"asuratoon.com"`
	ListPath string `json:"list_path" validate:"required" example:"/manga/list-mode/"`
	IsActive *bool  `json:"is_active" validate:"required" example:"true"`
} // @name CreateProviderRequest

type UpdateProviderRequest struct {
	Name     string `json:"name" validate:"required" example:"Asura Scans"`
	Scheme   string `json:"scheme" validate:"required" example:"https://"`
	Host     string `json:"host" validate:"required" example:"asuratoon.com"`
	ListPath string `json:"list_path" validate:"required" example:"/manga/list-mode/"`
	IsActive *bool  `json:"is_active" validate:"required" example:"true"`
} // @name UpdateProviderRequest

func newSentrySpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	span.Name = "fourleaves.studio/manga-scraper/internal/rest/v1/providers"

	return span
}
