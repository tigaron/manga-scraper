package internal

type Series struct {
	Provider      string   `json:"provider"`
	Slug          string   `json:"slug"`
	Title         string   `json:"title"`
	SourceURL     string   `json:"sourceURL"`
	CoverURL      string   `json:"coverURL"`
	Synopsis      string   `json:"synopsis"`
	Genres        []string `json:"genres"`
	ChaptersCount int      `json:"chaptersCount"`
	LatestChapter string   `json:"latestChapter"`
}
type SeriesBC struct {
	Provider Breadcrumb `json:"provider"`
	Series   Breadcrumb `json:"series"`
}

type CreateInitSeriesParams struct {
	Provider   string
	Slug       string
	Title      string
	SourcePath string
}

type FindSeriesParams struct {
	Provider string
	Slug     string
	Order    SortOrder
	Page     int
	Size     int
	Cursor   string
}

const (
	ASC  SortOrder = "asc"
	DESC SortOrder = "desc"
)

type SortOrder string

type UpdateInitSeriesParams struct {
	Provider     string
	Slug         string
	ThumbnailURL string
	Synopsis     string
	Genres       []byte
}

type UpdateLatestSeriesParams struct {
	Provider      string
	Slug          string
	AddChapters   int
	LatestChapter string
}

func NewSortOrder(s string) SortOrder {
	switch s {
	case "asc":
		return ASC
	case "desc":
		return DESC
	default:
		return ASC
	}
}

func (s *CreateInitSeriesParams) Validate() error {
	if s.Provider == "" {
		return NewErrorf(ErrInvalidInput, "provider is required")
	}

	if s.Slug == "" {
		return NewErrorf(ErrInvalidInput, "slug is required")
	}

	if s.Title == "" {
		return NewErrorf(ErrInvalidInput, "title is required")
	}

	if s.SourcePath == "" {
		return NewErrorf(ErrInvalidInput, "source path is required")
	}

	return nil
}

func (s *UpdateInitSeriesParams) Validate() error {
	if s.Provider == "" {
		return NewErrorf(ErrInvalidInput, "provider is required")
	}

	if s.Slug == "" {
		return NewErrorf(ErrInvalidInput, "slug is required")
	}

	if s.ThumbnailURL == "" {
		return NewErrorf(ErrInvalidInput, "thumbnail URL is required")
	}

	return nil
}

func (s *UpdateLatestSeriesParams) Validate() error {
	if s.Provider == "" {
		return NewErrorf(ErrInvalidInput, "provider is required")
	}

	if s.Slug == "" {
		return NewErrorf(ErrInvalidInput, "slug is required")
	}

	if s.AddChapters == 0 {
		return NewErrorf(ErrInvalidInput, "add chapters is required")
	}

	if s.LatestChapter == "" {
		return NewErrorf(ErrInvalidInput, "latest chapter is required")
	}

	return nil
}
