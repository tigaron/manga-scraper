package middlewares

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SentryMiddleware() echo.MiddlewareFunc {
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

func SentryHandleInternalError(c echo.Context, span *sentry.Span, err error, transaction string) {
	c.Logger().Debugj(map[string]interface{}{
		"_source": transaction,
		"error":   err.Error(),
	})

	span.Status = sentry.SpanStatusInternalError
	sentry.CurrentHub().WithScope(func(scope *sentry.Scope) {
		scope.SetRequest(c.Request())
		scope.SetTag("handled", "true")
		scope.SetTag("transaction", transaction)
		scope.SetTag("status_code", fmt.Sprintf("%d", http.StatusInternalServerError))
		sentry.CaptureException(err)
	})
}

func SentryHandleInternalErrorWithData(c echo.Context, span *sentry.Span, err error, transaction string, data interface{}) {
	c.Logger().Debugj(map[string]interface{}{
		"_source": transaction,
		"error":   err.Error(),
	})

	sentry.CurrentHub().WithScope(func(scope *sentry.Scope) {
		scope.SetRequest(c.Request())
		scope.SetTag("handled", "true")
		scope.SetTag("transaction", transaction)
		scope.SetTag("status_code", fmt.Sprintf("%d", http.StatusInternalServerError))
		scope.SetExtra("data", data)
		sentry.CaptureException(err)
	})
}
