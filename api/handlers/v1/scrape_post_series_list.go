package v1Handler

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Model "fourleaves.studio/manga-scraper/api/models/v1"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
	"fourleaves.studio/manga-scraper/internal/scraper"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Create request to scrape series list
// @Description	Create request to scrape series list
// @Security		TokenAuth
// @Tags			scrape-requests
// @Accept			json
// @Produce		json
// @Param			body	body		PostScrapeSeriesList	true	"Request body"
// @Success		201		{object}	ResponseV1
// @Failure		400		{object}	ResponseV1
// @Failure		401		{object}	ResponseV1
// @Failure		403		{object}	ResponseV1
// @Failure		404		{object}	ResponseV1
// @Failure		500		{object}	ResponseV1
// @Failure		503		{object}	ResponseV1
// @Router			/api/v1/scrape-requests/series/list [post]
func (h *Handler) PostScrapeSeriesList(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.PostScrapeSeriesList")
	span.Name = "v1.PostScrapeSeriesList"
	defer span.Finish()

	var req v1Binding.PostScrapeSeriesList
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

	c.Logger().Debugj(map[string]interface{}{
		"_source":  "v1.PostScrapeSeriesList",
		"provider": req.Provider,
	})

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

	receipt, err := h.prisma.CreateSeriesListScrapeRequestV1(c.Request().Context(), provider)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateScrapeRequestV1", req)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to create scrape request",
		})
	}

	startTime := time.Now()

	scrapeData, err := scraper.ScrapeSeriesList(c.Request().Context(), h.config.RodURL, req.Provider, provider.Scheme+provider.Host+provider.ListPath)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "scraper.ScrapeSeriesList", req)
		h.prisma.UpdateScrapeRequestUniqueV1(c.Request().Context(), receipt.ID, "failed", time.Since(startTime).Seconds(), err.Error())
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to scrape series list",
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

	result := make([]v1Response.SeriesData, 0, len(scrapeData))

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(len(scrapeData))

	for _, series := range scrapeData {
		go func(s v1Model.SeriesList) {
			defer wg.Done()

			_, err = h.prisma.CreateSeriesListRowV1(c.Request().Context(), receipt.ID, s)
			if err != nil {
				middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateSeriesListRowV1", req)
			}

			upsertSeries, err := h.prisma.UpsertSeriesRowV1(c.Request().Context(), req.Provider, s)
			if err != nil {
				middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.UpsertSeriesRowV1", req)
			}

			mu.Lock()
			result = append(result, v1Response.NewSeriesData(provider, upsertSeries))
			mu.Unlock()

			err = h.redis.UnsetSeriesV1(c.Request().Context(), req.Provider, s.Slug)
			if err != nil {
				middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetSeriesV1", req)
			}
		}(series)
	}

	wg.Wait()

	err = h.redis.UnsetSeriesListV1(c.Request().Context(), req.Provider)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetSeriesListV1", req)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusCreated, v1Response.Response{
		Error:   false,
		Message: "Created",
		Data:    result,
	})
}
