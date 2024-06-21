package v1Handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary		Get health check
// @Description	Get health check
// @Produce		json
// @Success		200	{object}	ResponseV1
// @Router			/health [get]
func GetHealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		Error:   false,
		Message: "OK",
	})
}
