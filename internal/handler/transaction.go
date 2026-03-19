package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TransactionHandler handles transaction endpoints.
type TransactionHandler struct {
	transactionService service.TransactionService
	userService        service.UserService
}

// NewTransactionHandler creates TransactionHandler.
func NewTransactionHandler(transactionService service.TransactionService, userService service.UserService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService, userService: userService}
}

func (h *TransactionHandler) resolveUserID(c *gin.Context) (uuid.UUID, bool) {
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

// List handles GET /api/v1/transactions.
func (h *TransactionHandler) List(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var filters dto.TransactionFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", err.Error())
		return
	}
	items, total, err := h.transactionService.List(c.Request.Context(), userID, filters)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to list transactions")
		return
	}
	responses := make([]dto.TransactionResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, dto.NewTransactionResponse(item))
	}
	page := filters.Page
	if page <= 0 {
		page = 1
	}
	perPage := filters.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	dto.Paginated(c, responses, page, perPage, total)
}

// Create handles POST /api/v1/transactions.
func (h *TransactionHandler) Create(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	item, err := h.transactionService.Create(c.Request.Context(), userID, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to create transaction")
		return
	}
	dto.Created(c, dto.NewTransactionResponse(*item))
}

// Get handles GET /api/v1/transactions/:id.
func (h *TransactionHandler) Get(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid transaction id")
		return
	}
	item, err := h.transactionService.Get(c.Request.Context(), userID, id)
	if err != nil {
		dto.Error(c, 404, "ERR_NOT_FOUND", "transaction not found")
		return
	}
	dto.Success(c, dto.NewTransactionResponse(*item))
}

// Update handles PATCH /api/v1/transactions/:id.
func (h *TransactionHandler) Update(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid transaction id")
		return
	}
	var req dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	item, err := h.transactionService.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to update transaction")
		return
	}
	dto.Success(c, dto.NewTransactionResponse(*item))
}

// Delete handles DELETE /api/v1/transactions/:id.
func (h *TransactionHandler) Delete(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid transaction id")
		return
	}
	if err := h.transactionService.Delete(c.Request.Context(), userID, id); err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to delete transaction")
		return
	}
	c.Status(204)
}
