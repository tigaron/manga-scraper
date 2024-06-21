package prisma

import (
	"context"
	"encoding/json"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/getsentry/sentry-go"
)

func newSentrySpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	span.Name = "fourleaves.studio/manga-scraper/internal/database/prisma"
	span.SetTag("db.system", "mysql")

	return span
}

func newSortOrder(order internal.SortOrder) SortOrder {
	switch order {
	case internal.ASC:
		return SortOrderAsc
	case internal.DESC:
		return SortOrderDesc
	default:
		return SortOrderAsc
	}
}

func newStringSliceFromBytes(b []byte) []string {
	var s []string
	_ = json.Unmarshal(b, &s)

	return s
}

func newContentURLsFromSlice(s []string, p string) []string {
	var paths []string
	for i := range s {
		paths = append(paths, p+s[i])
	}

	return paths
}
