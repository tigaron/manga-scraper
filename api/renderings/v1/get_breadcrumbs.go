package v1Response

type BreadcrumbsData struct {
	Provider string `json:"provider"`
	Series   string `json:"series,omitempty"`
	Chapter  string `json:"chapter,omitempty"`
}

func NewChapterBreadcrumbs(provider, series, chapter string) BreadcrumbsData {
	return BreadcrumbsData{
		Provider: provider,
		Series:   series,
		Chapter:  chapter,
	}
}

func NewSeriesBreadcrumbs(provider, series string) BreadcrumbsData {
	return BreadcrumbsData{
		Provider: provider,
		Series:   series,
	}
}

func NewProviderBreadcrumbs(provider string) BreadcrumbsData {
	return BreadcrumbsData{
		Provider: provider,
	}
}
