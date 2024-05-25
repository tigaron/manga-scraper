package v1Handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetLogin(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetLogin")
	span.Name = "v1.GetLogin"
	defer span.Finish()

	state, err := generateRandomState()
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "generateRandomState")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to generate random state.",
		})
	}

	// Save the state inside the session.
	sess, err := session.Get("session", c)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "session.Get")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to get session.",
		})
	}

	sess.Values["state"] = state
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "sess.Save")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to save session.",
		})
	}

	span.Status = sentry.SpanStatusUnauthenticated
	return c.Redirect(http.StatusTemporaryRedirect, h.auth.AuthCodeURL(state))
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
