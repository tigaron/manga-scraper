package middlewares

import (
	"net/http"

	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		span := sentry.StartSpan(c.Request().Context(), "IsAuthenticated")
		span.Name = "IsAuthenticated"
		defer span.Finish()

		sess, err := session.Get("session", c)
		if err != nil {
			SentryHandleInternalError(c, span, err, "session.Get")
			return c.JSON(http.StatusInternalServerError, v1Response.Response{
				Error:   true,
				Message: "Internal Server Error",
			})
		}

		if sess.Values["profile"] == nil {
			span.Status = sentry.SpanStatusUnauthenticated
			return c.JSON(http.StatusUnauthorized, v1Response.Response{
				Error:   true,
				Message: "Unauthorized",
				Detail:  "Please log in at /api/v1/users/login",
			})
		}

		span.Status = sentry.SpanStatusOK
		return next(c)
	}
}

func IsAdmin(adminSub string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			span := sentry.StartSpan(c.Request().Context(), "IsAdmin")
			span.Name = "IsAdmin"
			defer span.Finish()

			sess, err := session.Get("session", c)
			if err != nil {
				SentryHandleInternalError(c, span, err, "session.Get")
				return c.JSON(http.StatusInternalServerError, v1Response.Response{
					Error:   true,
					Message: "Internal Server Error",
				})
			}

			if sess.Values["profile"] == nil {
				span.Status = sentry.SpanStatusUnauthenticated
				return c.JSON(http.StatusUnauthorized, v1Response.Response{
					Error:   true,
					Message: "Unauthorized",
					Detail:  "Please log in at /api/v1/users/login",
				})
			}

			profile := sess.Values["profile"].(map[string]interface{})
			if profile["sub"] != adminSub {
				span.Status = sentry.SpanStatusPermissionDenied
				return c.JSON(http.StatusForbidden, v1Response.Response{
					Error:   true,
					Message: "Forbidden",
				})
			}

			span.Status = sentry.SpanStatusOK
			return next(c)
		}
	}
}
