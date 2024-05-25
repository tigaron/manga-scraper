package v1Handler

import (
	"net/http"

	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetProfile(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetProfile")
	span.Name = "v1.GetProfile"
	defer span.Finish()

	sess, err := session.Get("session", c)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "session.Get")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to get session.",
		})
	}

	profile := sess.Values["profile"]

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
		Data:    profile,
	})
}
