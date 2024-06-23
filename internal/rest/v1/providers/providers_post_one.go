package providers

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"

	"fourleaves.studio/manga-scraper/internal"
	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
)

// @Summary		Create provider
// @Description	Create provider
// @Security		TokenAuth
// @Tags			providers
// @Accept			json
// @Produce		json
// @Param			body	body		CreateProviderRequest	true	"Request body"
// @Success		201		{object}	ResponseV1
// @Failure		400		{object}	ResponseV1
// @Failure		401		{object}	ResponseV1
// @Failure		403		{object}	ResponseV1
// @Failure		409		{object}	ResponseV1
// @Failure		500		{object}	ResponseV1
// @Router			/api/v1/providers [post]
func (h *ProviderHandler) Create(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.Create")
	defer span.Finish()

	var req CreateProviderRequest
	err := c.Bind(&req)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Invalid request", internal.WrapErrorf(err, internal.ErrInvalidInput, "bind request"), span)
	}

	err = c.Validate(&req)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Invalid request", internal.WrapErrorf(err, internal.ErrInvalidInput, "validate request"), span)
	}

	provider, err := h.svc.Create(c.Request().Context(), internal.ProviderParams(req))
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to create provider", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusCreated, v1Handler.Response{
		Error:   false,
		Message: "Created",
		Data:    provider,
	})
}
