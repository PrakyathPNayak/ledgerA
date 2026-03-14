package dto

// PaginationQuery contains generic pagination fields.
type PaginationQuery struct {
	Page    int `form:"page" validate:"omitempty,min=1"`
	PerPage int `form:"per_page" validate:"omitempty,min=1,max=200"`
}

// PaginationMeta contains pagination metadata.
type PaginationMeta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

// PaginatedResponse wraps data with pagination metadata.
type PaginatedResponse[T any] struct {
	Data T              `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
