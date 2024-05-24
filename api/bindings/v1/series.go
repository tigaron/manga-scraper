package v1Binding

type PaginatedRequest struct {
	Page int `form:"page" validate:"required,gt=0" example:"1"`
	Size int `form:"size" validate:"required,gt=0,lte=100" example:"10"`
}
