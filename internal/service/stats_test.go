package service

import (
	"context"
	"errors"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"ledgerA/internal/repository"
	"ledgerA/internal/service/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStatsServiceSummary(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New()
	now := time.Now()
	accountID := uuid.New()
	idStr := accountID.String()

	tests := []struct {
		name    string
		query   dto.StatsQuery
		setup   func(*mocks.TransactionRepository)
		wantErr bool
	}{
		{
			name:  "ok totals",
			query: dto.StatsQuery{Period: "month", Value: "2026-03", AccountID: &idStr},
			setup: func(repo *mocks.TransactionRepository) {
				repo.On("ListByUserID", ctx, uid, repository.TransactionListFilter{AccountID: &accountID, Type: "all", SortBy: "transaction_date", SortDir: "asc", Page: 1, PerPage: 5000}).Return([]model.Transaction{{Amount: 10, TransactionDate: now}, {Amount: -4, TransactionDate: now}}, int64(2), nil)
			},
		},
		{name: "bad account id", query: dto.StatsQuery{Period: "month", Value: "2026-03", AccountID: ptr("bad")}, setup: func(_ *mocks.TransactionRepository) {}, wantErr: true},
		{name: "repo error", query: dto.StatsQuery{Period: "month", Value: "2026-03"}, setup: func(repo *mocks.TransactionRepository) {
			repo.On("ListByUserID", ctx, uid, repository.TransactionListFilter{Type: "all", SortBy: "transaction_date", SortDir: "asc", Page: 1, PerPage: 5000}).Return([]model.Transaction{}, int64(0), errors.New("db"))
		}, wantErr: true},
		{name: "ok empty", query: dto.StatsQuery{Period: "month", Value: "2026-03"}, setup: func(repo *mocks.TransactionRepository) {
			repo.On("ListByUserID", ctx, uid, repository.TransactionListFilter{Type: "all", SortBy: "transaction_date", SortDir: "asc", Page: 1, PerPage: 5000}).Return([]model.Transaction{}, int64(0), nil)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.TransactionRepository{}
			tt.setup(repo)
			svc := NewStatsService(repo)
			_, err := svc.Summary(ctx, uid, tt.query)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
