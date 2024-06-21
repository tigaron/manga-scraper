package chapters

import (
	"net/http"

	"fourleaves.studio/manga-scraper/internal"
	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get paginated chapter list
// @Description	Get paginated chapter list
// @Tags			chapters
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug"	example(asura)
// @Param			series_slug		path		string	true	"Series slug"	example(reincarnator)
// @Param			sort			query		string	false	"Sort order"	enum(asc, desc)	default(asc)
// @Param			page			query		string	true	"Page"			example(10)
// @Param			size			query		string	true	"Size"			example(100)
// @Success		200				{object}	ResponseV1
// @Failure		400				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/chapters/{provider_slug}/{series_slug} [get]
func (h *ChapterHandler) FindPaginated(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.FindPaginated")
	defer span.Finish()

	var req PaginatedRequest
	err := c.Bind(&req)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Invalid request", internal.WrapErrorf(err, internal.ErrInvalidInput, "bind request"), span)
	}

	err = c.Validate(&req)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Invalid request", internal.WrapErrorf(err, internal.ErrInvalidInput, "validate request"), span)
	}

	providerSlug := c.Param("provider_slug")
	seriesSlug := c.Param("series_slug")

	params := internal.FindChapterParams{
		Provider: providerSlug,
		Series:   seriesSlug,
		Order:    internal.NewSortOrder(req.Sort),
		Page:     req.Page,
		Size:     req.Size,
	}

	chapters, err := h.svc.FindPaginated(c.Request().Context(), params)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to get chapters", err, span)
	}

	var prevPage, nextPage, total int

	if req.Page >= 2 {
		prevPage = req.Page - 1
	}

	if len(chapters) == req.Size {
		nextPage = req.Page + 1
	}

	total = len(chapters)

	result := PaginatedResponse{
		PaginationData: PaginationData{
			PrevPage: prevPage,
			NextPage: nextPage,
			Total:    total,
		},
		Chapters: chapters,
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Handler.Response{
		Error:   false,
		Message: "OK",
		Data:    result,
	})
}
