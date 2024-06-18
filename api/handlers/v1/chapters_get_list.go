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

// @Summary		Get chapter list only
// @Description	Get chapter list only
// @Tags			chapters
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug"	example(asura)
// @Param			series_slug		path		string	true	"Series slug"	example(reincarnator)
// @Param			sort			query		string	false	"Sort order"	enum(asc, desc)	default(asc)
// @Success		200				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/chapters/{provider_slug}/{series_slug}/_list [get]
func (h *Handler) GetChapterList(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetChapterList")
	span.Name = "v1.GetChapterList"
	defer span.Finish()

	providerSlug := c.Param("provider_slug")
	seriesSlug := c.Param("series_slug")
	sort := c.QueryParam("sort")

	var order db.SortOrder

	switch sort {
	case "asc":
		order = db.ASC
	case "desc":
		order = db.DESC
	default:
		order = db.ASC
	}

	c.Logger().Debugj(map[string]interface{}{
		"_source":       "v1.GetChapterList",
		"provider_slug": providerSlug,
		"series_slug":   seriesSlug,
	})

	cache, err := h.redis.GetChapterListOnlyV1(c.Request().Context(), providerSlug, seriesSlug)
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

	chapterList, err := h.prisma.FindChaptersListV1(c.Request().Context(), providerSlug, seriesSlug)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "prisma.FindChaptersListV1")
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

	result := v1Response.NewListAllChapterData(provider, series, chapterList)

	err = h.redis.SetChapterListOnlyV1(c.Request().Context(), providerSlug, seriesSlug, result)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "redis.SetChapterListAllV1")
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
		Data:    result,
	})
}
