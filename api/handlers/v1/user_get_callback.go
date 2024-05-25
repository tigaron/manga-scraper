package v1Handler

import (
	"net/http"

	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Handler for our callback.
func (h *Handler) GetCallback(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetCallback")
	span.Name = "v1.GetCallback"
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

	if c.QueryParam("state") != sess.Values["state"] {
		span.Status = sentry.SpanStatusInvalidArgument
		return c.JSON(http.StatusBadRequest, v1Response.Response{
			Error:   true,
			Message: "Bad Request",
			Detail:  "Invalid state parameter.",
		})
	}

	// Exchange an authorization code for a token.
	token, err := h.auth.Exchange(c.Request().Context(), c.QueryParam("code"))
	if err != nil {
		span.Status = sentry.SpanStatusPermissionDenied
		return c.JSON(http.StatusUnauthorized, v1Response.Response{
			Error:   true,
			Message: "Unauthorized",
			Detail:  "Failed to exchange an authorization code for a token.",
		})
	}

	idToken, err := h.auth.VerifyIDToken(c.Request().Context(), token)
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "session.Get")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to verify ID Token.",
		})
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "idToken.Claims")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to get profile from ID Token.",
		})
	}

	sess.Values["access_token"] = token.AccessToken
	sess.Values["profile"] = profile
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "sess.Save")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to save session.",
		})
	}

	// Redirect to logged in page.
	span.Status = sentry.SpanStatusOK
	return c.Redirect(http.StatusTemporaryRedirect, "/api/v1/user/profile")
}
