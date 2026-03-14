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

func TestCategoryServiceCreate(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New()
	tests := []struct {
		name    string
		setup   func(*mocks.CategoryRepository, *mocks.AuditRepository)
		wantErr bool
	}{
		{"ok", func(c *mocks.CategoryRepository, a *mocks.AuditRepository) {
			c.On("Create", ctx, mock.AnythingOfType("*model.Category")).Return(nil)
			a.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(nil)
		}, false},
		{"create fail", func(c *mocks.CategoryRepository, _ *mocks.AuditRepository) {
			c.On("Create", ctx, mock.AnythingOfType("*model.Category")).Return(errors.New("db"))
		}, true},
		{"audit fail", func(c *mocks.CategoryRepository, a *mocks.AuditRepository) {
			c.On("Create", ctx, mock.AnythingOfType("*model.Category")).Return(nil)
			a.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(errors.New("db"))
		}, true},
		{"ok2", func(c *mocks.CategoryRepository, a *mocks.AuditRepository) {
			c.On("Create", ctx, mock.AnythingOfType("*model.Category")).Return(nil)
			a.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(nil)
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &mocks.CategoryRepository{}
			a := &mocks.AuditRepository{}
			tt.setup(c, a)
			svc := NewCategoryService(c, a)
			_, err := svc.Create(ctx, uid, dto.CreateCategoryRequest{Name: "Food"})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
