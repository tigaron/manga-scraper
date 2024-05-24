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

func NewChapterData(provider *db.ProviderModel, series *db.SeriesModel, chapter *db.ChapterModel) ChapterData {
	fullTitle, ok := chapter.FullTitle()
	if !ok {
		fullTitle = ""
	}

	var sourceURL string

	sourcePath, ok := chapter.SourcePath()
	if !ok {
		sourceURL = ""
	} else {
		sourceURL = provider.Scheme + provider.Host + sourcePath
	}

	var nextSlug, nextURL, prevSlug, prevURL string

	if nextChapter, ok := chapter.NextSlug(); ok {
		nextSlug = nextChapter
		if nextPath, ok := chapter.NextPath(); ok && nextPath != "" {
			nextURL = provider.Scheme + provider.Host + nextPath
		}
	}

	if prevChapter, ok := chapter.PrevSlug(); ok {
		prevSlug = prevChapter
		if prevPath, ok := chapter.PrevPath(); ok && prevPath != "" {
			prevURL = provider.Scheme + provider.Host + prevPath
		}
	}

	contentPathsJSON, ok := chapter.ContentPaths()
	if !ok {
		contentPathsJSON = []byte(`[]`)
	}

	var contentPaths []string
	err := json.Unmarshal([]byte(contentPathsJSON), &contentPaths)
	if err != nil {
		contentPaths = []string{}
	}

	var contentURLs []string
	for i := range contentPaths {
		contentURLs = append(contentURLs, provider.Scheme+provider.Host+contentPaths[i])
	}

	return ChapterData{
		Provider:   series.ProviderSlug,
		Series:     series.Slug,
		Slug:       chapter.Slug,
		FullTitle:  fullTitle,
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

func NewChapterListData(provider *db.ProviderModel, series *db.SeriesModel, chapterList []db.ChapterModel) []ChapterData {
	result := make([]ChapterData, 0, len(chapterList))
	for _, chapter := range chapterList {
		result = append(result, NewChapterData(provider, series, &chapter))
	}

	return result
}
