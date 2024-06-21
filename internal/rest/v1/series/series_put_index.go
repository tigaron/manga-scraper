package series

import (
	"net/http"

	"fourleaves.studio/manga-scraper/internal"
	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Feed the open search engine
// @Description	Feed the open search engine
// @Security		TokenAuth
// @Tags			series
// @Produce		json
// @Param			provider_slug path string true "Provider Slug"
// @Success		200	{object}	ResponseV1
// @Failure		401	{object}	ResponseV1
// @Failure		403	{object}	ResponseV1
// @Failure		404	{object}	ResponseV1
// @Failure		500	{object}	ResponseV1
// @Router			/api/v1/series/{provider_slug} [put]
func (h *SeriesHandler) Index(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.Index")
	defer span.Finish()

	providerSlug := c.Param("provider_slug")

	series, err := h.svc.FindAll(c.Request().Context(), internal.FindSeriesParams{
		Provider: providerSlug,
	})
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to get series", err, span)
	}

	err = h.svc.Index(c.Request().Context(), series)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to index series", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Handler.Response{
		Error:   false,
		Message: "OK",
	})
}
