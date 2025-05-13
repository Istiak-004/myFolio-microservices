package postgres

import (
	"context"

	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/models"
	"github.com/istiak-004/myFolio-microservices/pkg/database"
	"github.com/jmoiron/sqlx"
)

type VerificationRepository struct {
	db *sqlx.DB
}

func NewVerificationRepository(db *database.Client) *VerificationRepository {
	return &VerificationRepository{db: db.GetDB()}
}

func (r *VerificationRepository) Create(ctx context.Context, token *models.VerificationToken) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO verification_tokens (token, user_id, expires_at)
        VALUES ($1, $2, $3)`,
		token.Token, token.UserID, token.ExpiresAt,
	)
	return err
}

func (r *VerificationRepository) Get(ctx context.Context, token string) (*models.VerificationToken, error) {
	t := &models.VerificationToken{}
	err := r.db.QueryRowContext(ctx, `SELECT token, user_id, expires_at FROM verification_tokens WHERE token = $1`, token).
		Scan(&t.Token, &t.UserID, &t.ExpiresAt)
	return t, err
}

func (r *VerificationRepository) Delete(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM verification_tokens WHERE token = $1`, token)
	return err
}
