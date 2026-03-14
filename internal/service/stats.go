package service

import (
	"context"
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/repository"

	"github.com/google/uuid"
)

type statsService struct {
	txRepo repository.TransactionRepository
}

// NewStatsService creates a new StatsService.
func NewStatsService(txRepo repository.TransactionRepository) StatsService {
	return &statsService{txRepo: txRepo}
}

func (s *statsService) Summary(ctx context.Context, userID uuid.UUID, query dto.StatsQuery) (*dto.StatsSummaryResponse, error) {
	filter := repository.TransactionListFilter{Type: "all", SortBy: "transaction_date", SortDir: "asc", Page: 1, PerPage: 5000}
	if query.AccountID != nil && *query.AccountID != "" && *query.AccountID != "all" {
		parsedID, err := uuid.Parse(*query.AccountID)
		if err != nil {
			return nil, fmt.Errorf("statsService.Summary.ParseAccountID: %w", err)
		}
		filter.AccountID = &parsedID
	}

	items, _, err := s.txRepo.ListByUserID(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("statsService.Summary.List: %w", err)
	}

	resp := &dto.StatsSummaryResponse{
		CategoryBreakdownExpense: []dto.CategoryBreakdownItem{},
		CategoryBreakdownIncome:  []dto.CategoryBreakdownItem{},
		Timeseries:               []dto.TimeseriesPoint{},
	}

	for _, item := range items {
		if item.Amount >= 0 {
			resp.TotalIncome += item.Amount
		} else {
			resp.TotalExpense += -item.Amount
		}
	}
	resp.Net = resp.TotalIncome - resp.TotalExpense

	return resp, nil
}
