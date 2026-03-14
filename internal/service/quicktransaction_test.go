package service

import (
	"context"
	"errors"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"ledgerA/internal/service/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type txServiceStub struct{}

func (t txServiceStub) List(context.Context, uuid.UUID, dto.TransactionFilters) ([]model.Transaction, int64, error) {
	return nil, 0, nil
}
func (t txServiceStub) Get(context.Context, uuid.UUID, uuid.UUID) (*model.Transaction, error) {
	return nil, nil
}
func (t txServiceStub) Create(context.Context, uuid.UUID, dto.CreateTransactionRequest) (*model.Transaction, error) {
	return &model.Transaction{}, nil
}
func (t txServiceStub) Update(context.Context, uuid.UUID, uuid.UUID, dto.UpdateTransactionRequest) (*model.Transaction, error) {
	return nil, nil
}
func (t txServiceStub) Delete(context.Context, uuid.UUID, uuid.UUID) error { return nil }

func TestQuickTransactionServiceCreate(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New()
	tests := []struct {
		name    string
		setup   func(*mocks.QuickTransactionRepository, *mocks.AuditRepository)
		wantErr bool
	}{
		{"ok", func(q *mocks.QuickTransactionRepository, a *mocks.AuditRepository) {
			q.On("Create", ctx, mock.AnythingOfType("*model.QuickTransaction")).Return(nil)
			a.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(nil)
		}, false},
		{"create fail", func(q *mocks.QuickTransactionRepository, _ *mocks.AuditRepository) {
			q.On("Create", ctx, mock.AnythingOfType("*model.QuickTransaction")).Return(errors.New("db"))
		}, true},
		{"audit fail", func(q *mocks.QuickTransactionRepository, a *mocks.AuditRepository) {
			q.On("Create", ctx, mock.AnythingOfType("*model.QuickTransaction")).Return(nil)
			a.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(errors.New("db"))
		}, true},
		{"ok2", func(q *mocks.QuickTransactionRepository, a *mocks.AuditRepository) {
			q.On("Create", ctx, mock.AnythingOfType("*model.QuickTransaction")).Return(nil)
			a.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(nil)
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &mocks.QuickTransactionRepository{}
			a := &mocks.AuditRepository{}
			tt.setup(q, a)
			svc := NewQuickTransactionService(q, txServiceStub{}, a)
			_, err := svc.Create(ctx, uid, dto.CreateQuickTransactionRequest{Label: "Tea"})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
