package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/repository"
	"ledgerA/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuditHandler handles audit log endpoints.
type AuditHandler struct {
	auditService service.AuditService
	userService  service.UserService
}

// NewAuditHandler creates AuditHandler.
func NewAuditHandler(auditService service.AuditService, userService service.UserService) *AuditHandler {
	return &AuditHandler{auditService: auditService, userService: userService}
}

func (h *AuditHandler) resolveUserID(c *gin.Context) (uuid.UUID, bool) {
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

// List handles GET /api/v1/audit.
func (h *AuditHandler) List(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}

	filter := repository.AuditListFilter{
		Page:    1,
		PerPage: 20,
	}

	if v := c.Query("entity_type"); v != "" {
		filter.EntityType = &v
	}

	if v := c.Query("entity_id"); v != "" {
		parsed, err := uuid.Parse(v)
		if err != nil {
			dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid entity_id")
			return
		}
		filter.EntityID = &parsed
	}

	if v := c.Query("page"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil || p < 1 {
			dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid page")
			return
		}
		filter.Page = p
	}

	if v := c.Query("per_page"); v != "" {
		pp, err := strconv.Atoi(v)
		if err != nil || pp < 1 || pp > 100 {
			dto.Error(c, 400, "ERR_BAD_REQUEST", "per_page must be between 1 and 100")
			return
		}
		filter.PerPage = pp
	}

	logs, total, err := h.auditService.List(c.Request.Context(), userID, filter)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}

	dto.Paginated(c, logs, filter.Page, filter.PerPage, total)
}
