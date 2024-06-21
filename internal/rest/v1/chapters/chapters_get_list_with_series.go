package chapters

import (
	"net/http"

	"fourleaves.studio/manga-scraper/internal"
	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get chapter list with series
// @Description	Get chapter list with series
// @Tags			chapters
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug"	example(asura)
// @Param			series_slug		path		string	true	"Series slug"	example(reincarnator)
// @Success		200				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/chapters/{provider_slug}/{series_slug}/_list [get]
func (h *ChapterHandler) FindListWithRel(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.FindListWithRel")
	defer span.Finish()

	providerSlug := c.Param("provider_slug")
	seriesSlug := c.Param("series_slug")

	params := internal.FindChapterParams{
		Provider: providerSlug,
		Series:   seriesSlug,
	}

	chaptersList, err := h.svc.FindListWithRel(c.Request().Context(), params)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to get chapters", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Handler.Response{
		Error:   false,
		Message: "OK",
		Data:    chaptersList,
	})
}
