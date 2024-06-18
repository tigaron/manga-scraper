package v1Binding

type PostProviderRequest struct {
	Slug     string `json:"slug" validate:"required" example:"asura"`
	Name     string `json:"name" validate:"required" example:"Asura Scans"`
	Scheme   string `json:"scheme" validate:"required" example:"https://"`
	Host     string `json:"host" validate:"required" example:"asuratoon.com"`
	ListPath string `json:"list_path" validate:"required" example:"/manga/list-mode/"`
	IsActive *bool  `json:"is_active" validate:"required" example:"true"`
} // @name PostProviderRequest

type PutProviderRequest struct {
	Name     string `json:"name" validate:"required" example:"Asura Scans"`
	Scheme   string `json:"scheme" validate:"required" example:"https://"`
	Host     string `json:"host" validate:"required" example:"asuratoon.com"`
	ListPath string `json:"list_path" validate:"required" example:"/manga/list-mode/"`
	IsActive *bool  `json:"is_active" validate:"required" example:"true"`
} // @name PutProviderRequest
