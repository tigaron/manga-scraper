package scraper

import (
	"fmt"

	v1Model "fourleaves.studio/manga-scraper/api/models/v1"
	"fourleaves.studio/manga-scraper/internal/scraper/asura"
)

func ScrapeSeriesList(provider, url string) ([]v1Model.SeriesList, error) {
	switch provider {
	case "asura":
		return asura.ScrapeSeriesList(url)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func ScrapeSeriesDetail(provider, url string) (v1Model.SeriesDetail, error) {
	switch provider {
	case "asura":
		return asura.ScrapeSeriesDetail(url)
	default:
		return v1Model.SeriesDetail{}, fmt.Errorf("not implemented yet")
	}
}

func ScrapeChapterList(provider, url string) ([]v1Model.ChapterList, error) {
	switch provider {
	case "asura":
		return asura.ScrapeChapterList(url)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func ScrapeChapterDetail(provider, url string) (v1Model.ChapterDetail, error) {
	switch provider {
	case "asura":
		return asura.ScrapeChapterDetail(url)
	default:
		return v1Model.ChapterDetail{}, fmt.Errorf("not implemented yet")
	}
}
