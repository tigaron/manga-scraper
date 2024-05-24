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
