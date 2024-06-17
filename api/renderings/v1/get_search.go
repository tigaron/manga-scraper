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
	Provider      string   `json:"provider"`
	Title         string   `json:"title"`
	Synopsis      string   `json:"synopsis"`
	Genres        []string `json:"genres"`
	CoverURL      string   `json:"coverURL"`
	Status        string   `json:"status"`
	ChaptersCount int      `json:"chaptersCount"`
	LatestChapter string   `json:"latestChapter"`
}

func NewSeriesSearchData(provider *db.ProviderModel, series *db.SeriesModel) SeriesSearchData {
	var genres []string
	err := json.Unmarshal(series.Genres, &genres)
	if err != nil {
		genres = []string{}
	}

	return SeriesSearchData{
		Slug: series.Slug,
		Data: SearchData{
			Provider:      provider.Name,
			Title:         series.Title,
			Synopsis:      series.Synopsis,
			Genres:        genres,
			Status:        string(series.Status),
			CoverURL:      series.ThumbnailURL,
			ChaptersCount: series.ChaptersCount,
			LatestChapter: series.LatestChapter,
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
