package chapters

// import (
// 	"errors"
// 	"net/http"

// 	"fourleaves.studio/manga-scraper/api/middlewares"
// 	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
// 	db "fourleaves.studio/manga-scraper/internal/database/prisma"
// 	"github.com/getsentry/sentry-go"
// 	"github.com/labstack/echo/v4"
// )

// // @Summary		Get chapter list only
// // @Description	Get chapter list only
// // @Tags			chapters
// // @Produce		json
// // @Param			provider_slug	path		string	true	"Provider slug"	example(asura)
// // @Param			series_slug		path		string	true	"Series slug"	example(reincarnator)
// // @Param			sort			query		string	false	"Sort order"	enum(asc, desc)	default(asc)
// // @Success		200				{object}	ResponseV1
// // @Failure		404				{object}	ResponseV1
// // @Failure		500				{object}	ResponseV1
// // @Router			/api/v1/chapters/{provider_slug}/{series_slug}/_list [get]
// func (h *Handler) GetChapterList(c echo.Context) error {
// 	span := sentry.StartSpan(c.Request().Context(), "v1.GetChapterList")
// 	span.Name = "v1.GetChapterList"
// 	defer span.Finish()

// 	providerSlug := c.Param("provider_slug")
// 	seriesSlug := c.Param("series_slug")
// 	sort := c.QueryParam("sort")

// 	var order db.SortOrder

// 	switch sort {
// 	case "asc":
// 		order = db.ASC
// 	case "desc":
// 		order = db.DESC
// 	default:
// 		order = db.ASC
// 	}

// 	c.Logger().Debugj(map[string]interface{}{
// 		"_source":       "v1.GetChapterList",
// 		"provider_slug": providerSlug,
// 		"series_slug":   seriesSlug,
// 	})

// 	cache, err := h.redis.GetChaptersListWithRelV1(c.Request().Context(), providerSlug, seriesSlug, order)
// 	if err == nil {
// 		span.Status = sentry.SpanStatusOK
// 		return c.JSON(http.StatusOK, v1Handler.Response{
// 			Error:   false,
// 			Message: "OK",
// 			Data:    cache,
// 		})
// 	}

// 	series, err := h.prisma.FindChaptersListWithRelV1(c.Request().Context(), providerSlug, seriesSlug, order)
// 	if errors.Is(err, db.ErrNotFound) {
// 		span.Status = sentry.SpanStatusNotFound
// 		return c.JSON(http.StatusNotFound, v1Handler.Response{
// 			Error:   true,
// 			Message: "Not found",
// 			Detail:  "Series not found",
// 		})
// 	} else if err != nil {
// 		middlewares.SentryHandleInternalError(c, span, err, "prisma.FindChaptersListWithRelV1")
// 		return c.JSON(http.StatusInternalServerError, v1Handler.Response{
// 			Error:   true,
// 			Message: "Internal Server Error",
// 			Detail:  "Failed to find chapters",
// 		})
// 	}

// 	chapterList := series.Chapters()

// 	if len(chapterList) == 0 {
// 		span.Status = sentry.SpanStatusNotFound
// 		return c.JSON(http.StatusNotFound, v1Handler.Response{
// 			Error:   true,
// 			Message: "Not found",
// 			Detail:  "No chapters found",
// 		})
// 	}

// 	result := v1Response.NewListAllChapterData(series)

// 	err = h.redis.SetChaptersListWithRelV1(c.Request().Context(), providerSlug, seriesSlug, order, result)
// 	if err != nil {
// 		middlewares.SentryHandleInternalError(c, span, err, "redis.SetChaptersListWithRelV1")
// 	}

// 	span.Status = sentry.SpanStatusOK
// 	return c.JSON(http.StatusOK, v1Handler.Response{
// 		Error:   false,
// 		Message: "OK",
// 		Data:    result,
// 	})
// }
