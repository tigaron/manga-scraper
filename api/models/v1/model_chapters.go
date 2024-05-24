package v1Model

type ChapterList struct {
	ShortTitle string
	Slug       string
	Number     float64
	Href       string
}

type ChapterDetail struct {
	FullTitle    string
	SourcePath   string
	ContentPaths []byte
	NextPath     string
	NextSlug     string
	PrevPath     string
	PrevSlug     string
}

type TSReaderScript struct {
	PrevURL string `json:"prevUrl"`
	NextURL string `json:"nextUrl"`
	Sources []struct {
		Images []string `json:"images"`
	} `json:"sources"`
}
