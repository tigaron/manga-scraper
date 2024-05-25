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

// @Summary		Create request to scrape chapter list
// @Description	Create request to scrape chapter list
// @Tags			scrape-requests
// @Accept			json
// @Produce		json
// @Param			body	body		v1Binding.PostScrapeChapterList	true	"Request body"
// @Success		201		{object}	v1Response.Response
// @Failure		400		{object}	v1Response.Response
// @Failure		404		{object}	v1Response.Response
// @Failure		500		{object}	v1Response.Response
// @Router			/api/v1/scrape-requests/chapters/list [post]
func (h *Handler) PostScrapeChapterList(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.PostScrapeChapterList")
	span.Name = "v1.PostScrapeChapterList"
	defer span.Finish()

	var req v1Binding.PostScrapeChapterList
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

	receipt, err := h.prisma.CreateChapterListScrapeRequestV1(c.Request().Context(), provider, series)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateChapterListScrapeRequestV1", req)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to create scrape request",
		})
	}

	startTime := time.Now()

	scrapeData, err := scraper.ScrapeChapterList(h.config.RodURL, req.Provider, provider.Scheme+provider.Host+series.SourcePath)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "scraper.ScrapeChapterList", req)
		h.prisma.UpdateScrapeRequestUniqueV1(c.Request().Context(), receipt.ID, "failed", time.Since(startTime).Seconds(), err.Error())
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to scrape chapter list",
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

	result := make([]v1Response.ChapterData, 0, len(scrapeData))

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(len(scrapeData))

	for _, chapter := range scrapeData {
		go func(ch v1Model.ChapterList) {
			defer wg.Done()
			_, err = h.prisma.CreateChapterListRowV1(c.Request().Context(), receipt.ID, ch)
			if err != nil {
				middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateChapterListRowV1", req)
			}

			upsertChapter, err := h.prisma.UpsertChaptersRowV1(c.Request().Context(), req.Provider, req.Series, ch)
			if err != nil {
				middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.UpsertSeriesRowV1", req)
			}

			mu.Lock()
			result = append(result, v1Response.NewChapterData(provider, series, upsertChapter))
			mu.Unlock()

			err = h.redis.UnsetChapterV1(c.Request().Context(), req.Provider, req.Series, ch.Slug)
			if err != nil {
				middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetChapterV1", req)
			}
		}(chapter)
	}

	wg.Wait()

	err = h.redis.UnsetChapterListV1(c.Request().Context(), req.Provider, req.Series)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetChapterListV1", req)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusCreated, v1Response.Response{
		Error:   false,
		Message: "Created",
		Data:    result,
	})
}
