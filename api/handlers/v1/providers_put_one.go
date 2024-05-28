package v1Handler

import (
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"

	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

// @Summary		Update provider
// @Description	Update provider
// @Tags			providers
// @Accept			json
// @Produce		json
// @Param			provider_slug	path	string							true	"Provider slug"	example(asura)
// @Param			body			body	v1Binding.PutProviderRequest	true	"Request body"
// @Security		BearerAuth
// @Success		200	{object}	v1Response.Response
// @Failure		400	{object}	v1Response.Response
// @Failure		401	{object}	v1Response.Response
// @Failure		403	{object}	v1Response.Response
// @Failure		404	{object}	v1Response.Response
// @Failure		500	{object}	v1Response.Response
// @Router			/api/v1/providers/{provider_slug} [put]
func (h *Handler) PutProvider(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.PutProvider")
	span.Name = "v1.PutProvider"
	defer span.Finish()

	var req v1Binding.PutProviderRequest
	err := c.Bind(&req)
	if err != nil {
		span.Status = sentry.SpanStatusInvalidArgument
		return c.JSON(http.StatusBadRequest, v1Response.Response{
			Error:   true,
			Message: "Bad Request",
			Detail:  err.Error(),
		})
	}

	err = c.Validate(&req)
	if err != nil {
		span.Status = sentry.SpanStatusInvalidArgument
		return c.JSON(http.StatusBadRequest, v1Response.Response{
			Error:   true,
			Message: "Bad Request",
			Detail:  middlewares.FormatValidationError(err),
		})
	}

	providerSlug := c.Param("provider_slug")

	provider, err := h.prisma.UpdateProviderUniqueV1(c.Request().Context(), providerSlug, req)
	if errors.Is(err, db.ErrNotFound) {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  "Provider not found",
		})
	} else if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.UpdateProviderUniqueV1", req)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to update provider",
		})
	}

	err = h.redis.UnsetProviderV1(c.Request().Context(), providerSlug)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetProviderV1", req)
	}

	err = h.redis.UnsetProviderListV1(c.Request().Context())
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetProviderListV1", req)
	}

	result := v1Response.NewProviderData(provider)

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "Updated",
		Data:    result,
	})
}
