package v1Handler

import (
	"errors"
	"fmt"
	"net/http"

	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get provider breadcrumbs
// @Description	Get provider breadcrumbs
// @Tags			breadcrumbs
// @Produce		json
// @Param			provider_slug	path		string	true	"Provider slug" example(asura)
// @Success		200	{object}	ResponseV1
// @Failure		400	{object}	ResponseV1
// @Failure		404	{object}	ResponseV1
// @Failure		500	{object}	ResponseV1
// @Router			/api/v1/providers/{provider_slug}/_bc [get]
func (h *Handler) GetProviderBreadcrumbs(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetBreadcrumbs")
	span.Name = "v1.GetSearch"
	defer span.Finish()

	providerSlug := c.Param("provider_slug")

	c.Logger().Debugj(map[string]interface{}{
		"_source":       "v1.GetProviderBreadcrumbs",
		"provider_slug": providerSlug,
	})

	cache, err := h.redis.GetProviderBreadcrumbsV1(c.Request().Context(), providerSlug)
	if err == nil {
		span.Status = sentry.SpanStatusOK
		return c.JSON(http.StatusOK, v1Response.Response{
			Error:   false,
			Message: "OK",
			Data:    cache,
		})
	}

	provider, err := h.prisma.FindProviderUniqueV1(c.Request().Context(), providerSlug)
	if errors.Is(err, db.ErrNotFound) {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  fmt.Sprintf("Provider with slug '%s' not found", providerSlug),
		})
	} else if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "prisma.FindProviderUniqueV1")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to find provider",
		})
	}

	result := v1Response.NewProviderBreadcrumbs(provider.Name)

	err = h.redis.SetProviderBreadcrumbsV1(c.Request().Context(), providerSlug, result)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "redis.SetProviderBreadcrumbsV1")
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
		Data:    result,
	})
}
