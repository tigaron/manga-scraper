package scrapers

import (
	"net/http"

	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get scrape request by ID
// @Description	Get scrape request by ID
// @Security		TokenAuth
// @Tags			scrapers
// @Produce		json
// @Param			id	path		string	true	"Request ID"	example(550e8400-e29b-41d4-a716-446655440000)
// @Success		200	{object}	ResponseV1
// @Failure		401	{object}	ResponseV1
// @Failure		403	{object}	ResponseV1
// @Failure		404	{object}	ResponseV1
// @Failure		500	{object}	ResponseV1
// @Router			/api/v1/scrapers/{id} [get]
func (h *ScraperHandler) Find(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.Find")
	defer span.Finish()

	id := c.Param("id")

	receipt, err := h.svc.Find(c.Request().Context(), id)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to get scrape request", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Handler.Response{
		Error:   false,
		Message: "OK",
		Data:    receipt,
	})
}
