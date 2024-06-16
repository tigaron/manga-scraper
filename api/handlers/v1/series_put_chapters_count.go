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

// @Summary		Update series chapters count
// @Description	Update series chapters count
// @Security		TokenAuth
// @Tags			series
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug"	example(asura)
// @Param			series_slug		path		string	true	"Series slug"	example(reincarnator)
// @Success		200				{object}	ResponseV1
// @Failure		401				{object}	ResponseV1
// @Failure		403				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/series/{provider_slug}/{series_slug}/_chc [put]
func (h *Handler) PutSeriesChaptersCount(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.PutSeriesLastChapter")
	span.Name = "v1.PutSeriesLastChapter"
	defer span.Finish()

	providerSlug := c.Param("provider_slug")
	seriesSlug := c.Param("series_slug")

	c.Logger().Debugj(map[string]interface{}{
		"_source":  "v1.PutSeriesLastChapter",
		"provider": providerSlug,
		"series":   seriesSlug,
	})

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

	series, err := h.prisma.FindSeriesUniqueV1(c.Request().Context(), provider.Slug, seriesSlug)
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
			Detail:  "Failed to find series",
		})
	}

	chaptersCount, err := h.prisma.CountChaptersV1(c.Request().Context(), provider.Slug, series.Slug)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "prisma.CountChaptersV1")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to count chapters",
		})
	}

	_, err = h.prisma.UpdateChaptersCountV1(c.Request().Context(), provider.Slug, series.Slug, chaptersCount)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "prisma.UpdateChaptersCountV1")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to update chapters count",
		})
	}

	_ = h.redis.UnsetSeriesV1(c.Request().Context(), provider.Slug, series.Slug)
	_ = h.redis.UnsetSeriesListV1(c.Request().Context(), provider.Slug)

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "Success",
		Detail:  "Chapters count updated",
	})
}
