package service

import (
	"context"
	"ledgerA/internal/model"
	"ledgerA/internal/repository"

	"github.com/google/uuid"
)

type auditService struct {
	repo repository.AuditRepository
}

// NewAuditService returns a new AuditService backed by repo.
func NewAuditService(repo repository.AuditRepository) AuditService {
	return &auditService{repo: repo}
}

// List returns a paginated audit log for the given user.
func (s *auditService) List(ctx context.Context, userID uuid.UUID, filter repository.AuditListFilter) ([]model.AuditLog, int64, error) {
	return s.repo.ListByUserID(ctx, userID, filter)
}
