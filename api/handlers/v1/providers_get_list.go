package v1Handler

import (
	"net/http"

	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get provider list
// @Description	Get provider list
// @Tags			providers
// @Produce		json
// @Success		200	{object}	v1Response.Response
// @Failure		404	{object}	v1Response.Response
// @Failure		500	{object}	v1Response.Response
// @Router			/api/v1/providers [get]
func (h *Handler) GetProvidersList(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetProvidersList")
	span.Name = "v1.GetProvidersList"
	defer span.Finish()

	cache, err := h.redis.GetProviderListV1(c.Request().Context())
	if err == nil {
		span.Status = sentry.SpanStatusOK
		return c.JSON(http.StatusOK, v1Response.Response{
			Error:   false,
			Message: "OK",
			Data:    cache,
		})
	}

	providers, err := h.prisma.FindProviderManyV1(c.Request().Context())
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "db.Provider.FindMany")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to find providers",
		})
	}

	if len(providers) == 0 {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  "No provider found",
		})
	}

	result := v1Response.NewProvidersListData(providers)

	err = h.redis.SetProviderListV1(c.Request().Context(), result)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "redis.SetProviderListV1")
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
		Data:    result,
	})
}
