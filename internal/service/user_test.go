package service

import (
	"context"
	"errors"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"ledgerA/internal/service/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserServiceSync(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		uid     string
		setup   func(*mocks.UserRepository)
		wantErr bool
	}{
		{
			name: "update existing",
			uid:  "u1",
			setup: func(repo *mocks.UserRepository) {
				repo.On("FindByFirebaseUID", ctx, "u1").Return(&model.User{FirebaseUID: "u1"}, nil)
				repo.On("Update", ctx, &model.User{FirebaseUID: "u1", DisplayName: "N", Email: "a@b.com"}).Return(nil)
			},
		},
		{
			name: "create new",
			uid:  "u2",
			setup: func(repo *mocks.UserRepository) {
				repo.On("FindByFirebaseUID", ctx, "u2").Return((*model.User)(nil), errors.New("not found"))
				repo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)
			},
		},
		{
			name: "update fails",
			uid:  "u3",
			setup: func(repo *mocks.UserRepository) {
				repo.On("FindByFirebaseUID", ctx, "u3").Return(&model.User{FirebaseUID: "u3"}, nil)
				repo.On("Update", ctx, &model.User{FirebaseUID: "u3", DisplayName: "N", Email: "a@b.com"}).Return(errors.New("db"))
			},
			wantErr: true,
		},
		{
			name: "create fails",
			uid:  "u4",
			setup: func(repo *mocks.UserRepository) {
				repo.On("FindByFirebaseUID", ctx, "u4").Return((*model.User)(nil), errors.New("not found"))
				repo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(errors.New("db"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UserRepository{}
			tt.setup(repo)
			svc := NewUserService(repo)
			_, err := svc.Sync(ctx, tt.uid, dto.SyncRequest{DisplayName: "N", Email: "a@b.com"})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
