package scraper

import (
	"context"
	"fmt"

	"fourleaves.studio/manga-scraper/internal"
	"fourleaves.studio/manga-scraper/internal/scraper/agscomics"
	"fourleaves.studio/manga-scraper/internal/scraper/anigliscans"
	"fourleaves.studio/manga-scraper/internal/scraper/asura"
	"fourleaves.studio/manga-scraper/internal/scraper/flame"
	"fourleaves.studio/manga-scraper/internal/scraper/luminous"
	"fourleaves.studio/manga-scraper/internal/scraper/mangagalaxy"
	"fourleaves.studio/manga-scraper/internal/scraper/nightscans"
	"fourleaves.studio/manga-scraper/internal/scraper/surya"
)

func ScrapeSeriesList(ctx context.Context, browserUrl, provider, url string) ([]internal.SeriesListResult, error) {
	switch provider {
	case "asura":
		return asura.ScrapeSeriesList(ctx, browserUrl, url)
	case "surya":
		return surya.ScrapeSeriesList(ctx, browserUrl, url)
	case "flame":
		return flame.ScrapeSeriesList(ctx, browserUrl, url)
	case "luminous":
		return luminous.ScrapeSeriesList(ctx, browserUrl, url)
	case "anigliscans":
		return anigliscans.ScrapeSeriesList(ctx, browserUrl, url)
	case "agscomics":
		return agscomics.ScrapeSeriesList(ctx, browserUrl, url)
	case "nightscans":
		return nightscans.ScrapeSeriesList(ctx, browserUrl, url)
	case "mangagalaxy":
		return mangagalaxy.ScrapeSeriesList(ctx, browserUrl, url)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func ScrapeSeriesDetail(ctx context.Context, browserUrl, provider, url string) (internal.SeriesDetailResult, error) {
	switch provider {
	case "asura":
		return asura.ScrapeSeriesDetail(ctx, browserUrl, url)
	case "surya":
		return surya.ScrapeSeriesDetail(ctx, browserUrl, url)
	case "flame":
		return flame.ScrapeSeriesDetail(ctx, browserUrl, url)
	case "luminous":
		return luminous.ScrapeSeriesDetail(ctx, browserUrl, url)
	case "anigliscans":
		return anigliscans.ScrapeSeriesDetail(ctx, browserUrl, url)
	case "agscomics":
		return agscomics.ScrapeSeriesDetail(ctx, browserUrl, url)
	case "nightscans":
		return nightscans.ScrapeSeriesDetail(ctx, browserUrl, url)
	case "mangagalaxy":
		return mangagalaxy.ScrapeSeriesDetail(ctx, browserUrl, url)
	default:
		return internal.SeriesDetailResult{}, fmt.Errorf("not implemented yet")
	}
}

func ScrapeChapterList(ctx context.Context, browserUrl, provider, url string) ([]internal.ChapterListResult, error) {
	switch provider {
	case "asura":
		return asura.ScrapeChapterList(ctx, browserUrl, url)
	case "surya":
		return surya.ScrapeChapterList(ctx, browserUrl, url)
	case "flame":
		return flame.ScrapeChapterList(ctx, browserUrl, url)
	case "luminous":
		return luminous.ScrapeChapterList(ctx, browserUrl, url)
	case "anigliscans":
		return anigliscans.ScrapeChapterList(ctx, browserUrl, url)
	case "agscomics":
		return agscomics.ScrapeChapterList(ctx, browserUrl, url)
	case "nightscans":
		return nightscans.ScrapeChapterList(ctx, browserUrl, url)
	case "mangagalaxy":
		return mangagalaxy.ScrapeChapterList(ctx, browserUrl, url)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func ScrapeChapterDetail(ctx context.Context, browserUrl, provider, url string) (internal.ChapterDetailResult, error) {
	switch provider {
	case "asura":
		return asura.ScrapeChapterDetail(ctx, browserUrl, url)
	case "surya":
		return surya.ScrapeChapterDetail(ctx, browserUrl, url)
	case "flame":
		return flame.ScrapeChapterDetail(ctx, browserUrl, url)
	case "luminous":
		return luminous.ScrapeChapterDetail(ctx, browserUrl, url)
	case "anigliscans":
		return anigliscans.ScrapeChapterDetail(ctx, browserUrl, url)
	case "agscomics":
		return agscomics.ScrapeChapterDetail(ctx, browserUrl, url)
	case "nightscans":
		return nightscans.ScrapeChapterDetail(ctx, browserUrl, url)
	case "mangagalaxy":
		return mangagalaxy.ScrapeChapterDetail(ctx, browserUrl, url)
	default:
		return internal.ChapterDetailResult{}, fmt.Errorf("not implemented yet")
	}
}
