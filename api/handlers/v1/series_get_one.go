package v1Handler

import (
	"errors"
	"fmt"
	"net/http"

	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get series by slug
// @Description	Get series by slug
// @Tags			series
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug" example(asura)
// @Param			series_slug	path		string	true	"Series slug" example(reincarnator)
// @Success		200				{object}	v1Response.Response
// @Failure		404				{object}	v1Response.Response
// @Failure		500				{object}	v1Response.Response
// @Router			/api/v1/series/{provider_slug}/{series_slug} [get]
func (h *Handler) GetSeries(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetProvider")
	span.Name = "v1.GetProvider"
	defer span.Finish()

	providerSlug := c.Param("provider_slug")
	seriesSlug := c.Param("series_slug")

	cache, err := h.redis.GetSeriesV1(c.Request().Context(), providerSlug, seriesSlug)
	if err == nil {
		span.Status = sentry.SpanStatusOK
		return c.JSON(http.StatusOK, v1Response.Response{
			Error:   false,
			Message: "OK",
			Data:    cache,
		})
	}

	provider, err := h.prisma.FindProviderUniqueV1(c.Request().Context(), providerSlug)
	if errors.Is(err, db.ErrNotFound) {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  fmt.Sprintf("Provider with slug '%s' not found", providerSlug),
		})
	} else if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "prisma.FindProviderUniqueV1")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
		})
	}

	series, err := h.prisma.FindSeriesUniqueV1(c.Request().Context(), providerSlug, seriesSlug)
	if errors.Is(err, db.ErrNotFound) {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  fmt.Sprintf("Series with slug '%s' not found", seriesSlug),
		})
	} else if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "prisma.FindSeriesUniqueV1")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
		})
	}

	result := v1Response.NewSeriesData(provider, series)

	err = h.redis.SetSeriesV1(c.Request().Context(), providerSlug, seriesSlug, result)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "redis.SetSeriesV1")
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
		Data:    result,
	})
}
