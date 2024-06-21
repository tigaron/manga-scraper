package internal

type ScrapeRequest struct {
	ID          string              `json:"id"`
	Type        ScrapeRequestType   `json:"type"`
	Status      ScrapeRequestStatus `json:"status"`
	BaseURL     string              `json:"baseURL"`
	RequestPath string              `json:"requestPath"`
	Provider    string              `json:"provider"`
	Series      string              `json:"series,omitempty"`
	Chapter     string              `json:"chapter,omitempty"`
	Retries     int                 `json:"retries,omitempty"`
	TotalTime   float64             `json:"totalTime,omitempty"`
	Error       bool                `json:"error,omitempty"`
	Message     string              `json:"message,omitempty"`
}

type CreateScrapeRequestParams struct {
	ID          string
	Type        ScrapeRequestType   `json:"type"`
	Status      ScrapeRequestStatus `json:"status"`
	BaseURL     string              `json:"baseURL"`
	RequestPath string              `json:"requestPath"`
	Provider    string              `json:"provider"`
	Series      string              `json:"series,omitempty"`
	Chapter     string              `json:"chapter,omitempty"`
}

type UpdateScrapeRequestParams struct {
	ID        string
	Status    ScrapeRequestStatus
	Retries   int
	TotalTime float64
	Error     bool
	Message   string
}

type FindScrapeRequestParams struct {
	Status ScrapeRequestStatus
	Order  SortOrder
	Page   int
	Size   int
	Cursor string
}

type SeriesListResult struct {
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	SourcePath string `json:"sourcePath"`
}

type SeriesDetailResult struct {
	ThumbnailURL string `json:"thumbnailURL"`
	Synopsis     string `json:"synopsis"`
	Genres       []byte `json:"genres"`
}

type ChapterListResult struct {
	ShortTitle string  `json:"shortTitle"`
	Slug       string  `json:"slug"`
	Number     float64 `json:"number"`
	Href       string  `json:"href"`
}

type ChapterDetailResult struct {
	FullTitle    string `json:"fullTitle"`
	SourcePath   string `json:"sourcePath"`
	ContentPaths []byte `json:"contentPaths"`
	NextPath     string `json:"nextPath"`
	NextSlug     string `json:"nextSlug"`
	PrevPath     string `json:"prevPath"`
	PrevSlug     string `json:"prevSlug"`
}

type TSReaderScript struct {
	PrevURL string `json:"prevUrl"`
	NextURL string `json:"nextUrl"`
	Sources []struct {
		Images []string `json:"images"`
	} `json:"sources"`
}

type ScrapeRequestStatus string

const (
	PendingRequestStatus   ScrapeRequestStatus = "PENDING"
	CompletedRequestStatus ScrapeRequestStatus = "COMPLETED"
	FailedRequestStatus    ScrapeRequestStatus = "FAILED"
)

type ScrapeRequestType string

const (
	SeriesListRequestType    ScrapeRequestType = "SERIES_LIST"
	SeriesDetailRequestType  ScrapeRequestType = "SERIES_DETAIL"
	ChapterListRequestType   ScrapeRequestType = "CHAPTER_LIST"
	ChapterDetailRequestType ScrapeRequestType = "CHAPTER_DETAIL"
)

func (s *CreateScrapeRequestParams) Validate() error {
	if s.Type == "" {
		return NewErrorf(ErrInvalidInput, "type is required")
	}

	if s.BaseURL == "" {
		return NewErrorf(ErrInvalidInput, "baseURL is required")
	}

	if s.RequestPath == "" {
		return NewErrorf(ErrInvalidInput, "requestPath is required")
	}

	if s.Status == "" {
		return NewErrorf(ErrInvalidInput, "status is required")
	}

	if s.Provider == "" {
		return NewErrorf(ErrInvalidInput, "provider is required")
	}

	switch s.Type {
	case ChapterDetailRequestType:
		if s.Series == "" {
			return NewErrorf(ErrInvalidInput, "series is required")
		}

		if s.Chapter == "" {
			return NewErrorf(ErrInvalidInput, "chapter is required")
		}
	case ChapterListRequestType:
		fallthrough
	case SeriesDetailRequestType:
		if s.Series == "" {
			return NewErrorf(ErrInvalidInput, "series is required")
		}
	}

	return nil
}

func (s *UpdateScrapeRequestParams) Validate() error {
	if s.ID == "" {
		return NewErrorf(ErrInvalidInput, "id is required")
	}

	if s.Status == "" {
		return NewErrorf(ErrInvalidInput, "status is required")
	}

	return nil
}
