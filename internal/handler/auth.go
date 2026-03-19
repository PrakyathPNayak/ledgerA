package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/service"

	firebasepkg "ledgerA/pkg/firebase"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles auth endpoints.
type AuthHandler struct {
	userService service.UserService
}

// NewAuthHandler creates AuthHandler.
func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Sync handles POST /api/v1/auth/sync.
func (h *AuthHandler) Sync(c *gin.Context) {
	var req dto.SyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 401, "ERR_UNAUTHORIZED", "missing or invalid token")
		return
	}

	verified, err := firebasepkg.VerifyIDToken(c.Request.Context(), req.FirebaseToken)
	if err != nil {
		dto.Error(c, 401, "ERR_UNAUTHORIZED", "token verification failed")
		return
	}

	user, err := h.userService.Sync(c.Request.Context(), verified.UID, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to sync user")
		return
	}

	dto.Success(c, dto.SyncResponse{User: dto.NewUserResponse(*user)})
}
