package scraper

import (
	"fmt"

	v1Model "fourleaves.studio/manga-scraper/api/models/v1"
	"fourleaves.studio/manga-scraper/internal/scraper/asura"
)

func ScrapeSeriesList(browserUrl, provider, url string) ([]v1Model.SeriesList, error) {
	switch provider {
	case "asura":
		return asura.ScrapeSeriesList(browserUrl, url)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func ScrapeSeriesDetail(browserUrl, provider, url string) (v1Model.SeriesDetail, error) {
	switch provider {
	case "asura":
		return asura.ScrapeSeriesDetail(browserUrl, url)
	default:
		return v1Model.SeriesDetail{}, fmt.Errorf("not implemented yet")
	}
}

func ScrapeChapterList(browserUrl, provider, url string) ([]v1Model.ChapterList, error) {
	switch provider {
	case "asura":
		return asura.ScrapeChapterList(browserUrl, url)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func ScrapeChapterDetail(browserUrl, provider, url string) (v1Model.ChapterDetail, error) {
	switch provider {
	case "asura":
		return asura.ScrapeChapterDetail(browserUrl, url)
	default:
		return v1Model.ChapterDetail{}, fmt.Errorf("not implemented yet")
	}
}
