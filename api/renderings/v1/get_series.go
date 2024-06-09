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

	return SeriesData{
		Provider:  series.ProviderSlug,
		Slug:      series.Slug,
		Title:     series.Title,
		SourceURL: provider.Scheme + provider.Host + series.SourcePath,
		CoverURL:  thumbnailUrl,
		Synopsis:  synopsis,
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

type PaginationData struct {
	PrevPage int `json:"prevPage,omitempty"`
	NextPage int `json:"nextPage,omitempty"`
	Total    int `json:"total,omitempty"`
}

type PaginatedSeriesData struct {
	PaginationData
	Series []SeriesData `json:"series"`
}

func NewSeriesListPaginatedData(provider *db.ProviderModel, seriesList []db.SeriesModel, paginationData PaginationData) PaginatedSeriesData {
	result := make([]SeriesData, 0, len(seriesList))
	for _, series := range seriesList {
		result = append(result, NewSeriesData(provider, &series))
	}

	return PaginatedSeriesData{
		PaginationData: paginationData,
		Series:         result,
	}
}
