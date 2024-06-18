package v1Handler

import (
	"fmt"
	"net/http"

	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get all chapter list
// @Description	Get all chapter list
// @Tags			chapters
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug"	example(asura)
// @Param			series_slug		path		string	true	"Series slug"	example(reincarnator)
// @Param			sort			query		string	false	"Sort order"	enum(asc, desc)	default(asc)
// @Success		200				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/chapters/{provider_slug}/{series_slug}/_all [get]
func (h *Handler) GetChapterListAll(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetChapterListAll")
	span.Name = "v1.GetChapterListAll"
	defer span.Finish()

	providerSlug := c.Param("provider_slug")
	seriesSlug := c.Param("series_slug")
	sort := c.QueryParam("sort")

	c.Logger().Debugj(map[string]interface{}{
		"_source":       "v1.GetChapterListAll",
		"provider_slug": providerSlug,
		"series_slug":   seriesSlug,
		"sort":          sort,
	})

	var order db.SortOrder

	switch sort {
	case "asc":
		order = db.ASC
	case "desc":
		order = db.DESC
	default:
		order = db.ASC
	}

	cache, err := h.redis.GetChaptersListAllV1(c.Request().Context(), providerSlug, seriesSlug, order)
	if err == nil {
		span.Status = sentry.SpanStatusOK
		return c.JSON(http.StatusOK, v1Response.Response{
			Error:   false,
			Message: "OK",
			Data:    cache,
		})
	}

	series, err := h.prisma.FindChaptersListAllV1(c.Request().Context(), providerSlug, seriesSlug, order)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "prisma.FindChaptersListAllV1")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to find chapters",
		})
	}

	chapterList := series.Chapters()

	if len(chapterList) == 0 {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  fmt.Sprintf("Chapters with provider slug '%s' and series slug '%s' not found", providerSlug, seriesSlug),
		})
	}

	result := v1Response.NewChapterListData(series)

	err = h.redis.SetChapterListAllV1(c.Request().Context(), providerSlug, seriesSlug, result)
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
