package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"ledgerA/internal/repository"

	"github.com/google/uuid"
)

type accountService struct {
	accountRepo repository.AccountRepository
	auditRepo   repository.AuditRepository
}

// NewAccountService creates a new AccountService.
func NewAccountService(accountRepo repository.AccountRepository, auditRepo repository.AuditRepository) AccountService {
	return &accountService{accountRepo: accountRepo, auditRepo: auditRepo}
}

func (s *accountService) List(ctx context.Context, userID uuid.UUID) ([]model.Account, int64, error) {
	accounts, total, err := s.accountRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("accountService.List: %w", err)
	}
	return accounts, total, nil
}

func (s *accountService) Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Account, error) {
	account, err := s.accountRepo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("accountService.Get: %w", err)
	}
	return account, nil
}

func (s *accountService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateAccountRequest) (*model.Account, error) {
	account := req.ToModel(userID)
	if err := s.accountRepo.Create(ctx, &account); err != nil {
		return nil, fmt.Errorf("accountService.Create: %w", err)
	}
	if err := s.writeAudit(ctx, userID, "account", account.ID, "create", nil, account); err != nil {
		return nil, fmt.Errorf("accountService.Create.Audit: %w", err)
	}
	return &account, nil
}

func (s *accountService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateAccountRequest) (*model.Account, error) {
	account, err := s.accountRepo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("accountService.Update.Find: %w", err)
	}
	before := *account
	if req.Name != "" {
		account.Name = req.Name
	}
	if req.IsArchived != nil {
		account.IsArchived = *req.IsArchived
	}
	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("accountService.Update.Save: %w", err)
	}
	if err := s.writeAudit(ctx, userID, "account", account.ID, "update", before, account); err != nil {
		return nil, fmt.Errorf("accountService.Update.Audit: %w", err)
	}
	return account, nil
}

func (s *accountService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	account, err := s.accountRepo.FindByID(ctx, userID, id)
	if err != nil {
		return fmt.Errorf("accountService.Delete.Find: %w", err)
	}
	if err := s.accountRepo.Delete(ctx, userID, id); err != nil {
		return fmt.Errorf("accountService.Delete.Delete: %w", err)
	}
	if err := s.writeAudit(ctx, userID, "account", id, "delete", account, nil); err != nil {
		return fmt.Errorf("accountService.Delete.Audit: %w", err)
	}
	return nil
}

func (s *accountService) writeAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string, before any, after any) error {
	diffPayload := map[string]any{"before": before, "after": after}
	diffBytes, err := json.Marshal(diffPayload)
	if err != nil {
		return fmt.Errorf("accountService.writeAudit.Marshal: %w", err)
	}

	entry := model.AuditLog{
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Diff:       string(diffBytes),
	}
	if err := s.auditRepo.Create(ctx, &entry); err != nil {
		return fmt.Errorf("accountService.writeAudit.Create: %w", err)
	}
	return nil
}
