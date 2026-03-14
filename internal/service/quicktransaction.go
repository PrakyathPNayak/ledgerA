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

type quickTransactionService struct {
	quickRepo repository.QuickTransactionRepository
	txService TransactionService
	auditRepo repository.AuditRepository
}

// NewQuickTransactionService creates a new QuickTransactionService.
func NewQuickTransactionService(quickRepo repository.QuickTransactionRepository, txService TransactionService, auditRepo repository.AuditRepository) QuickTransactionService {
	return &quickTransactionService{quickRepo: quickRepo, txService: txService, auditRepo: auditRepo}
}

func (s *quickTransactionService) List(ctx context.Context, userID uuid.UUID) ([]model.QuickTransaction, int64, error) {
	items, total, err := s.quickRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("quickTransactionService.List: %w", err)
	}
	return items, total, nil
}

func (s *quickTransactionService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateQuickTransactionRequest) (*model.QuickTransaction, error) {
	item := req.ToModel(userID)
	if err := s.quickRepo.Create(ctx, &item); err != nil {
		return nil, fmt.Errorf("quickTransactionService.Create: %w", err)
	}
	if err := s.writeAudit(ctx, userID, item.ID, "create", nil, item); err != nil {
		return nil, fmt.Errorf("quickTransactionService.Create.Audit: %w", err)
	}
	return &item, nil
}

func (s *quickTransactionService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateQuickTransactionRequest) (*model.QuickTransaction, error) {
	item, err := s.quickRepo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("quickTransactionService.Update.Find: %w", err)
	}
	before := *item
	if req.Label != nil {
		item.Label = *req.Label
	}
	item.AccountID = req.AccountID
	item.CategoryID = req.CategoryID
	item.SubcategoryID = req.SubcategoryID
	item.Amount = req.Amount
	item.Notes = req.Notes
	if err := s.quickRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("quickTransactionService.Update.Save: %w", err)
	}
	if err := s.writeAudit(ctx, userID, item.ID, "update", before, item); err != nil {
		return nil, fmt.Errorf("quickTransactionService.Update.Audit: %w", err)
	}
	return item, nil
}

func (s *quickTransactionService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	item, err := s.quickRepo.FindByID(ctx, userID, id)
	if err != nil {
		return fmt.Errorf("quickTransactionService.Delete.Find: %w", err)
	}
	if err := s.quickRepo.Delete(ctx, userID, id); err != nil {
		return fmt.Errorf("quickTransactionService.Delete.Delete: %w", err)
	}
	if err := s.writeAudit(ctx, userID, id, "delete", item, nil); err != nil {
		return fmt.Errorf("quickTransactionService.Delete.Audit: %w", err)
	}
	return nil
}

func (s *quickTransactionService) Reorder(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	if err := s.quickRepo.Reorder(ctx, userID, ids); err != nil {
		return fmt.Errorf("quickTransactionService.Reorder: %w", err)
	}
	return nil
}

func (s *quickTransactionService) Execute(ctx context.Context, userID uuid.UUID, id uuid.UUID, txDate string) (*model.Transaction, error) {
	item, err := s.quickRepo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("quickTransactionService.Execute.Find: %w", err)
	}
	if item.AccountID == nil || item.CategoryID == nil || item.Amount == nil {
		return nil, fmt.Errorf("quickTransactionService.Execute: quick transaction missing required account/category/amount")
	}

	req := dto.CreateTransactionRequest{
		AccountID:       *item.AccountID,
		CategoryID:      *item.CategoryID,
		SubcategoryID:   item.SubcategoryID,
		Name:            item.Label,
		Amount:          *item.Amount,
		TransactionDate: txDate,
		Notes:           item.Notes,
		IsScheduled:     false,
	}
	return s.txService.Create(ctx, userID, req)
}

func (s *quickTransactionService) writeAudit(ctx context.Context, userID uuid.UUID, entityID uuid.UUID, action string, before any, after any) error {
	diffPayload := map[string]any{"before": before, "after": after}
	diffBytes, err := json.Marshal(diffPayload)
	if err != nil {
		return fmt.Errorf("quickTransactionService.writeAudit.Marshal: %w", err)
	}
	entry := model.AuditLog{UserID: userID, EntityType: "quick_transaction", EntityID: entityID, Action: action, Diff: string(diffBytes)}
	if err := s.auditRepo.Create(ctx, &entry); err != nil {
		return fmt.Errorf("quickTransactionService.writeAudit.Create: %w", err)
	}
	return nil
}
