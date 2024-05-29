package scraper

import (
	"context"
	"fmt"

	v1Model "fourleaves.studio/manga-scraper/api/models/v1"
	"fourleaves.studio/manga-scraper/internal/scraper/asura"
)

func ScrapeSeriesList(ctx context.Context, browserUrl, provider, url string) ([]v1Model.SeriesList, error) {
	switch provider {
	case "asura":
		return asura.ScrapeSeriesList(ctx, browserUrl, url)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func ScrapeSeriesDetail(ctx context.Context, browserUrl, provider, url string) (v1Model.SeriesDetail, error) {
	switch provider {
	case "asura":
		return asura.ScrapeSeriesDetail(ctx, browserUrl, url)
	default:
		return v1Model.SeriesDetail{}, fmt.Errorf("not implemented yet")
	}
}

func ScrapeChapterList(ctx context.Context, browserUrl, provider, url string) ([]v1Model.ChapterList, error) {
	switch provider {
	case "asura":
		return asura.ScrapeChapterList(ctx, browserUrl, url)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func ScrapeChapterDetail(ctx context.Context, browserUrl, provider, url string) (v1Model.ChapterDetail, error) {
	switch provider {
	case "asura":
		return asura.ScrapeChapterDetail(ctx, browserUrl, url)
	default:
		return v1Model.ChapterDetail{}, fmt.Errorf("not implemented yet")
	}
}
