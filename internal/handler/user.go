package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user endpoints.
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates UserHandler.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Me handles GET /api/v1/users/me.
func (h *UserHandler) Me(c *gin.Context) {
	uid, err := middleware.GetFirebaseUID(c)
	if err != nil {
		dto.Error(c, 401, "ERR_UNAUTHORIZED", err.Error())
		return
	}

	user, err := h.userService.GetMe(c.Request.Context(), uid)
	if err != nil {
		dto.Error(c, 404, "ERR_NOT_FOUND", err.Error())
		return
	}

	dto.Success(c, dto.NewUserResponse(*user))
}

// UpdateMe handles PATCH /api/v1/users/me.
func (h *UserHandler) UpdateMe(c *gin.Context) {
	uid, err := middleware.GetFirebaseUID(c)
	if err != nil {
		dto.Error(c, 401, "ERR_UNAUTHORIZED", err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}

	user, err := h.userService.UpdateMe(c.Request.Context(), uid, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}

	dto.Success(c, dto.NewUserResponse(*user))
}
