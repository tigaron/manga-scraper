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

// @Summary		Get all series list
// @Description	Get all series list
// @Tags			series
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug"	example(asura)
// @Success		200				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/series/{provider_slug}/_all [get]
func (h *Handler) GetSeriesListAll(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetSeriesListAll")
	span.Name = "v1.GetSeriesListAll"
	defer span.Finish()

	providerSlug := c.Param("provider_slug")

	c.Logger().Debugj(map[string]interface{}{
		"_source":       "v1.GetSeriesListAll",
		"provider_slug": providerSlug,
	})

	cache, err := h.redis.GetSeriesListAllV1(c.Request().Context(), providerSlug)
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
			Detail:  "Failed to find provider",
		})
	}

	series, err := h.prisma.FindSeriesManyV1(c.Request().Context(), providerSlug)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "prisma.FindSeriesManyV1")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to find series",
		})
	}

	if len(series) == 0 {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  fmt.Sprintf("Series with provider slug '%s' not found", providerSlug),
		})
	}

	result := v1Response.NewSeriesListData(provider, series)

	err = h.redis.SetSeriesListAllV1(c.Request().Context(), providerSlug, result)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "redis.SetSeriesListAllV1")
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
		Data:    result,
	})
}
