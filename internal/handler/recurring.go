package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RecurringHandler handles recurring transaction endpoints.
type RecurringHandler struct {
	recurringService service.RecurringService
	userService      service.UserService
}

// NewRecurringHandler creates RecurringHandler.
func NewRecurringHandler(recurringService service.RecurringService, userService service.UserService) *RecurringHandler {
	return &RecurringHandler{recurringService: recurringService, userService: userService}
}

func (h *RecurringHandler) resolveUserID(c *gin.Context) (uuid.UUID, bool) {
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

// List handles GET /api/v1/recurring.
func (h *RecurringHandler) List(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	items, total, err := h.recurringService.List(c.Request.Context(), userID)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to list recurring transactions")
		return
	}
	responses := make([]dto.RecurringResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, dto.NewRecurringResponse(item))
	}
	dto.Paginated(c, responses, 1, len(responses), total)
}

// Create handles POST /api/v1/recurring.
func (h *RecurringHandler) Create(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var req dto.CreateRecurringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 400, "ERR_VALIDATION", err.Error())
		return
	}
	item, err := h.recurringService.Create(c.Request.Context(), userID, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Created(c, dto.NewRecurringResponse(*item))
}

// Update handles PATCH /api/v1/recurring/:id.
func (h *RecurringHandler) Update(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_VALIDATION", "invalid id")
		return
	}
	var req dto.UpdateRecurringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 400, "ERR_VALIDATION", err.Error())
		return
	}
	item, err := h.recurringService.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Success(c, dto.NewRecurringResponse(*item))
}

// Delete handles DELETE /api/v1/recurring/:id.
func (h *RecurringHandler) Delete(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_VALIDATION", "invalid id")
		return
	}
	if err := h.recurringService.Delete(c.Request.Context(), userID, id); err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Success(c, gin.H{"deleted": true})
}

// ProcessDue handles POST /api/v1/recurring/process — triggers processing of due recurring transactions.
func (h *RecurringHandler) ProcessDue(c *gin.Context) {
	count, err := h.recurringService.ProcessDue(c.Request.Context())
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Success(c, gin.H{"processed": count})
}
