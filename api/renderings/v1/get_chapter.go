package v1Response

import (
	"encoding/json"

	db "fourleaves.studio/manga-scraper/internal/database/prisma"
)

type ChapterData struct {
	Provider    string     `json:"provider"`
	Series      string     `json:"series"`
	Slug        string     `json:"slug"`
	FullTitle   string     `json:"fullTitle"`
	ShortTitle  string     `json:"shortTitle"`
	Number      float64    `json:"number"`
	SourceURL   string     `json:"sourceURL"`
	ChapterNav  ChapterNav `json:"chapterNav"`
	ContentURLs []string   `json:"contentURLs"`
}

type ChapterNav struct {
	NextSlug string `json:"nextSlug"`
	NextURL  string `json:"nextURL"`
	PrevSlug string `json:"prevSlug"`
	PrevURL  string `json:"prevURL"`
}

func NewChapterData(provider *db.ProviderModel, series string, chapter *db.ChapterModel) ChapterData {
	var sourceURL string

	sourcePath := chapter.SourcePath
	if sourcePath != "" {
		sourceURL = provider.Scheme + provider.Host + sourcePath
	}

	var nextSlug, nextURL, prevSlug, prevURL string

	if chapter.NextSlug != "" {
		nextSlug = chapter.NextSlug
		if chapter.NextPath != "" {
			nextURL = provider.Scheme + provider.Host + chapter.NextPath
		}
	}

	if chapter.PrevSlug != "" {
		prevSlug = chapter.PrevSlug
		if chapter.PrevPath != "" {
			prevURL = provider.Scheme + provider.Host + chapter.PrevPath
		}
	}

	var contentPaths []string
	err := json.Unmarshal(chapter.ContentPaths, &contentPaths)
	if err != nil {
		contentPaths = []string{}
	}

	var contentURLs []string
	for i := range contentPaths {
		contentURLs = append(contentURLs, provider.Scheme+provider.Host+contentPaths[i])
	}

	return ChapterData{
		Provider:   provider.Slug,
		Series:     series,
		Slug:       chapter.Slug,
		FullTitle:  chapter.FullTitle,
		ShortTitle: chapter.ShortTitle,
		Number:     chapter.Number,
		SourceURL:  sourceURL,
		ChapterNav: ChapterNav{
			NextSlug: nextSlug,
			NextURL:  nextURL,
			PrevSlug: prevSlug,
			PrevURL:  prevURL,
		},
		ContentURLs: contentURLs,
	}
}

func NewChapterListData(series *db.SeriesModel) []ChapterData {
	provider := series.Provider()
	chapterList := series.Chapters()

	result := make([]ChapterData, 0, len(chapterList))
	for _, chapter := range chapterList {
		result = append(result, NewChapterData(provider, series.Slug, &chapter))
	}

	return result
}
