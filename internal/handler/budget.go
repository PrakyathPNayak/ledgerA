package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BudgetHandler handles budget endpoints.
type BudgetHandler struct {
	budgetService service.BudgetService
	userService   service.UserService
}

// NewBudgetHandler creates BudgetHandler.
func NewBudgetHandler(budgetService service.BudgetService, userService service.UserService) *BudgetHandler {
	return &BudgetHandler{budgetService: budgetService, userService: userService}
}

func (h *BudgetHandler) resolveUserID(c *gin.Context) (uuid.UUID, bool) {
	uid, err := middleware.GetFirebaseUID(c)
	if err != nil {
		dto.Error(c, 401, "ERR_UNAUTHORIZED", "unauthorized")
		return uuid.Nil, false
	}
	user, err := h.userService.GetMe(c.Request.Context(), uid)
	if err != nil {
		dto.Error(c, 404, "ERR_NOT_FOUND", "user not found")
		return uuid.Nil, false
	}
	return user.ID, true
}

// List handles GET /api/v1/budgets.
func (h *BudgetHandler) List(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	items, total, err := h.budgetService.List(c.Request.Context(), userID)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to list budgets")
		return
	}
	responses := make([]dto.BudgetResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, dto.NewBudgetResponse(item))
	}
	dto.Paginated(c, responses, 1, len(responses), total)
}

// Create handles POST /api/v1/budgets.
func (h *BudgetHandler) Create(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var req dto.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 400, "ERR_VALIDATION", err.Error())
		return
	}
	item, err := h.budgetService.Create(c.Request.Context(), userID, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Created(c, dto.NewBudgetResponse(*item))
}

// Update handles PATCH /api/v1/budgets/:id.
func (h *BudgetHandler) Update(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_VALIDATION", "invalid id")
		return
	}
	var req dto.UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 400, "ERR_VALIDATION", err.Error())
		return
	}
	item, err := h.budgetService.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Success(c, dto.NewBudgetResponse(*item))
}

// Delete handles DELETE /api/v1/budgets/:id.
func (h *BudgetHandler) Delete(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_VALIDATION", "invalid id")
		return
	}
	if err := h.budgetService.Delete(c.Request.Context(), userID, id); err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Success(c, gin.H{"deleted": true})
}

// Progress handles GET /api/v1/budgets/progress.
func (h *BudgetHandler) Progress(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	results, err := h.budgetService.Progress(c.Request.Context(), userID)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Success(c, results)
}
