package v1Response

import (
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

type PaginatedChapterListData struct {
	PaginationData
	Chapters []ChapterData `json:"chapters"`
}

func NewChapterListPaginatedData(series *db.SeriesModel, paginationData PaginationData) PaginatedChapterListData {
	provider := series.Provider()
	chapterList := series.Chapters()

	result := make([]ChapterData, 0, len(chapterList))
	for _, chapter := range chapterList {
		result = append(result, NewChapterData(provider, series.Slug, &chapter))
	}

	return PaginatedChapterListData{
		PaginationData: paginationData,
		Chapters:       result,
	}
}
