package service

import (
	"context"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"ledgerA/internal/repository"

	"github.com/google/uuid"
)

// UserService defines user business logic.
type UserService interface {
	Sync(ctx context.Context, firebaseUID string, req dto.SyncRequest) (*model.User, error)
	GetMe(ctx context.Context, firebaseUID string) (*model.User, error)
	UpdateMe(ctx context.Context, firebaseUID string, req dto.UpdateUserRequest) (*model.User, error)
}

// AccountService defines account business logic.
type AccountService interface {
	List(ctx context.Context, userID uuid.UUID) ([]model.Account, int64, error)
	Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Account, error)
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateAccountRequest) (*model.Account, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateAccountRequest) (*model.Account, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

// CategoryService defines category business logic.
type CategoryService interface {
	List(ctx context.Context, userID uuid.UUID) ([]model.Category, int64, error)
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateCategoryRequest) (*model.Category, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateCategoryRequest) (*model.Category, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

// SubcategoryService defines subcategory business logic.
type SubcategoryService interface {
	ListByCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]model.Subcategory, int64, error)
	Create(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, req dto.CreateSubcategoryRequest) (*model.Subcategory, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateSubcategoryRequest) (*model.Subcategory, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	GetOrCreateFromCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, categoryName string) (*model.Subcategory, error)
}

// TransactionService defines transaction business logic.
type TransactionService interface {
	List(ctx context.Context, userID uuid.UUID, filters dto.TransactionFilters) ([]model.Transaction, int64, error)
	Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Transaction, error)
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateTransactionRequest) (*model.Transaction, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateTransactionRequest) (*model.Transaction, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	Transfer(ctx context.Context, userID uuid.UUID, req dto.TransferRequest) (*model.Transaction, *model.Transaction, error)
}

// QuickTransactionService defines quick transaction business logic.
type QuickTransactionService interface {
	List(ctx context.Context, userID uuid.UUID) ([]model.QuickTransaction, int64, error)
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateQuickTransactionRequest) (*model.QuickTransaction, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateQuickTransactionRequest) (*model.QuickTransaction, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	Reorder(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	Execute(ctx context.Context, userID uuid.UUID, id uuid.UUID, txDate string) (*model.Transaction, error)
}

// StatsService defines stats business logic.
type StatsService interface {
	Summary(ctx context.Context, userID uuid.UUID, query dto.StatsQuery) (*dto.StatsSummaryResponse, error)
	ExportPDF(ctx context.Context, userID uuid.UUID, query dto.StatsQuery) ([]byte, error)
	Compare(ctx context.Context, userID uuid.UUID, query dto.CompareQuery) (*dto.CompareResponse, error)
	Monthly(ctx context.Context, userID uuid.UUID, query dto.MonthlyQuery) (*dto.MonthlyReportResponse, error)
}

// AuditService defines audit query business logic.
type AuditService interface {
	List(ctx context.Context, userID uuid.UUID, filter repository.AuditListFilter) ([]model.AuditLog, int64, error)
}

// RecurringService defines recurring transaction business logic.
type RecurringService interface {
	List(ctx context.Context, userID uuid.UUID) ([]model.RecurringTransaction, int64, error)
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateRecurringRequest) (*model.RecurringTransaction, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateRecurringRequest) (*model.RecurringTransaction, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	ProcessDue(ctx context.Context) (int, error)
}

// BudgetService defines budget business logic.
type BudgetService interface {
	List(ctx context.Context, userID uuid.UUID) ([]model.Budget, int64, error)
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateBudgetRequest) (*model.Budget, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateBudgetRequest) (*model.Budget, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	Progress(ctx context.Context, userID uuid.UUID) ([]dto.BudgetProgressResponse, error)
}
