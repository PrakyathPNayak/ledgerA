package mocks

import (
	"context"
	"ledgerA/internal/model"
	"ledgerA/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct{ mock.Mock }

func (m *UserRepository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*model.User, error) {
	args := m.Called(ctx, firebaseUID)
	if u := args.Get(0); u != nil {
		return u.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if u := args.Get(0); u != nil {
		return u.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *UserRepository) Create(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *UserRepository) Update(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}

type AccountRepository struct{ mock.Mock }

func (m *AccountRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Account, int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Account), args.Get(1).(int64), args.Error(2)
}
func (m *AccountRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Account, error) {
	args := m.Called(ctx, userID, id)
	if a := args.Get(0); a != nil {
		return a.(*model.Account), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *AccountRepository) Create(ctx context.Context, account *model.Account) error {
	return m.Called(ctx, account).Error(0)
}
func (m *AccountRepository) Update(ctx context.Context, account *model.Account) error {
	return m.Called(ctx, account).Error(0)
}
func (m *AccountRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return m.Called(ctx, userID, id).Error(0)
}
func (m *AccountRepository) UpdateBalance(ctx context.Context, userID uuid.UUID, id uuid.UUID, delta float64) error {
	return m.Called(ctx, userID, id, delta).Error(0)
}

type CategoryRepository struct{ mock.Mock }

func (m *CategoryRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Category, int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Category), args.Get(1).(int64), args.Error(2)
}
func (m *CategoryRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Category, error) {
	args := m.Called(ctx, userID, id)
	if c := args.Get(0); c != nil {
		return c.(*model.Category), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *CategoryRepository) FindByName(ctx context.Context, userID uuid.UUID, name string) (*model.Category, error) {
	args := m.Called(ctx, userID, name)
	if c := args.Get(0); c != nil {
		return c.(*model.Category), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *CategoryRepository) Create(ctx context.Context, category *model.Category) error {
	return m.Called(ctx, category).Error(0)
}
func (m *CategoryRepository) Update(ctx context.Context, category *model.Category) error {
	return m.Called(ctx, category).Error(0)
}
func (m *CategoryRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return m.Called(ctx, userID, id).Error(0)
}

type SubcategoryRepository struct{ mock.Mock }

func (m *SubcategoryRepository) ListByCategoryID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]model.Subcategory, int64, error) {
	args := m.Called(ctx, userID, categoryID)
	return args.Get(0).([]model.Subcategory), args.Get(1).(int64), args.Error(2)
}
func (m *SubcategoryRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Subcategory, error) {
	args := m.Called(ctx, userID, id)
	if s := args.Get(0); s != nil {
		return s.(*model.Subcategory), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *SubcategoryRepository) FindByNameForCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, name string) (*model.Subcategory, error) {
	args := m.Called(ctx, userID, categoryID, name)
	if s := args.Get(0); s != nil {
		return s.(*model.Subcategory), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *SubcategoryRepository) Create(ctx context.Context, subcategory *model.Subcategory) error {
	return m.Called(ctx, subcategory).Error(0)
}
func (m *SubcategoryRepository) Update(ctx context.Context, subcategory *model.Subcategory) error {
	return m.Called(ctx, subcategory).Error(0)
}
func (m *SubcategoryRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return m.Called(ctx, userID, id).Error(0)
}

type TransactionRepository struct{ mock.Mock }

func (m *TransactionRepository) ListByUserID(ctx context.Context, userID uuid.UUID, filter repository.TransactionListFilter) ([]model.Transaction, int64, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).([]model.Transaction), args.Get(1).(int64), args.Error(2)
}
func (m *TransactionRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Transaction, error) {
	args := m.Called(ctx, userID, id)
	if t := args.Get(0); t != nil {
		return t.(*model.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *TransactionRepository) Create(ctx context.Context, tx *model.Transaction) error {
	return m.Called(ctx, tx).Error(0)
}
func (m *TransactionRepository) Update(ctx context.Context, tx *model.Transaction) error {
	return m.Called(ctx, tx).Error(0)
}
func (m *TransactionRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return m.Called(ctx, userID, id).Error(0)
}

type QuickTransactionRepository struct{ mock.Mock }

func (m *QuickTransactionRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.QuickTransaction, int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.QuickTransaction), args.Get(1).(int64), args.Error(2)
}
func (m *QuickTransactionRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.QuickTransaction, error) {
	args := m.Called(ctx, userID, id)
	if q := args.Get(0); q != nil {
		return q.(*model.QuickTransaction), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *QuickTransactionRepository) Create(ctx context.Context, quick *model.QuickTransaction) error {
	return m.Called(ctx, quick).Error(0)
}
func (m *QuickTransactionRepository) Update(ctx context.Context, quick *model.QuickTransaction) error {
	return m.Called(ctx, quick).Error(0)
}
func (m *QuickTransactionRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return m.Called(ctx, userID, id).Error(0)
}
func (m *QuickTransactionRepository) Reorder(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	return m.Called(ctx, userID, ids).Error(0)
}

type AuditRepository struct{ mock.Mock }

func (m *AuditRepository) Create(ctx context.Context, entry *model.AuditLog) error {
	return m.Called(ctx, entry).Error(0)
}
func (m *AuditRepository) ListByUserID(ctx context.Context, userID uuid.UUID, filter repository.AuditListFilter) ([]model.AuditLog, int64, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).([]model.AuditLog), args.Get(1).(int64), args.Error(2)
}

func MustDate(value string) time.Time {
	t, _ := time.Parse("2006-01-02", value)
	return t
}
