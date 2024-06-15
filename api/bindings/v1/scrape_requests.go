package v1Binding

type PostScrapeSeriesList struct {
	Provider string `json:"provider" validate:"required" example:"asura"`
} // @name PostScrapeSeriesList

type PostScrapeChapterList struct {
	Provider string `json:"provider" validate:"required" example:"asura"`
	Series   string `json:"series" validate:"required" example:"reincarnator"`
} // @name PostScrapeChapterList

type PutScrapeSeriesDetail struct {
	Provider string `json:"provider" validate:"required" example:"asura"`
	Series   string `json:"series" validate:"required" example:"reincarnator"`
} // @name PutScrapeSeriesDetail

type PutScrapeChapterDetail struct {
	Provider string `json:"provider" validate:"required" example:"asura"`
	Series   string `json:"series" validate:"required" example:"reincarnator"`
	Chapter  string `json:"chapter" validate:"required" example:"reincarnator-chapter-27"`
} // @name PutScrapeChapterDetail
