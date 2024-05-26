package v1Handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
	"fourleaves.studio/manga-scraper/internal/scraper"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Create request to scrape series detail
// @Description	Create request to scrape series detail
// @Tags			scrape-requests
// @Accept			json
// @Produce		json
// @Param			body	body		v1Binding.PutScrapeSeriesDetail	true	"Request body"
// @Success		200		{object}	v1Response.Response
// @Failure		400		{object}	v1Response.Response
// @Failure		403		{object}	v1Response.Response
// @Failure		404		{object}	v1Response.Response
// @Failure		500		{object}	v1Response.Response
// @Failure		503		{object}	v1Response.Response
// @Router			/api/v1/scrape-requests/series/detail [put]
func (h *Handler) PutScrapeSeriesDetail(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.PutScrapeSeriesDetail")
	span.Name = "v1.PutScrapeSeriesDetail"
	defer span.Finish()

	var req v1Binding.PutScrapeSeriesDetail
	err := c.Bind(&req)
	if err != nil {
		span.Status = sentry.SpanStatusInvalidArgument
		return c.JSON(http.StatusBadRequest, v1Response.Response{
			Error:   true,
			Message: "Bad Request",
			Detail:  err.Error(),
		})
	}

	err = c.Validate(&req)
	if err != nil {
		span.Status = sentry.SpanStatusInvalidArgument
		return c.JSON(http.StatusBadRequest, v1Response.Response{
			Error:   true,
			Message: "Bad Request",
			Detail:  middlewares.FormatValidationError(err),
		})
	}

	provider, err := h.prisma.FindProviderUniqueV1(c.Request().Context(), req.Provider)
	if errors.Is(err, db.ErrNotFound) {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  fmt.Sprintf("Provider with slug '%s' not found", req.Provider),
		})
	} else if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.FindProviderUniqueV1", req)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to find provider",
		})
	}

	series, err := h.prisma.FindSeriesUniqueV1(c.Request().Context(), req.Provider, req.Series)
	if errors.Is(err, db.ErrNotFound) {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  fmt.Sprintf("Series with slug '%s' not found", req.Series),
		})
	} else if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "prisma.FindSeriesUniqueV1")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to find series",
		})
	}

	receipt, err := h.prisma.CreateSeriesDetailScrapeRequestV1(c.Request().Context(), provider, series)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateSeriesDetailScrapeRequestV1", req)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to create scrape request",
		})
	}

	startTime := time.Now()

	scrapeData, err := scraper.ScrapeSeriesDetail(h.config.RodURL, req.Provider, provider.Scheme+provider.Host+series.SourcePath)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "scraper.ScrapeSeriesDetail", req)
		h.prisma.UpdateScrapeRequestUniqueV1(c.Request().Context(), receipt.ID, "failed", time.Since(startTime).Seconds(), err.Error())
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to scrape series detail",
		})
	}

	_, err = h.prisma.UpdateScrapeRequestUniqueV1(c.Request().Context(), receipt.ID, "success", time.Since(startTime).Seconds(), "Completed successfully")
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.UpdateScrapeRequestUniqueV1", req)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to update scrape request",
		})
	}

	_, err = h.prisma.CreateSeriesDetailRowV1(c.Request().Context(), receipt.ID, scrapeData)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateSeriesDetailRowV1", req)
	}

	updatedSeries, err := h.prisma.UpdateDetailSeriesRowV1(c.Request().Context(), req.Provider, series.Slug, scrapeData)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.UpdateDetailSeriesRowV1", req)
	}

	err = h.redis.UnsetSeriesV1(c.Request().Context(), req.Provider, series.Slug)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetSeriesV1", series)
	}

	err = h.redis.UnsetSeriesListV1(c.Request().Context(), req.Provider)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetSeriesListV1", req)
	}

	result := v1Response.NewSeriesData(provider, updatedSeries)

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
		Data:    result,
	})
}
