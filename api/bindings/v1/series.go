package v1Binding

type PaginatedRequest struct {
	Sort string `query:"sort" validate:"omitempty,oneof=asc desc" example:"asc"`
	Page int    `query:"page" validate:"required,gt=0" example:"1"`
	Size int    `query:"size" validate:"required,gt=0,lte=100" example:"10"`
}
