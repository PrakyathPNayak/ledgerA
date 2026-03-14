package service

import (
	"context"
	"errors"
	"ledgerA/internal/model"
	"ledgerA/internal/service/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSubcategoryServiceGetOrCreateFromCategory(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New()
	cid := uuid.New()
	tests := []struct {
		name    string
		setup   func(*mocks.SubcategoryRepository, *mocks.CategoryRepository, *mocks.AuditRepository)
		wantErr bool
	}{
		{"found existing", func(s *mocks.SubcategoryRepository, _ *mocks.CategoryRepository, _ *mocks.AuditRepository) {
			s.On("FindByNameForCategory", ctx, uid, cid, "Food").Return(&model.Subcategory{Name: "Food"}, nil)
		}, false},
		{"create when missing", func(s *mocks.SubcategoryRepository, _ *mocks.CategoryRepository, a *mocks.AuditRepository) {
			s.On("FindByNameForCategory", ctx, uid, cid, "Food").Return((*model.Subcategory)(nil), errors.New("missing"))
			s.On("Create", ctx, mock.AnythingOfType("*model.Subcategory")).Return(nil)
			a.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(nil)
		}, false},
		{"create fail", func(s *mocks.SubcategoryRepository, _ *mocks.CategoryRepository, _ *mocks.AuditRepository) {
			s.On("FindByNameForCategory", ctx, uid, cid, "Food").Return((*model.Subcategory)(nil), errors.New("missing"))
			s.On("Create", ctx, mock.AnythingOfType("*model.Subcategory")).Return(errors.New("db"))
		}, true},
		{"audit fail", func(s *mocks.SubcategoryRepository, _ *mocks.CategoryRepository, a *mocks.AuditRepository) {
			s.On("FindByNameForCategory", ctx, uid, cid, "Food").Return((*model.Subcategory)(nil), errors.New("missing"))
			s.On("Create", ctx, mock.AnythingOfType("*model.Subcategory")).Return(nil)
			a.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(errors.New("db"))
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mocks.SubcategoryRepository{}
			c := &mocks.CategoryRepository{}
			a := &mocks.AuditRepository{}
			tt.setup(s, c, a)
			svc := NewSubcategoryService(s, c, a)
			_, err := svc.GetOrCreateFromCategory(ctx, uid, cid, "Food")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
