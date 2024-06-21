package providersHandler

import (
	"net/http"

	"fourleaves.studio/manga-scraper/internal"
	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Get provider list
// @Description	Get provider list
// @Tags			providers
// @Produce		json
// @Success		200	{object}	ResponseV1
// @Failure		404	{object}	ResponseV1
// @Failure		500	{object}	ResponseV1
// @Router			/api/v1/providers [get]
func (h *ProviderHandler) FindAll(c echo.Context) error {
	span := newSentrySpan(c.Request().Context(), "v1.FindAll")
	defer span.Finish()

	providers, err := h.svc.FindAll(c.Request().Context(), internal.ASC)
	if err != nil {
		return v1Handler.RenderErrorResponse(c, "Failed to get providers", err, span)
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Handler.Response{
		Error:   false,
		Message: "OK",
		Data:    providers,
	})
}
