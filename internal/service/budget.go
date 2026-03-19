package service

import (
	"context"
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"ledgerA/internal/repository"
	"math"
	"time"

	"github.com/google/uuid"
)

type budgetService struct {
	repo    repository.BudgetRepository
	txnRepo repository.TransactionRepository
}

// NewBudgetService creates a new BudgetService.
func NewBudgetService(repo repository.BudgetRepository, txnRepo repository.TransactionRepository) BudgetService {
	return &budgetService{repo: repo, txnRepo: txnRepo}
}

func (s *budgetService) List(ctx context.Context, userID uuid.UUID) ([]model.Budget, int64, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *budgetService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateBudgetRequest) (*model.Budget, error) {
	item := &model.Budget{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Period:     req.Period,
		IsActive:   true,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("budgetService.Create: %w", err)
	}
	return item, nil
}

func (s *budgetService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateBudgetRequest) (*model.Budget, error) {
	item, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	if req.Amount != nil {
		item.Amount = *req.Amount
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *budgetService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return s.repo.Delete(ctx, userID, id)
}

func (s *budgetService) Progress(ctx context.Context, userID uuid.UUID) ([]dto.BudgetProgressResponse, error) {
	budgets, _, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var results []dto.BudgetProgressResponse

	for _, b := range budgets {
		if !b.IsActive {
			continue
		}

		var dateFrom, dateTo time.Time
		switch b.Period {
		case "monthly":
			dateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
			dateTo = dateFrom.AddDate(0, 1, -1)
		case "yearly":
			dateFrom = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
			dateTo = time.Date(now.Year(), 12, 31, 0, 0, 0, 0, time.UTC)
		}

		filter := repository.TransactionListFilter{
			CategoryID: &b.CategoryID,
			DateFrom:   &dateFrom,
			DateTo:     &dateTo,
			Page:       1,
			PerPage:    10000,
		}

		txns, _, err := s.txnRepo.ListByUserID(ctx, userID, filter)
		if err != nil {
			continue
		}

		var spent float64
		for _, tx := range txns {
			if tx.Amount < 0 {
				spent += math.Abs(tx.Amount)
			}
		}

		remaining := b.Amount - spent
		if remaining < 0 {
			remaining = 0
		}
		percent := 0.0
		if b.Amount > 0 {
			percent = (spent / b.Amount) * 100
		}

		results = append(results, dto.BudgetProgressResponse{
			BudgetResponse: dto.NewBudgetResponse(b),
			Spent:          spent,
			Remaining:      remaining,
			Percent:        percent,
		})
	}

	return results, nil
}
