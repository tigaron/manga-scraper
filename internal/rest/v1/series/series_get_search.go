package series

import (
	"net/http"

	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get series search result
// @Description	Get series search result
// @Tags			series
// @Produce		json
// @Param			q	query		string	true	"Query"	example(warrior high school)
// @Success		200	{object}	ResponseV1
// @Failure		400	{object}	ResponseV1
// @Failure		404	{object}	ResponseV1
// @Failure		500	{object}	ResponseV1
// @Router			/api/v1/series [get]
func (h *SeriesHandler) Search(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.Search")
	defer span.Finish()

	q := c.QueryParam("q")

	result, err := h.svc.Search(c.Request().Context(), q)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to search series", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Handler.Response{
		Error:   false,
		Message: "OK",
		Data:    result,
	})
}
