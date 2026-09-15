package pgstore

import (
	"context"
	"errors"
	"fmt"

	"{{MODULE_PATH}}/internal/modules/auth/adapters/pgstore/models"
	"{{MODULE_PATH}}/internal/modules/auth/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueViolationCode   = "23505"
	emailUniqueConstraint = "users_email_unique"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	err := models.New(r.pool).CreateUser(ctx, models.CreateUserParams{
		UserUuid:     user.UUID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	})
	if isUniqueViolation(err, emailUniqueConstraint) {
		return domain.ErrEmailTaken
	}
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByUUID(ctx context.Context, id domain.UserUUID) (*domain.User, error) {
	row, err := models.New(r.pool).GetUserByUUID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("getting user by uuid: %w", err)
	}

	return userFromRow(row), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := models.New(r.pool).GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("getting user by email: %w", err)
	}

	return userFromRow(row), nil
}

func userFromRow(row models.AuthUser) *domain.User {
	return &domain.User{
		UUID:         row.UserUuid,
		Name:         row.Name,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

func isUniqueViolation(err error, constraint string) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	return ok && pgErr.Code == uniqueViolationCode && pgErr.ConstraintName == constraint
}
