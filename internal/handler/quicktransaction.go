package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// QuickTransactionHandler handles quick transaction endpoints.
type QuickTransactionHandler struct {
	quickService service.QuickTransactionService
	userService  service.UserService
}

// NewQuickTransactionHandler creates QuickTransactionHandler.
func NewQuickTransactionHandler(quickService service.QuickTransactionService, userService service.UserService) *QuickTransactionHandler {
	return &QuickTransactionHandler{quickService: quickService, userService: userService}
}

func (h *QuickTransactionHandler) resolveUserID(c *gin.Context) (uuid.UUID, bool) {
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

// List handles GET /api/v1/quick-transactions.
func (h *QuickTransactionHandler) List(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	items, total, err := h.quickService.List(c.Request.Context(), userID)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to list quick transactions")
		return
	}
	responses := make([]dto.QuickTransactionResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, dto.NewQuickTransactionResponse(item))
	}
	dto.Paginated(c, responses, 1, len(responses), total)
}

// Create handles POST /api/v1/quick-transactions.
func (h *QuickTransactionHandler) Create(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var req dto.CreateQuickTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	item, err := h.quickService.Create(c.Request.Context(), userID, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to create quick transaction")
		return
	}
	dto.Created(c, dto.NewQuickTransactionResponse(*item))
}

// Update handles PATCH /api/v1/quick-transactions/:id.
func (h *QuickTransactionHandler) Update(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid quick transaction id")
		return
	}
	var req dto.UpdateQuickTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	item, err := h.quickService.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to update quick transaction")
		return
	}
	dto.Success(c, dto.NewQuickTransactionResponse(*item))
}

// Delete handles DELETE /api/v1/quick-transactions/:id.
func (h *QuickTransactionHandler) Delete(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid quick transaction id")
		return
	}
	if err := h.quickService.Delete(c.Request.Context(), userID, id); err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to delete quick transaction")
		return
	}
	c.Status(204)
}

// Execute handles POST /api/v1/quick-transactions/:id/execute.
func (h *QuickTransactionHandler) Execute(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid quick transaction id")
		return
	}
	txDate := time.Now().Format("2006-01-02")
	item, err := h.quickService.Execute(c.Request.Context(), userID, id, txDate)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to execute quick transaction")
		return
	}
	dto.Created(c, dto.NewTransactionResponse(*item))
}

// Reorder handles PATCH /api/v1/quick-transactions/reorder.
func (h *QuickTransactionHandler) Reorder(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var req dto.ReorderQuickTransactionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	if err := h.quickService.Reorder(c.Request.Context(), userID, req.IDs); err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to reorder quick transactions")
		return
	}
	c.Status(204)
}
