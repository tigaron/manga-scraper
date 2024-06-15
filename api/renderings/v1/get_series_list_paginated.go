package v1Response

import db "fourleaves.studio/manga-scraper/internal/database/prisma"

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
