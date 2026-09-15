package pgstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"{{MODULE_PATH}}/internal/modules/auth/adapters/pgstore/models"
	"{{MODULE_PATH}}/internal/modules/auth/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{pool: pool}
}

func (r *TokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	if err := createRefreshToken(ctx, models.New(r.pool), token); err != nil {
		return fmt.Errorf("creating refresh token: %w", err)
	}

	return nil
}

func (r *TokenRepository) GetByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	row, err := models.New(r.pool).GetRefreshTokenByHash(ctx, hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("getting refresh token: %w", err)
	}

	token := &domain.RefreshToken{
		UUID:      row.TokenUuid,
		UserUUID:  row.UserUuid,
		Hash:      row.TokenHash,
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}
	if row.RevokedAt != nil {
		token.RevokedAt = *row.RevokedAt
	}

	return token, nil
}

func (r *TokenRepository) Revoke(ctx context.Context, token *domain.RefreshToken) error {
	err := models.New(r.pool).RevokeRefreshToken(ctx, models.RevokeRefreshTokenParams{
		TokenUuid: token.UUID,
		RevokedAt: new(token.RevokedAt),
	})
	if err != nil {
		return fmt.Errorf("revoking refresh token: %w", err)
	}

	return nil
}

// Rotate revokes current only while it is still active and stores the
// replacement in the same transaction. The conditional update takes the row
// lock, so two concurrent rotations of one token cannot both succeed.
func (r *TokenRepository) Rotate(ctx context.Context, current, replacement *domain.RefreshToken) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("beginning rotation: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	queries := models.New(tx)

	rows, err := queries.RevokeRefreshTokenIfActive(ctx, models.RevokeRefreshTokenIfActiveParams{
		TokenUuid: current.UUID,
		RevokedAt: new(current.RevokedAt),
	})
	if err != nil {
		return fmt.Errorf("revoking rotated refresh token: %w", err)
	}

	if rows == 0 {
		return domain.ErrRefreshTokenRotated
	}

	if err := createRefreshToken(ctx, queries, replacement); err != nil {
		return fmt.Errorf("creating replacement refresh token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing rotation: %w", err)
	}

	return nil
}

func (r *TokenRepository) RevokeAllForUser(ctx context.Context, userUUID domain.UserUUID, now time.Time) error {
	err := models.New(r.pool).RevokeAllUserRefreshTokens(ctx, models.RevokeAllUserRefreshTokensParams{
		UserUuid:  userUUID,
		RevokedAt: new(now),
	})
	if err != nil {
		return fmt.Errorf("revoking user refresh tokens: %w", err)
	}

	return nil
}

func createRefreshToken(ctx context.Context, queries *models.Queries, token *domain.RefreshToken) error {
	return queries.CreateRefreshToken(ctx, models.CreateRefreshTokenParams{
		TokenUuid: token.UUID,
		UserUuid:  token.UserUUID,
		TokenHash: token.Hash,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: token.CreatedAt,
	})
}
