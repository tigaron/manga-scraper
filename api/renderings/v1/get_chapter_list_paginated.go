package v1Response

import db "fourleaves.studio/manga-scraper/internal/database/prisma"

type PaginatedChapterListData struct {
	PaginationData
	Chapters []ChapterData `json:"chapters"`
}

func NewChapterListPaginatedData(provider *db.ProviderModel, series *db.SeriesModel, chapterList []db.ChapterModel, paginationData PaginationData) PaginatedChapterListData {
	result := make([]ChapterData, 0, len(chapterList))
	for _, chapter := range chapterList {
		result = append(result, NewChapterData(provider, series, &chapter))
	}

	return PaginatedChapterListData{
		PaginationData: paginationData,
		Chapters:       result,
	}
}
