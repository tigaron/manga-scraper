package chapters

import (
	"net/http"

	"fourleaves.studio/manga-scraper/internal"
	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get chapter by slug
// @Description	Get chapter by slug
// @Tags			chapters
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug"	example(asura)
// @Param			series_slug		path		string	true	"Series slug"	example(reincarnator)
// @Param			chapter_slug	path		string	true	"Chapter slug"	example(reincarnator-chapter-0)
// @Success		200				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/chapters/{provider_slug}/{series_slug}/{chapter_slug} [get]
func (h *ChapterHandler) Find(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.Find")
	defer span.Finish()

	providerSlug := c.Param("provider_slug")
	seriesSlug := c.Param("series_slug")
	chapterSlug := c.Param("chapter_slug")

	params := internal.FindChapterParams{
		Provider: providerSlug,
		Series:   seriesSlug,
		Slug:     chapterSlug,
	}

	chapter, err := h.svc.Find(c.Request().Context(), params)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to get chapter", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Handler.Response{
		Error:   false,
		Message: "OK",
		Data:    chapter,
	})
}
