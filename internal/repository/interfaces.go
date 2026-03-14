package repository

import (
	"context"
	"ledgerA/internal/model"
	"time"

	"github.com/google/uuid"
)

// TransactionListFilter defines filters for transaction lists.
type TransactionListFilter struct {
	AccountID     *uuid.UUID
	CategoryID    *uuid.UUID
	SubcategoryID *uuid.UUID
	DateFrom      *time.Time
	DateTo        *time.Time
	Search        *string
	Type          string
	SortBy        string
	SortDir       string
	Page          int
	PerPage       int
	PassbookMode  bool
}

// AuditListFilter defines filters for audit list.
type AuditListFilter struct {
	EntityType *string
	EntityID   *uuid.UUID
	Page       int
	PerPage    int
}

// UserRepository defines user persistence methods.
type UserRepository interface {
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
}

// AccountRepository defines account persistence methods.
type AccountRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Account, int64, error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Account, error)
	Create(ctx context.Context, account *model.Account) error
	Update(ctx context.Context, account *model.Account) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	UpdateBalance(ctx context.Context, userID uuid.UUID, id uuid.UUID, delta float64) error
}

// CategoryRepository defines category persistence methods.
type CategoryRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Category, int64, error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Category, error)
	FindByName(ctx context.Context, userID uuid.UUID, name string) (*model.Category, error)
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

// SubcategoryRepository defines subcategory persistence methods.
type SubcategoryRepository interface {
	ListByCategoryID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]model.Subcategory, int64, error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Subcategory, error)
	FindByNameForCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, name string) (*model.Subcategory, error)
	Create(ctx context.Context, subcategory *model.Subcategory) error
	Update(ctx context.Context, subcategory *model.Subcategory) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

// TransactionRepository defines transaction persistence methods.
type TransactionRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID, filter TransactionListFilter) ([]model.Transaction, int64, error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Transaction, error)
	Create(ctx context.Context, tx *model.Transaction) error
	Update(ctx context.Context, tx *model.Transaction) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

// QuickTransactionRepository defines quick transaction persistence methods.
type QuickTransactionRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.QuickTransaction, int64, error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.QuickTransaction, error)
	Create(ctx context.Context, quick *model.QuickTransaction) error
	Update(ctx context.Context, quick *model.QuickTransaction) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	Reorder(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
}

// AuditRepository defines audit persistence methods.
type AuditRepository interface {
	Create(ctx context.Context, entry *model.AuditLog) error
	ListByUserID(ctx context.Context, userID uuid.UUID, filter AuditListFilter) ([]model.AuditLog, int64, error)
}
