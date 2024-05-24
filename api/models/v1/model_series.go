package v1Model

type SeriesList struct {
	Title      string
	Slug       string
	SourcePath string
}

type SeriesDetail struct {
	ThumbnailURL string
	Synopsis     string
	Genres       []byte
}
