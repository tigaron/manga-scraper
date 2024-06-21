package scrapers

import (
	"net/http"

	"fourleaves.studio/manga-scraper/internal"
	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Create scrape request
// @Description	Create scrape request
// @Security		TokenAuth
// @Tags			scrapers
// @Accept			json
// @Produce		json
// @Param			body	body		CreateScrapeRequest	true	"Request body"
// @Success		200		{object}	ResponseV1
// @Failure		400		{object}	ResponseV1
// @Failure		401		{object}	ResponseV1
// @Failure		403		{object}	ResponseV1
// @Failure		404		{object}	ResponseV1
// @Failure		500		{object}	ResponseV1
// @Router			/api/v1/scrapers [post]
func (h *ScraperHandler) Create(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.Create")
	defer span.Finish()

	var req CreateScrapeRequest
	err := c.Bind(&req)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Invalid request", internal.WrapErrorf(err, internal.ErrInvalidInput, "bind request"), span)
	}

	err = c.Validate(&req)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Invalid request", internal.WrapErrorf(err, internal.ErrInvalidInput, "validate request"), span)
	}

	params := internal.CreateScrapeRequestParams{
		Type:        internal.ScrapeRequestType(req.Type),
		Status:      internal.ScrapeRequestStatus(req.Status),
		BaseURL:     req.BaseURL,
		RequestPath: req.RequestPath,
		Provider:    req.Provider,
		Series:      req.Series,
		Chapter:     req.Chapter,
	}

	scrapeRequest, err := h.svc.Create(c.Request().Context(), params)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to create scrape request", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusAccepted, v1Handler.Response{
		Error:   false,
		Message: "Accepted",
		Data:    scrapeRequest,
	})
}
