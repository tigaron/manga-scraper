package v1Handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// @Summary		Feed the open search engine
// @Description	Feed the open search engine
// @Tags			search
// @Produce		json
// @Security		cookieAuth
// @Success		200	{object}	v1Response.Response
// @Failure		401	{object}	v1Response.Response
// @Failure		403	{object}	v1Response.Response
// @Failure		404	{object}	v1Response.Response
// @Failure		500	{object}	v1Response.Response
// @Router			/api/v1/search [put]
func (h *Handler) PutSearch(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.PutSearch")
	span.Name = "v1.PutSearch"
	defer span.Finish()

	providers, err := h.prisma.FindProviderManyV1(c.Request().Context())
	if err != nil {
		middlewares.SentryHandleInternalError(c, span, err, "db.Provider.FindMany")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to find providers",
		})
	}

	if len(providers) == 0 {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not found",
			Detail:  "No provider found",
		})
	}

	var errList []string

	for _, provider := range providers {
		series, _ := h.prisma.FindSeriesManyV1(c.Request().Context(), provider.Slug)
		result := v1Response.NewSeriesListSearchData(&provider, series)

		for _, res := range result {
			index := provider.Slug
			url := h.config.SearchURL + index + "/_doc/" + res.Slug
			body, _ := json.Marshal(res.Data)

			req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			_, err := client.Do(req)
			if err != nil {
				errList = append(errList, err.Error())
			}
		}
	}

	if len(errList) > 0 {
		middlewares.SentryHandleInternalError(c, span, fmt.Errorf(strings.Join(errList, ", ")), "client.Do")
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  strings.Join(errList, ", "),
		})
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
	})
}
