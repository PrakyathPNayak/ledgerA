package dto

import "github.com/gin-gonic/gin"

// ErrorBody contains standardized error details.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SuccessEnvelope contains standardized success response data.
type SuccessEnvelope struct {
	Data any `json:"data"`
}

// ErrorEnvelope contains standardized error response payload.
type ErrorEnvelope struct {
	Error ErrorBody `json:"error"`
}

// Success writes a success envelope with status 200.
func Success(c *gin.Context, data any) {
	c.JSON(200, SuccessEnvelope{Data: data})
}

// Created writes a success envelope with status 201.
func Created(c *gin.Context, data any) {
	c.JSON(201, SuccessEnvelope{Data: data})
}

// Paginated writes a paginated success envelope with status 200.
func Paginated(c *gin.Context, data any, page int, perPage int, total int64) {
	c.JSON(200, gin.H{"data": data, "meta": PaginationMeta{Page: page, PerPage: perPage, Total: total}})
}

// Error writes a standardized error envelope.
func Error(c *gin.Context, status int, code string, message string) {
	c.JSON(status, ErrorEnvelope{Error: ErrorBody{Code: code, Message: message}})
}
