package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/getsentry/sentry-go"
	"github.com/opensearch-project/opensearch-go/v2"
	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type SeriesSearchRepository struct {
	client *opensearch.Client
}

func NewSeriesSearchRepository(client *opensearch.Client) *SeriesSearchRepository {
	return &SeriesSearchRepository{
		client: client,
	}
}

func (s *SeriesSearchRepository) Index(ctx context.Context, series internal.Series) error {
	defer newSentrySpan(ctx, "SeriesSearchRepository.Index").Finish()

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(series); err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "json.NewEncoder.Encode")
	}

	req := opensearchapi.IndexRequest{
		Index:      series.Provider,
		Body:       &buf,
		DocumentID: series.Slug,
		Refresh:    "true",
	}

	resp, err := req.Do(ctx, s.client)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "IndexRequest.Do")
	}

	defer resp.Body.Close()

	if resp.IsError() {
		return internal.NewErrorf(internal.ErrUnknown, "IndexRequest.Do %d", resp.StatusCode)
	}

	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}

func (s *SeriesSearchRepository) Delete(ctx context.Context, provider, slug string) error {
	defer newSentrySpan(ctx, "SeriesSearchRepository.Delete").Finish()

	req := opensearchapi.DeleteRequest{
		Index:      provider,
		DocumentID: slug,
	}

	resp, err := req.Do(ctx, s.client)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "DeleteRequest.Do")
	}

	defer resp.Body.Close()

	if resp.IsError() {
		return internal.NewErrorf(internal.ErrUnknown, "DeleteRequest.Do %d", resp.StatusCode)
	}

	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}

func (s *SeriesSearchRepository) Search(ctx context.Context, q string) ([]internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesSearchRepository.Search").Finish()

	var buf bytes.Buffer

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  q,
				"fields": []string{"title", "synopsis", "genres"},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		fmt.Println(err)
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "json.NewEncoder.Encode")
	}

	req := opensearchapi.SearchRequest{
		Index: []string{},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, s.client)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "SearchRequest.Do")
	}

	defer resp.Body.Close()

	if resp.IsError() {
		return nil, internal.NewErrorf(internal.ErrUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source internal.Series `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "json.NewDecoder.Decode")
	}

	if len(hits.Hits.Hits) == 0 {
		return nil, internal.NewErrorf(internal.ErrNotFound, "no results found")
	}

	result := make([]internal.Series, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		result[i] = hit.Source
	}

	return result, nil
}

func newSentrySpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	span.Name = "fourleaves.studio/manga-scraper/internal/elasticsearch"
	span.SetTag("db.system", "elasticsearch")

	return span
}
