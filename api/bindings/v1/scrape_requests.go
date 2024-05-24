package v1Binding

type PostScrapeSeries struct {
	Provider string `json:"provider" validate:"required" example:"asura"`
}
