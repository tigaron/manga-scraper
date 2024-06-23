package providers

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"

	"fourleaves.studio/manga-scraper/internal"
	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
)

// @Summary		Update provider
// @Description	Update provider
// @Security		TokenAuth
// @Tags			providers
// @Accept			json
// @Produce		json
// @Param			provider_slug	path		string					true	"Provider slug"	example(asura)
// @Param			body			body		UpdateProviderRequest	true	"Request body"
// @Success		200				{object}	ResponseV1
// @Failure		400				{object}	ResponseV1
// @Failure		401				{object}	ResponseV1
// @Failure		403				{object}	ResponseV1
// @Failure		404				{object}	ResponseV1
// @Failure		500				{object}	ResponseV1
// @Router			/api/v1/providers/{provider_slug} [put]
func (h *ProviderHandler) Update(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.Update")
	defer span.Finish()

	var req UpdateProviderRequest
	err := c.Bind(&req)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Invalid request", internal.WrapErrorf(err, internal.ErrInvalidInput, "bind request"), span)
	}

	err = c.Validate(&req)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Invalid request", internal.WrapErrorf(err, internal.ErrInvalidInput, "validate request"), span)
	}

	providerSlug := c.Param("provider_slug")

	params := internal.ProviderParams{
		Slug:     providerSlug,
		Name:     req.Name,
		Scheme:   req.Scheme,
		Host:     req.Host,
		ListPath: req.ListPath,
		IsActive: req.IsActive,
	}

	provider, err := h.svc.Update(c.Request().Context(), params)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to update provider", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Handler.Response{
		Error:   false,
		Message: "Updated",
		Data:    provider,
	})
}
