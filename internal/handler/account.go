package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AccountHandler handles account endpoints.
type AccountHandler struct {
	accountService service.AccountService
	userService    service.UserService
}

// NewAccountHandler creates AccountHandler.
func NewAccountHandler(accountService service.AccountService, userService service.UserService) *AccountHandler {
	return &AccountHandler{accountService: accountService, userService: userService}
}

func (h *AccountHandler) resolveUserID(c *gin.Context) (uuid.UUID, bool) {
	uid, err := middleware.GetFirebaseUID(c)
	if err != nil {
		dto.Error(c, 401, "ERR_UNAUTHORIZED", err.Error())
		return uuid.Nil, false
	}
	user, err := h.userService.GetMe(c.Request.Context(), uid)
	if err != nil {
		dto.Error(c, 404, "ERR_NOT_FOUND", "user not found")
		return uuid.Nil, false
	}
	return user.ID, true
}

// List handles GET /api/v1/accounts.
func (h *AccountHandler) List(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	items, total, err := h.accountService.List(c.Request.Context(), userID)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	responses := make([]dto.AccountResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, dto.NewAccountResponse(item))
	}
	dto.Paginated(c, responses, 1, len(responses), total)
}

// Create handles POST /api/v1/accounts.
func (h *AccountHandler) Create(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	item, err := h.accountService.Create(c.Request.Context(), userID, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Created(c, dto.NewAccountResponse(*item))
}

// Update handles PATCH /api/v1/accounts/:id.
func (h *AccountHandler) Update(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid account id")
		return
	}
	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	item, err := h.accountService.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Success(c, dto.NewAccountResponse(*item))
}

// Delete handles DELETE /api/v1/accounts/:id.
func (h *AccountHandler) Delete(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid account id")
		return
	}
	if err := h.accountService.Delete(c.Request.Context(), userID, id); err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	c.Status(204)
}
