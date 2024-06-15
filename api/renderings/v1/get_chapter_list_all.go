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

func NewListChapterData(provider *db.ProviderModel, series *db.SeriesModel, chapter *db.ChapterModel) ListChapterData {
	return ListChapterData{
		Provider:   series.ProviderSlug,
		Series:     series.Slug,
		Slug:       chapter.Slug,
		ShortTitle: chapter.ShortTitle,
		Number:     chapter.Number,
	}
}

func NewListAllChapterData(provider *db.ProviderModel, series *db.SeriesModel, chapterList []db.ChapterModel) ListChapterResult {
	result := make([]ListChapterData, 0, len(chapterList))
	for _, chapter := range chapterList {
		result = append(result, NewListChapterData(provider, series, &chapter))
	}

	return ListChapterResult{
		Series:   NewSeriesData(provider, series),
		Chapters: result,
	}
}
