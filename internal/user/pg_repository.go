package user

import (
	"car-mobile-project/internal/models"
	"context"
)

type Repository interface {
	Create(ctx context.Context, user models.User) (*models.User, error)
	GetByPhoneNumber(ctx context.Context, phone string) (*models.User, error)
}
