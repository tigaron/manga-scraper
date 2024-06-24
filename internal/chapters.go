package internal

type Chapter struct {
	Provider    string      `json:"provider"`
	Series      string      `json:"series"`
	Slug        string      `json:"slug"`
	Number      float64     `json:"number"`
	FullTitle   string      `json:"fullTitle,omitempty"`
	ShortTitle  string      `json:"shortTitle"`
	SourceURL   string      `json:"sourceURL,omitempty"`
	ChapterNav  *ChapterNav `json:"chapterNav,omitempty"`
	ContentURLs []string    `json:"contentURLs,omitempty"`
	SourceHref  string      `json:"sourceHref,omitempty"`
}

type ChapterList struct {
	Series   Series    `json:"series"`
	Chapters []Chapter `json:"chapters"`
}

type ChapterNav struct {
	NextSlug string `json:"nextSlug,omitempty"`
	NextURL  string `json:"nextURL,omitempty"`
	PrevSlug string `json:"prevSlug,omitempty"`
	PrevURL  string `json:"prevURL,omitempty"`
}

type Breadcrumb struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type ChapterBC struct {
	Provider Breadcrumb `json:"provider"`
	Series   Breadcrumb `json:"series"`
	Chapter  Breadcrumb `json:"chapter"`
}

type CreateInitChapterParams struct {
	Provider   string
	Series     string
	Slug       string
	Number     float64
	ShortTitle string
	SourceHref string
}

type FindChapterParams struct {
	Provider string
	Series   string
	Slug     string
	Order    SortOrder
	Page     int
	Size     int
	Cursor   string
}

type UpdateInitChapterParams struct {
	Provider     string
	Series       string
	Slug         string
	FullTitle    string
	SourcePath   string
	ContentPaths []byte
	NextSlug     string
	NextPath     string
	PrevSlug     string
	PrevPath     string
}

func (c *CreateInitChapterParams) Validate() error {
	if c.Provider == "" {
		return NewErrorf(ErrInvalidInput, "provider is required")
	}

	if c.Series == "" {
		return NewErrorf(ErrInvalidInput, "series is required")
	}

	if c.Slug == "" {
		return NewErrorf(ErrInvalidInput, "slug is required")
	}

	if c.ShortTitle == "" {
		return NewErrorf(ErrInvalidInput, "shortTitle is required")
	}

	if c.SourceHref == "" {
		return NewErrorf(ErrInvalidInput, "sourceHref is required")
	}

	return nil
}

func (u *UpdateInitChapterParams) Validate() error {
	if u.Provider == "" {
		return NewErrorf(ErrInvalidInput, "provider is required")
	}

	if u.Series == "" {
		return NewErrorf(ErrInvalidInput, "series is required")
	}

	if u.Slug == "" {
		return NewErrorf(ErrInvalidInput, "slug is required")
	}

	if u.FullTitle == "" {
		return NewErrorf(ErrInvalidInput, "fullTitle is required")
	}

	if u.SourcePath == "" {
		return NewErrorf(ErrInvalidInput, "sourcePath is required")
	}

	return nil
}

func CcreateValidCreateInitChapterParams() *CreateInitChapterParams {
	return &CreateInitChapterParams{
		Provider:   "validProvider",
		Series:     "validSeries",
		Slug:       "validSlug",
		Number:     1,
		ShortTitle: "validShortTitle",
		SourceHref: "validSourceHref",
	}
}

func CreateValidUpdateInitChapterParams() *UpdateInitChapterParams {
	return &UpdateInitChapterParams{
		Provider:   "validProvider",
		Series:     "validSeries",
		Slug:       "validSlug",
		FullTitle:  "validFullTitle",
		SourcePath: "validSourcePath",
	}
}
