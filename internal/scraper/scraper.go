package scraper

import (
	"context"
	"fmt"

	v1Model "fourleaves.studio/manga-scraper/api/models/v1"
	"fourleaves.studio/manga-scraper/internal/scraper/asura"
	"fourleaves.studio/manga-scraper/internal/scraper/flame"
	"fourleaves.studio/manga-scraper/internal/scraper/luminous"
	"fourleaves.studio/manga-scraper/internal/scraper/surya"
)

func ScrapeSeriesList(ctx context.Context, browserUrl, provider, url string) ([]v1Model.SeriesList, error) {
	switch provider {
	case "asura":
		return asura.ScrapeSeriesList(ctx, browserUrl, url)
	case "surya":
		return surya.ScrapeSeriesList(ctx, browserUrl, url)
	case "flame":
		return flame.ScrapeSeriesList(ctx, browserUrl, url)
	case "luminous":
		return luminous.ScrapeSeriesList(ctx, browserUrl, url)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func ScrapeSeriesDetail(ctx context.Context, browserUrl, provider, url string) (v1Model.SeriesDetail, error) {
	switch provider {
	case "asura":
		return asura.ScrapeSeriesDetail(ctx, browserUrl, url)
	case "surya":
		return surya.ScrapeSeriesDetail(ctx, browserUrl, url)
	case "flame":
		return flame.ScrapeSeriesDetail(ctx, browserUrl, url)
	case "luminous":
		return luminous.ScrapeSeriesDetail(ctx, browserUrl, url)
	default:
		return v1Model.SeriesDetail{}, fmt.Errorf("not implemented yet")
	}
}

func ScrapeChapterList(ctx context.Context, browserUrl, provider, url string) ([]v1Model.ChapterList, error) {
	switch provider {
	case "asura":
		return asura.ScrapeChapterList(ctx, browserUrl, url)
	case "surya":
		return surya.ScrapeChapterList(ctx, browserUrl, url)
	case "flame":
		return flame.ScrapeChapterList(ctx, browserUrl, url)
	case "luminous":
		return luminous.ScrapeChapterList(ctx, browserUrl, url)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func ScrapeChapterDetail(ctx context.Context, browserUrl, provider, url string) (v1Model.ChapterDetail, error) {
	switch provider {
	case "asura":
		return asura.ScrapeChapterDetail(ctx, browserUrl, url)
	case "surya":
		return surya.ScrapeChapterDetail(ctx, browserUrl, url)
	case "flame":
		return flame.ScrapeChapterDetail(ctx, browserUrl, url)
	case "luminous":
		return luminous.ScrapeChapterDetail(ctx, browserUrl, url)
	default:
		return v1Model.ChapterDetail{}, fmt.Errorf("not implemented yet")
	}
}
