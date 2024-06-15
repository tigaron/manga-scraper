package v1Handler

import (
	"errors"
	"fmt"
	"net/http"

	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get paginated chapter list
// @Description	Get paginated chapter list
// @Tags			chapters
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug"	example(asura)
// @Param			series_slug		path		string	true	"Series slug"	example(reincarnator)
// @Param			page			query		string	true	"Page"			example(10)
// @Param			size			query		string	true	"Size"			example(100)
// @Success		200				{object}	ResponseV1
// @Failure		400				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/chapters/{provider_slug}/{series_slug} [get]
func (h *Handler) GetChapterListPaginated(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetChapterListPaginated")
	span.Name = "v1.GetChapterListPaginated"
	defer span.Finish()

	var req v1Binding.PaginatedRequest
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

	providerSlug := c.Param("provider_slug")
	seriesSlug := c.Param("series_slug")

	c.Logger().Debugj(map[string]interface{}{
		"_source":       "v1.GetChapterListPaginated",
		"provider_slug": providerSlug,
		"series_slug":   seriesSlug,
		"page":          req.Page,
		"size":          req.Size,
	})

	cache, err := h.redis.GetChapterListV1(c.Request().Context(), providerSlug, seriesSlug, req.Page, req.Size)
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
			Detail:  "Failed to find series",
		})
	}

	chapterList, err := h.prisma.FindChaptersManyPaginatedV1(c.Request().Context(), providerSlug, seriesSlug, req)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "prisma.FindChaptersManyPaginatedV1")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to find chapters",
		})
	}

	if len(chapterList) == 0 {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  fmt.Sprintf("Chapters with provider slug '%s' and series slug '%s' not found", providerSlug, seriesSlug),
		})
	}

	var prevPage, nextPage, total int

	if req.Page >= 2 {
		prevPage = req.Page - 1
	}

	if len(chapterList) == req.Size {
		nextPage = req.Page + 1
	}

	total = len(chapterList)

	paginationData := v1Response.PaginationData{
		PrevPage: prevPage,
		NextPage: nextPage,
		Total:    total,
	}

	result := v1Response.NewChapterListPaginatedData(provider, series, chapterList, paginationData)

	err = h.redis.SetChapterListV1(c.Request().Context(), providerSlug, seriesSlug, req.Page, req.Size, result)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "redis.SetChapterListV1")
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
		Data:    result,
	})
}
