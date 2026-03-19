package handler

import (
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/service"
	"net/http"

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

// Summary handles GET /api/v1/stats/summary.
func (h *StatsHandler) Summary(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var query dto.StatsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid query parameters")
		return
	}
	result, err := h.statsService.Summary(c.Request.Context(), userID, query)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to generate stats summary")
		return
	}
	dto.Success(c, result)
}

// ExportPDF handles GET /api/v1/stats/export/pdf.
func (h *StatsHandler) ExportPDF(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var query dto.StatsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid query parameters")
		return
	}

	pdfBytes, err := h.statsService.ExportPDF(c.Request.Context(), userID, query)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to generate PDF report")
		return
	}

	filename := fmt.Sprintf("ledgerA-stats-%s-%s.pdf", query.Period, query.Value)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// Compare handles GET /api/v1/stats/compare.
func (h *StatsHandler) Compare(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}

	var query dto.CompareQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid query parameters")
		return
	}

	result, err := h.statsService.Compare(c.Request.Context(), userID, query)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to generate comparison")
		return
	}
	dto.Success(c, result)
}

// Monthly handles GET /api/v1/stats/monthly.
func (h *StatsHandler) Monthly(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}

	var query dto.MonthlyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid query parameters")
		return
	}

	result, err := h.statsService.Monthly(c.Request.Context(), userID, query)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to generate monthly report")
		return
	}
	dto.Success(c, result)
}
