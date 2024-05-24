package v1Handler

import (
	"errors"
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

// @Summary		Create request to scrape series list
// @Description	Create request to scrape series list
// @Tags			scrape-requests
// @Accept			json
// @Produce		json
// @Param			body	body		v1Binding.PostScrapeSeries	true	"Request body"
// @Success		201		{object}	v1Response.Response
// @Failure		400		{object}	v1Response.Response
// @Failure		404		{object}	v1Response.Response
// @Failure		500		{object}	v1Response.Response
// @Router			/api/v1/scrape-requests/series/list [post]
func (h *Handler) PostScrapeSeriesList(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.PostScrapeSeriesList")
	span.Name = "v1.PostScrapeSeriesList"
	defer span.Finish()

	var req v1Binding.PostScrapeSeries
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
		})
	} else if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.FindProviderUniqueV1", req)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
		})
	}

	receipt, err := h.prisma.CreateScrapeRequestV1(c.Request().Context(), provider)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateScrapeRequestV1", req)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
		})
	}

	startTime := time.Now()

	scrapeData, err := scraper.ScrapeSeriesList(req.Provider, provider.Scheme+provider.Host+provider.ListPath)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "scraper.ScrapeSeriesList", req)
		h.prisma.UpdateScrapeRequestUniqueV1(c.Request().Context(), receipt.ID, "failed", time.Since(startTime).Seconds(), err.Error())
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
		})
	}

	_, err = h.prisma.UpdateScrapeRequestUniqueV1(c.Request().Context(), receipt.ID, "success", time.Since(startTime).Seconds(), "Completed successfully")
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.UpdateScrapeRequestUniqueV1", req)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
		})
	}

	seriesList := make([]db.SeriesModel, len(scrapeData))

	for _, series := range scrapeData {
		_, err = h.prisma.CreateSeriesListRowV1(c.Request().Context(), receipt.ID, series)
		if err != nil {
			middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateSeriesListRowV1", series)
		}

		res, err := h.prisma.UpsertSeriesRowV1(c.Request().Context(), req.Provider, series)
		if err != nil {
			middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.UpsertSeriesRowV1", series)
		}

		seriesList = append(seriesList, *res)

		err = h.redis.UnsetSeriesV1(c.Request().Context(), req.Provider, series.Slug)
		if err != nil {
			middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetSeriesV1", series)
		}
	}

	err = h.redis.UnsetSeriesListV1(c.Request().Context(), req.Provider)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetSeriesListV1", req)
	}

	result := v1Response.NewGetSeriesListData(provider, seriesList)

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusCreated, v1Response.Response{
		Error:   false,
		Message: "Created",
		Data:    result,
	})
}
