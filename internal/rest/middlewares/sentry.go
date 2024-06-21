package middlewares

import (
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (m *Middleware) SentryMiddleware() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll:   false,
		DisablePrintStack: true,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			// Capture exception to Sentry
			sentry.CurrentHub().WithScope(func(scope *sentry.Scope) {
				scope.SetRequest(c.Request())
				scope.SetContext("StackTrace", map[string]interface{}{
					"trace": string(stack),
				})
				sentry.CaptureException(err)
			})
			return err
		},
	})
}
