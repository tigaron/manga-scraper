package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

func (m *Middleware) WithHeaderAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		span := sentry.StartSpan(c.Request().Context(), "WithHeaderAuth")
		span.Name = "WithHeaderAuth"
		defer span.Finish()

		authHeader := c.Request().Header.Get("Authorization")
		sessionToken := strings.TrimPrefix(authHeader, "Bearer ")

		c.Logger().Debugj(map[string]interface{}{
			"_source":       "middlewares.WithHeaderAuth",
			"session_token": sessionToken,
		})

		claims, err := jwt.Verify(c.Request().Context(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			span.Status = sentry.SpanStatusUnauthenticated
			return c.JSON(http.StatusUnauthorized, v1Handler.Response{
				Error:   true,
				Message: "Unauthorized",
				Detail:  "Invalid session token",
			})
		}

		profile, err := json.Marshal(claims)
		if err != nil {
			// SentryHandleInternalError(c, span, err, "json.Marshal")
			return c.JSON(http.StatusInternalServerError, v1Handler.Response{
				Error:   true,
				Message: "Internal Server Error",
			})
		}

		c.Logger().Debugj(map[string]interface{}{
			"_source": "middlewares.WithHeaderAuth",
			"profile": string(profile),
		})

		span.Status = sentry.SpanStatusOK
		return next(c)
	}
}

func (m *Middleware) IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		span := sentry.StartSpan(c.Request().Context(), "IsAdmin")
		span.Name = "IsAdmin"
		defer span.Finish()

		authHeader := c.Request().Header.Get("Authorization")
		sessionToken := strings.TrimPrefix(authHeader, "Bearer ")

		c.Logger().Debugj(map[string]interface{}{
			"_source":       "middlewares.IsAdmin",
			"session_token": sessionToken,
		})

		claims, err := jwt.Verify(c.Request().Context(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			span.Status = sentry.SpanStatusUnauthenticated
			return c.JSON(http.StatusUnauthorized, v1Handler.Response{
				Error:   true,
				Message: "Unauthorized",
				Detail:  "Invalid session token",
			})
		}

		profile, err := json.Marshal(claims)
		if err != nil {
			// SentryHandleInternalError(c, span, err, "json.Marshal")
			return c.JSON(http.StatusInternalServerError, v1Handler.Response{
				Error:   true,
				Message: "Internal Server Error",
			})
		}

		c.Logger().Debugj(map[string]interface{}{
			"_source": "middlewares.IsAdmin",
			"profile": string(profile),
		})

		if claims.Subject != m.config.AdminSub {
			span.Status = sentry.SpanStatusPermissionDenied
			return c.JSON(http.StatusForbidden, v1Handler.Response{
				Error:   true,
				Message: "Forbidden",
				Detail:  "Admin role required",
			})
		}

		span.Status = sentry.SpanStatusOK
		return next(c)
	}
}
