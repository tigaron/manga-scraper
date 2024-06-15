package v1Response

import (
	"encoding/json"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

type SeriesSearchData struct {
	Slug string     `json:"slug"`
	Data SearchData `json:"data"`
}

type SearchData struct {
	Title    string   `json:"title"`
	Synopsis string   `json:"synopsis"`
	Genres   []string `json:"genres"`
	CoverURL string   `json:"coverURL"`
}

func NewSeriesSearchData(provider *db.ProviderModel, series *db.SeriesModel) SeriesSearchData {
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

	return SeriesSearchData{
		Slug: series.Slug,
		Data: SearchData{
			Title:    series.Title,
			Synopsis: synopsis,
			Genres:   genres,
			CoverURL: thumbnailUrl,
		},
	}
}

func NewSeriesListSearchData(provider *db.ProviderModel, seriesList []db.SeriesModel) []SeriesSearchData {
	result := make([]SeriesSearchData, 0, len(seriesList))
	for _, series := range seriesList {
		result = append(result, NewSeriesSearchData(provider, &series))
	}

	return result
}