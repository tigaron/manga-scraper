package v1Response

import db "fourleaves.studio/manga-scraper/internal/database/prisma"

type ListChapterData struct {
	Provider   string  `json:"provider"`
	Series     string  `json:"series"`
	Slug       string  `json:"slug"`
	ShortTitle string  `json:"shortTitle"`
	Number     float64 `json:"number"`
}

type ListChapterResult struct {
	Series   SeriesData        `json:"series"`
	Chapters []ListChapterData `json:"chapters"`
}

func NewListChapterData(provider string, series string, chapter db.ChapterModel) ListChapterData {
	return ListChapterData{
		Provider:   provider,
		Series:     series,
		Slug:       chapter.Slug,
		ShortTitle: chapter.ShortTitle,
		Number:     chapter.Number,
	}
}

func NewListAllChapterData(series *db.SeriesModel) ListChapterResult {
	provider := series.Provider()
	chapterList := series.Chapters()

	result := make([]ListChapterData, 0, len(chapterList))
	for _, chapter := range chapterList {
		result = append(result, NewListChapterData(provider.Slug, series.Slug, chapter))
	}

	return ListChapterResult{
		Series:   NewSeriesData(provider, series),
		Chapters: result,
	}
}
