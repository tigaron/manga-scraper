package v1

import (
	"errors"
	"net/http"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Detail  interface{} `json:"detail,omitempty"`
} // @name ResponseV1

func RenderErrorResponse(c echo.Context, msg string, err error, span *sentry.Span) error {
	status := http.StatusInternalServerError
	response := Response{
		Error:   true,
		Message: "Internal Server Error",
		Detail:  msg,
	}

	var iErr *internal.Error
	for errors.As(err, &iErr) {
		err = iErr.Unwrap()
	}

	if iErr == nil {
		response.Detail = "Something went wrong"
		sentry.CaptureException(err)
		c.Logger().Errorj(map[string]interface{}{
			"_source": "renderErrorResponse",
			"_msg":    msg,
			"error":   err.Error(),
		})
	} else {
		c.Logger().Errorj(map[string]interface{}{
			"_source": "renderErrorResponse",
			"_msg":    msg,
			"error":   iErr.Error(),
		})
		switch iErr.Code() {
		case internal.ErrNotFound:
			status = http.StatusNotFound
			span.Status = sentry.SpanStatusNotFound
			response.Message = "Not Found"
		case internal.ErrInvalidInput:
			status = http.StatusBadRequest
			span.Status = sentry.SpanStatusInvalidArgument
			response.Message = "Bad Request"
		case internal.ErrUniqueConstraint:
			status = http.StatusConflict
			span.Status = sentry.SpanStatusAlreadyExists
			response.Message = "Conflict"
		case internal.ErrUnknown:
			fallthrough
		default:
			status = http.StatusInternalServerError
			span.Status = sentry.SpanStatusInternalError
		}
	}

	return c.JSON(status, response)
}
