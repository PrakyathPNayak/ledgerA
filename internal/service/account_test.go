package service

import (
	"context"
	"errors"
	"ledgerA/internal/dto"
	"ledgerA/internal/service/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAccountServiceCreate(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New()
	tests := []struct {
		name    string
		setup   func(*mocks.AccountRepository, *mocks.AuditRepository)
		wantErr bool
	}{
		{"ok", func(a *mocks.AccountRepository, au *mocks.AuditRepository) {
			a.On("Create", ctx, mock.AnythingOfType("*model.Account")).Return(nil)
			au.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(nil)
		}, false},
		{"account create fail", func(a *mocks.AccountRepository, _ *mocks.AuditRepository) {
			a.On("Create", ctx, mock.AnythingOfType("*model.Account")).Return(errors.New("db"))
		}, true},
		{"audit fail", func(a *mocks.AccountRepository, au *mocks.AuditRepository) {
			a.On("Create", ctx, mock.AnythingOfType("*model.Account")).Return(nil)
			au.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(errors.New("db"))
		}, true},
		{"different name", func(a *mocks.AccountRepository, au *mocks.AuditRepository) {
			a.On("Create", ctx, mock.AnythingOfType("*model.Account")).Return(nil)
			au.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(nil)
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &mocks.AccountRepository{}
			au := &mocks.AuditRepository{}
			tt.setup(a, au)
			svc := NewAccountService(a, au)
			_, err := svc.Create(ctx, uid, dto.CreateAccountRequest{Name: "Wallet", OpeningBalance: 1})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
