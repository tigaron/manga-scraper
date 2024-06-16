package v1Response

import (
	"encoding/json"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

type SeriesData struct {
	Provider  string   `json:"provider"`
	Slug      string   `json:"slug"`
	Title     string   `json:"title"`
	SourceURL string   `json:"sourceURL"`
	CoverURL  string   `json:"coverURL"`
	Synopsis  string   `json:"synopsis"`
	Genres    []string `json:"genres"`
}

func NewSeriesData(provider *db.ProviderModel, series *db.SeriesModel) SeriesData {
	var genres []string
	err := json.Unmarshal(series.Genres, &genres)
	if err != nil {
		genres = []string{}
	}

	return SeriesData{
		Provider:  series.ProviderSlug,
		Slug:      series.Slug,
		Title:     series.Title,
		SourceURL: provider.Scheme + provider.Host + series.SourcePath,
		CoverURL:  series.ThumbnailURL,
		Synopsis:  series.Synopsis,
		Genres:    genres,
	}
}

func NewSeriesListData(provider *db.ProviderModel, seriesList []db.SeriesModel) []SeriesData {
	result := make([]SeriesData, 0, len(seriesList))
	for _, series := range seriesList {
		result = append(result, NewSeriesData(provider, &series))
	}

	return result
}
