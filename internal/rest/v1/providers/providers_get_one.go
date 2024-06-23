package providers

import (
	"net/http"

	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get provider by slug
// @Description	Get provider by slug
// @Tags			providers
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug" example(asura)
// @Success		200				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/providers/{provider_slug} [get]
func (h *ProviderHandler) Find(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.Find")
	defer span.Finish()

	providerSlug := c.Param("provider_slug")

	provider, err := h.svc.Find(c.Request().Context(), providerSlug)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to get providers", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Handler.Response{
		Error:   false,
		Message: "OK",
		Data:    provider,
	})
}
