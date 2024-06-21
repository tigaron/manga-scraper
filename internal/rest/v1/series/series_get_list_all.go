package series

import (
	"net/http"

	"fourleaves.studio/manga-scraper/internal"
	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get all series list
// @Description	Get all series list
// @Tags			series
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug"	example(asura)
// @Param			sort			query		string	false	"Sort order"	enum(asc, desc)	default(asc)
// @Success		200				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/series/{provider_slug}/_all [get]
func (h *SeriesHandler) FindAll(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.FindAll")
	defer span.Finish()

	providerSlug := c.Param("provider_slug")
	sort := c.QueryParam("sort")

	params := internal.FindSeriesParams{
		Provider: providerSlug,
		Order:    internal.NewSortOrder(sort),
	}

	seriesList, err := h.svc.FindAll(c.Request().Context(), params)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to get series", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Handler.Response{
		Error:   false,
		Message: "OK",
		Data:    seriesList,
	})
}
