package v1Binding

type PostScrapeSeriesList struct {
	Provider string `json:"provider" validate:"required" example:"asura"`
}

type PutScrapeSeriesDetail struct {
	Provider string `json:"provider" validate:"required" example:"asura"`
	Series   string `json:"series" validate:"required" example:"reincarnator"`
}
