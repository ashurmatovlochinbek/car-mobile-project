package repository

import (
	"car-mobile-project/internal/models"
	"car-mobile-project/internal/user"
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) user.Repository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user models.User) (*models.User, error) {
	newUUID := uuid.New()
	u := &models.User{}
	createUser := `INSERT INTO users (user_id, name, phone_number) VALUES ($1, $2, $3) RETURNING *;`

	if err := r.db.QueryRowxContext(
		ctx,
		createUser,
		&newUUID,
		&user.Name,
		&user.PhoneNumber,
	).Scan(&u.UserId, &u.Name, &u.PhoneNumber); err != nil {
		return nil, errors.Wrap(err, "userRepo.Create.StructScan")
	}
	return u, nil
}

func (r *userRepo) GetByPhoneNumber(ctx context.Context, phone string) (*models.User, error) {
	getByPhoneNumber := `SELECT * FROM users WHERE phone_number = $1;`
	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getByPhoneNumber, phone).Scan(&u.UserId, &u.Name, &u.PhoneNumber); err != nil {
		return nil, errors.Wrap(err, "userRepo.GetByPhoneNumber")
	}
	return u, nil
}
