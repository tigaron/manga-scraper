package v1Response

import (
	"encoding/json"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

type GetSeriesData struct {
	Provider  string   `json:"provider"`
	Slug      string   `json:"slug"`
	Title     string   `json:"title"`
	SourceURL string   `json:"sourceURL"`
	CoverURL  string   `json:"coverURL"`
	Synopsis  string   `json:"synopsis"`
	Genres    []string `json:"genres"`
}

func NewGetSeriesData(provider *db.ProviderModel, series *db.SeriesModel) GetSeriesData {
	thumbnailUrl, ok := series.ThumbnailURL()
	if !ok {
		thumbnailUrl = ""
	}

	synopsis, ok := series.Synopsis()
	if !ok {
		synopsis = ""
	}

	genresJson, ok := series.Genres()
	if !ok {
		genresJson = []byte(`[]`)
	}

	var genres []string
	err := json.Unmarshal(genresJson, &genres)
	if err != nil {
		genres = []string{}
	}

	return GetSeriesData{
		Provider:  series.ProviderSlug,
		Slug:      series.Slug,
		Title:     series.Title,
		SourceURL: provider.Scheme + provider.Host + series.SourcePath,
		CoverURL:  thumbnailUrl,
		Synopsis:  synopsis,
		Genres:    genres,
	}
}

func NewGetSeriesListData(provider *db.ProviderModel, seriesList []db.SeriesModel) []GetSeriesData {
	result := make([]GetSeriesData, 0, len(seriesList))
	for _, series := range seriesList {
		result = append(result, NewGetSeriesData(provider, &series))
	}

	return result
}
