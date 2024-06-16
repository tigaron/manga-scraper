package v1Handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"fourleaves.studio/manga-scraper/api/middlewares"
	v1Response "fourleaves.studio/manga-scraper/api/renderings/v1"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

type SearchHitsTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type SearchHitsSource struct {
	Title         string   `json:"title"`
	Synopsis      string   `json:"synopsis"`
	Genres        []string `json:"genres"`
	CoverURL      string   `json:"coverURL"`
	ChaptersCount int      `json:"chaptersCount"`
	LatestChapter string   `json:"latestChapter"`
}

type SearchHitsData struct {
	Index  string           `json:"_index"`
	ID     string           `json:"_id"`
	Score  float64          `json:"_score"`
	Source SearchHitsSource `json:"_source"`
}

type SearchHits struct {
	Total    SearchHitsTotal  `json:"total"`
	MaxScore float64          `json:"max_score"`
	Hits     []SearchHitsData `json:"hits"`
}

type SearchShards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type SearchResponse struct {
	Took     int          `json:"took"`
	TimedOut bool         `json:"timed_out"`
	Shards   SearchShards `json:"_shards"`
	Hits     SearchHits   `json:"hits"`
}

type SearchResult struct {
	Provider      string   `json:"provider"`
	Slug          string   `json:"slug"`
	Title         string   `json:"title"`
	Synopsis      string   `json:"synopsis"`
	Genres        []string `json:"genres"`
	CoverURL      string   `json:"coverURL"`
	ChaptersCount int      `json:"chaptersCount"`
	LatestChapter string   `json:"latestChapter"`
}

// @Summary		Get series search result
// @Description	Get series search result
// @Tags			search
// @Produce		json
// @Param			q	query		string	true	"Query"	example(high school)
// @Success		200	{object}	ResponseV1
// @Failure		400	{object}	ResponseV1
// @Failure		404	{object}	ResponseV1
// @Failure		500	{object}	ResponseV1
// @Router			/api/v1/search [get]
func (h *Handler) GetSearch(c echo.Context) error {
	span := sentry.StartSpan(c.Request().Context(), "v1.GetSearch")
	span.Name = "v1.GetSearch"
	defer span.Finish()

	q := c.QueryParam("q")

	if q == "" {
		span.Status = sentry.SpanStatusInvalidArgument
		return c.JSON(http.StatusBadRequest, v1Response.Response{
			Error:   true,
			Message: "Bad Request",
			Detail:  "Query is required",
		})
	}

	c.Logger().Debugj(map[string]interface{}{
		"_source": "v1.GetSearch",
		"query":   q,
	})

	encodedQ := url.QueryEscape(q)
	url := h.config.SearchURL + "_search" + "?q=" + encodedQ
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "http.NewRequest", q)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to create request",
		})
	}

	resp, err := client.Do(req)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "client.Do", q)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to send request",
		})
	}

	var res SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		middlewares.SentryHandleInternalErrorWithData(c, span, err, "json.NewDecoder", q)
		return c.JSON(http.StatusInternalServerError, v1Response.Response{
			Error:   true,
			Message: "Internal Server Error",
			Detail:  "Failed to decode response",
		})
	}

	if len(res.Hits.Hits) == 0 {
		span.Status = sentry.SpanStatusNotFound
		return c.JSON(http.StatusNotFound, v1Response.Response{
			Error:   true,
			Message: "Not Found",
			Detail:  "No result found",
		})
	}

	var result []SearchResult

	for _, hit := range res.Hits.Hits {
		result = append(result, SearchResult{
			Provider:      hit.Index,
			Slug:          hit.ID,
			Title:         hit.Source.Title,
			Synopsis:      hit.Source.Synopsis,
			Genres:        hit.Source.Genres,
			CoverURL:      hit.Source.CoverURL,
			ChaptersCount: hit.Source.ChaptersCount,
			LatestChapter: hit.Source.LatestChapter,
		})
	}

	span.Status = sentry.SpanStatusOK
	return c.JSON(http.StatusOK, v1Response.Response{
		Error:   false,
		Message: "OK",
		Data:    result,
	})
}
