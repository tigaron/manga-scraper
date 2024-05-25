package v1Handler

import (
	"net/http"
	"net/url"

	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Handler for our logout.
func (h *Handler) GetLogout(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetLogout")
	span.Name = "v1.GetLogout"
	defer span.Finish()

	// Get the session
	sess, err := session.Get("session", c)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "session.Get")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to get session.",
		})
	}

	// Invalidate the session
	sess.Options.MaxAge = -1
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "sess.Save")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to save session.",
		})
	}

	logoutUrl, err := url.Parse("https://" + h.config.Auth0Domain + "/v2/logout")
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "url.Parse")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to parse logout URL.",
		})
	}

	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + c.Request().Host + "/api/v1")
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "url.Parse")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to parse returnTo URL.",
		})
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", h.config.Auth0ClientID)
	logoutUrl.RawQuery = parameters.Encode()

	span.Status = sentry.SpanStatusOK
	return c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}
