package v1Handler

import (
	"net/http"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	"github.com/labstack/echo/v4"
)

// @Summary		Get health check
// @Description	Get health check
// @Produce		json
// @Success		200	{object}	ResponseV1
// @Router			/health [get]
func (h *Handler) GetHealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
	})
}
