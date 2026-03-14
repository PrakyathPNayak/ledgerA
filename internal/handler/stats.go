package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StatsHandler handles stats endpoints.
type StatsHandler struct {
	statsService service.StatsService
	userService  service.UserService
}

// NewStatsHandler creates StatsHandler.
func NewStatsHandler(statsService service.StatsService, userService service.UserService) *StatsHandler {
	return &StatsHandler{statsService: statsService, userService: userService}
}

func (h *StatsHandler) resolveUserID(c *gin.Context) (uuid.UUID, bool) {
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

// Summary handles GET /api/v1/stats/summary.
func (h *StatsHandler) Summary(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var query dto.StatsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", err.Error())
		return
	}
	result, err := h.statsService.Summary(c.Request.Context(), userID, query)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", err.Error())
		return
	}
	dto.Success(c, result)
}

// ExportPDF handles GET /api/v1/stats/export/pdf.
func (h *StatsHandler) ExportPDF(c *gin.Context) {
	dto.Error(c, 501, "ERR_NOT_IMPLEMENTED", "PDF export service is wired in TODO-014")
}

// Compare handles GET /api/v1/stats/compare.
func (h *StatsHandler) Compare(c *gin.Context) {
	dto.Error(c, 501, "ERR_NOT_IMPLEMENTED", "compare stats is not implemented yet")
}
