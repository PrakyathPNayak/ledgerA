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
		setup   func(*mocks.TransactionRepository, *mocks.CategoryRepository, *mocks.SubcategoryRepository)
		wantErr bool
	}{
		{
			name:  "ok totals",
			query: dto.StatsQuery{Period: "month", Value: "2026-03", AccountID: &idStr},
			setup: func(repo *mocks.TransactionRepository, catRepo *mocks.CategoryRepository, subRepo *mocks.SubcategoryRepository) {
				categoryID := uuid.New()
				subcategoryID := uuid.New()
				category := model.Category{BaseModel: model.BaseModel{ID: categoryID}, Name: "Food"}
				subcategory := model.Subcategory{BaseModel: model.BaseModel{ID: subcategoryID}, CategoryID: categoryID, Name: "Groceries"}

				repo.On("ListByUserID", ctx, uid, repository.TransactionListFilter{AccountID: &accountID, Type: "all", SortBy: "transaction_date", SortDir: "asc", Page: 1, PerPage: 5000}).Return([]model.Transaction{{Amount: 10, TransactionDate: now}, {Amount: -4, TransactionDate: now}}, int64(2), nil)
				catRepo.On("ListByUserID", ctx, uid).Return([]model.Category{category}, int64(1), nil)
				subRepo.On("ListByCategoryID", ctx, uid, categoryID).Return([]model.Subcategory{subcategory}, int64(1), nil)
			},
		},
		{name: "bad account id", query: dto.StatsQuery{Period: "month", Value: "2026-03", AccountID: ptr("bad")}, setup: func(_ *mocks.TransactionRepository, _ *mocks.CategoryRepository, _ *mocks.SubcategoryRepository) {}, wantErr: true},
		{name: "repo error", query: dto.StatsQuery{Period: "month", Value: "2026-03"}, setup: func(repo *mocks.TransactionRepository, _ *mocks.CategoryRepository, _ *mocks.SubcategoryRepository) {
			repo.On("ListByUserID", ctx, uid, repository.TransactionListFilter{Type: "all", SortBy: "transaction_date", SortDir: "asc", Page: 1, PerPage: 5000}).Return([]model.Transaction{}, int64(0), errors.New("db"))
		}, wantErr: true},
		{name: "ok empty", query: dto.StatsQuery{Period: "month", Value: "2026-03"}, setup: func(repo *mocks.TransactionRepository, catRepo *mocks.CategoryRepository, _ *mocks.SubcategoryRepository) {
			repo.On("ListByUserID", ctx, uid, repository.TransactionListFilter{Type: "all", SortBy: "transaction_date", SortDir: "asc", Page: 1, PerPage: 5000}).Return([]model.Transaction{}, int64(0), nil)
			catRepo.On("ListByUserID", ctx, uid).Return([]model.Category{}, int64(0), nil)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.TransactionRepository{}
			catRepo := &mocks.CategoryRepository{}
			subRepo := &mocks.SubcategoryRepository{}
			tt.setup(repo, catRepo, subRepo)
			svc := NewStatsService(repo, catRepo, subRepo)
			result, err := svc.Summary(ctx, uid, tt.query)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotNil(t, result.Timeseries)
				require.NotNil(t, result.CategoryBreakdownExpense)
				require.NotNil(t, result.CategoryBreakdownIncome)
			}
		})
	}
}
