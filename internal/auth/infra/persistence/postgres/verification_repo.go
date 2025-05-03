package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type VerificationRepository struct {
	db *sql.DB
}

func NewVerificationRepository(db *sql.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) CreateEmailVerification(ctx context.Context, userID, email, token string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO email_verifications 
		(id, user_id, email, token, expires_at) 
		VALUES ($1, $2, $3, $4, $5)`,
		uuid.New().String(),
		userID,
		email,
		token,
		time.Now().Add(24*time.Hour),
	)
	return err
}

func (r *VerificationRepository) VerifyEmail(ctx context.Context, token string) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Get verification record
	var userID, email string
	err = tx.QueryRowContext(ctx,
		`SELECT user_id, email FROM email_verifications 
		WHERE token = $1 AND expires_at > NOW() AND used_at IS NULL`,
		token,
	).Scan(&userID, &email)
	if err != nil {
		return "", err
	}

	// Mark as used
	_, err = tx.ExecContext(ctx,
		`UPDATE email_verifications SET used_at = NOW() WHERE token = $1`,
		token,
	)
	if err != nil {
		return "", err
	}

	// Update user
	_, err = tx.ExecContext(ctx,
		`UPDATE users SET is_verified = true, updated_at = NOW() WHERE id = $1`,
		userID,
	)
	if err != nil {
		return "", err
	}

	return userID, tx.Commit()
}
