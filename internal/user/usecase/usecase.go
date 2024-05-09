package usecase

import (
	"car-mobile-project/config"
	"car-mobile-project/internal/models"
	"car-mobile-project/internal/user"
	"car-mobile-project/pkg/httpErrors"
	"context"
)

type userUC struct {
	cfg      *config.Config
	userRepo user.Repository
}

func NewUserUseCase(cfg *config.Config, userRepo user.Repository) user.UseCase {
	return &userUC{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (uc *userUC) Create(ctx context.Context, user models.User) (*models.User, error) {
	u, err := uc.userRepo.Create(ctx, user)

	if err != nil {
		return nil, httpErrors.ParseErrors(err)
	}
	return u, nil
}

func (uc *userUC) GetByPhoneNumber(ctx context.Context, phone string) (*models.User, error) {
	u, err := uc.userRepo.GetByPhoneNumber(ctx, phone)

	if err != nil {
		return nil, httpErrors.ParseErrors(err)
	}
	return u, nil
}
