package scrapers

// import (
// 	"errors"
// 	"net/http"
// 	"strings"
// 	"time"

// 	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
// 	"fourleaves.studio/manga-scraper/api/middlewares"
// 	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
// 	db "fourleaves.studio/manga-scraper/internal/database/prisma"
// 	"fourleaves.studio/manga-scraper/internal/scraper"
// 	"github.com/getsentry/sentry-go"
// 	"github.com/labstack/echo/v4"
// )

// // @Summary		Create request to scrape chapter detail
// // @Description	Create request to scrape chapter detail
// // @Security		TokenAuth
// // @Tags			scrape-requests
// // @Accept			json
// // @Produce		json
// // @Param			body	body		PutScrapeChapterDetail	true	"Request body"
// // @Success		200		{object}	ResponseV1
// // @Failure		400		{object}	ResponseV1
// // @Failure		401		{object}	ResponseV1
// // @Failure		403		{object}	ResponseV1
// // @Failure		404		{object}	ResponseV1
// // @Failure		500		{object}	ResponseV1
// // @Failure		503		{object}	ResponseV1
// // @Router			/api/v1/scrape-requests/chapters/detail [put]
// func (h *Handler) PutScrapeChapterDetail(c echo.Context) error {
// 	span := sentry.StartSpan(c.Request().Context(), "v1.PutScrapeChapterDetail")
// 	span.Name = "v1.PutScrapeChapterDetail"
// 	defer span.Finish()

// 	var req v1Binding.PutScrapeChapterDetail
// 	err := c.Bind(&req)
// 	if err != nil {
// 		span.Status = sentry.SpanStatusInvalidArgument
// 		return c.JSON(http.StatusBadRequest, v1Response.Response{
// 			Error:   true,
// 			Message: "Bad Request",
// 			Detail:  err.Error(),
// 		})
// 	}

// 	err = c.Validate(&req)
// 	if err != nil {
// 		span.Status = sentry.SpanStatusInvalidArgument
// 		return c.JSON(http.StatusBadRequest, v1Response.Response{
// 			Error:   true,
// 			Message: "Bad Request",
// 			Detail:  middlewares.FormatValidationError(err),
// 		})
// 	}

// 	c.Logger().Debugj(map[string]interface{}{
// 		"_source":  "v1.PutScrapeChapterDetail",
// 		"provider": req.Provider,
// 		"series":   req.Series,
// 		"chapter":  req.Chapter,
// 	})

// 	chapter, err := h.prisma.FindChapterUniqueWithRelV1(c.Request().Context(), req.Provider, req.Series, req.Chapter)
// 	if errors.Is(err, db.ErrNotFound) {
// 		span.Status = sentry.SpanStatusNotFound
// 		return c.JSON(http.StatusNotFound, v1Response.Response{
// 			Error:   true,
// 			Message: "Not found",
// 			Detail:  "Chapter not found",
// 		})
// 	} else if err != nil {
// 		middlewares.SentryHandleInternalError(c, span, err, "prisma.FindChapterUniqueWithRelV1")
// 		return c.JSON(http.StatusInternalServerError, v1Response.Response{
// 			Error:   true,
// 			Message: "Internal Server Error",
// 			Detail:  "Failed to find chapter",
// 		})
// 	}

// 	provider := chapter.Provider()
// 	series := chapter.Series()

// 	receipt, err := h.prisma.CreateChapterDetailScrapeRequestV1(c.Request().Context(), provider, series, chapter)
// 	if err != nil {
// 		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateChapterDetailScrapeRequestV1", req)
// 		return c.JSON(http.StatusInternalServerError, v1Response.Response{
// 			Error:   true,
// 			Message: "Internal Server Error",
// 			Detail:  "Failed to create scrape request",
// 		})
// 	}

// 	startTime := time.Now()

// 	var reqURL string

// 	if chapter.SourcePath != "" {
// 		reqURL = provider.Scheme + provider.Host + chapter.SourcePath
// 	} else {
// 		hrefArr := strings.Split(chapter.SourceHref, "/")
// 		reqURL = provider.Scheme + provider.Host + "/" + strings.Join(hrefArr[3:], "/")
// 	}

// 	scrapeData, err := scraper.ScrapeChapterDetail(c.Request().Context(), h.config.RodURL, req.Provider, reqURL)
// 	if err != nil {
// 		middlewares.SentryHandleInternalErrorWithData(c, span, err, "scraper.ScrapeChapterList", req)
// 		h.prisma.UpdateScrapeRequestUniqueV1(c.Request().Context(), receipt.ID, "failed", time.Since(startTime).Seconds(), err.Error())
// 		return c.JSON(http.StatusInternalServerError, v1Response.Response{
// 			Error:   true,
// 			Message: "Internal Server Error",
// 			Detail:  "Failed to scrape chapter detail",
// 		})
// 	}

// 	_, err = h.prisma.UpdateScrapeRequestUniqueV1(c.Request().Context(), receipt.ID, "success", time.Since(startTime).Seconds(), "Completed successfully")
// 	if err != nil {
// 		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.UpdateScrapeRequestUniqueV1", req)
// 		return c.JSON(http.StatusInternalServerError, v1Response.Response{
// 			Error:   true,
// 			Message: "Internal Server Error",
// 			Detail:  "Failed to update scrape request",
// 		})
// 	}

// 	_, err = h.prisma.CreateChapterDetailRowV1(c.Request().Context(), receipt.ID, scrapeData)
// 	if err != nil {
// 		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateChapterDetailRowV1", req)
// 	}

// 	updatedChapter, err := h.prisma.UpdateDetailChapterRowV1(c.Request().Context(), req.Provider, req.Series, req.Chapter, scrapeData)
// 	if err != nil {
// 		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.UpdateDetailChapterRowV1", req)
// 	}

// 	err = h.redis.UnsetChapterV1(c.Request().Context(), req.Provider, req.Series, chapter.Slug)
// 	if err != nil {
// 		middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetChapterV1", req)
// 	}

// 	err = h.redis.UnsetChapterListV1(c.Request().Context(), req.Provider, req.Series)
// 	if err != nil {
// 		middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetChapterListV1", req)
// 	}

// 	result := v1Response.NewChapterData(provider, series.Slug, updatedChapter)

// 	span.Status = sentry.SpanStatusOK
// 	return c.JSON(http.StatusOK, v1Response.Response{
// 		Error:   false,
// 		Message: "OK",
// 		Data:    result,
// 	})
// }
