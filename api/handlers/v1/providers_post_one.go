package v1Handler

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"

	v1Binding "fourleaves.studio/manga-scraper/api/bindings/v1"
	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

// @Summary		Create provider
// @Description	Create provider
// @Security		TokenAuth
// @Tags			providers
// @Accept			json
// @Produce		json
// @Param			body	body		PostProviderRequest	true	"Request body"
// @Success		201		{object}	ResponseV1
// @Failure		400		{object}	ResponseV1
// @Failure		401		{object}	ResponseV1
// @Failure		403		{object}	ResponseV1
// @Failure		409		{object}	ResponseV1
// @Failure		500		{object}	ResponseV1
// @Router			/api/v1/providers [post]
func (h *Handler) PostProvider(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.PostProvider")
	span.Name = "v1.PostProvider"
	defer span.Finish()

	var req v1Binding.PostProviderRequest
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

	c.Logger().Debugj(map[string]interface{}{
		"_source":   "v1.PostProvider",
		"slug":      req.Slug,
		"name":      req.Name,
		"scheme":    req.Scheme,
		"host":      req.Host,
		"list_path": req.ListPath,
		"is_active": req.IsActive,
	})

	provider, err := h.prisma.CreateProviderV1(c.Request().Context(), req)
	if err != nil {
		if info, e := db.IsErrUniqueConstraint(err); e {
			span.Status = sentry.SpanStatusAlreadyExists
			return c.JSON(http.StatusConflict, v1Response.Response{
				Error:   true,
				Message: "Conflict",
				Detail:  fmt.Sprintf("Unique constraint violation on %v", info.Key),
			})
		}

		middlewares.SentryHandleInternalErrorWithData(c, span, err, "prisma.CreateProviderV1", req)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to create provider",
		})
	}

	err = h.redis.UnsetProviderListV1(c.Request().Context())
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "redis.UnsetProviderListV1", req)
	}

	result := v1Response.NewProviderData(provider)

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusCreated, v1Response.Response{
		Error:   false,
		Message: "Created",
		Data:    result,
	})
}
