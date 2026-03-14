package service

import (
	"context"
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"ledgerA/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Sync(ctx context.Context, firebaseUID string, req dto.SyncRequest) (*model.User, error) {
	user, err := s.userRepo.FindByFirebaseUID(ctx, firebaseUID)
	if err == nil {
		user.DisplayName = req.DisplayName
		user.Email = req.Email
		if errUpdate := s.userRepo.Update(ctx, user); errUpdate != nil {
			return nil, fmt.Errorf("userService.Sync.Update: %w", errUpdate)
		}
		return user, nil
	}

	newUser := req.ToModel(firebaseUID, "")
	if newUser.CurrencyCode == "" {
		newUser.CurrencyCode = "INR"
	}
	if errCreate := s.userRepo.Create(ctx, &newUser); errCreate != nil {
		return nil, fmt.Errorf("userService.Sync.Create: %w", errCreate)
	}
	return &newUser, nil
}

func (s *userService) GetMe(ctx context.Context, firebaseUID string) (*model.User, error) {
	user, err := s.userRepo.FindByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("userService.GetMe: %w", err)
	}
	return user, nil
}

func (s *userService) UpdateMe(ctx context.Context, firebaseUID string, req dto.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("userService.UpdateMe.Find: %w", err)
	}
	user.DisplayName = req.DisplayName
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("userService.UpdateMe.Update: %w", err)
	}
	return user, nil
}
