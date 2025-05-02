package postgres

import (
	"auth-service/internal/services/auth-service/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database.
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsAdmin = false
	user.IsSuperAdmin = false
	user.Role = "visitor"

	query := `INSERT INTO users 
		(id, first_name,last_name,role,is_admin,is_super_admin,email, password_hash, is_verified, is_active, created_at, updated_at) 
		ON CONFLICT (email) DO NOTHING
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsAdmin,
		user.IsSuperAdmin,
		user.Email,
		user.PasswordHash,
		user.IsVerified,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetByEmail retrieves a user by email from the database.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, first_name, last_name, email, password_hash, role, is_admin, is_super_admin, is_verified, is_active, created_at, updated_at 
		FROM users WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsAdmin,
		&user.IsSuperAdmin,
		&user.IsVerified,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user by ID from the database.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, first_name, last_name, email, password_hash, role, is_admin, is_super_admin, is_verified, is_active, created_at, updated_at 
		FROM users WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsAdmin,
		&user.IsSuperAdmin,
		&user.IsVerified,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user in the database.
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users 
		SET first_name = $1, last_name = $2, email = $3, password_hash = $4, is_verified = $5, is_active = $6, updated_at = $7 
		WHERE id = $8`

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.IsVerified,
		user.IsActive,
		time.Now(),
		user.ID,
	)

	return err
}

func (r *UserRepository) FindOrCreateOAuthUser(ctx context.Context, provider, providerID, email string) (*models.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check if the user already exists
	query := `SELECT u.id, u.email, u.is_verified, u.is_active 
		FROM users u
		JOIN oauth_providers op ON u.id = op.user_id
		WHERE op.provider = $1 AND op.provider_id = $2`

	row := tx.QueryRowContext(ctx, query, provider, providerID)

	var user models.User

	err = row.Scan(&user.ID, &user.Email, &user.IsVerified, &user.IsActive)
	if err == nil {
		return &user, tx.Commit()
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// If the user doesn't exist, create a new user

	firstName, lastName := ExtractPotentialNames(email) // This function should be defined to extract names from the email

	user.ID = uuid.New().String()
	user.Email = email
	user.FirstName = firstName
	user.LastName = lastName
	user.Role = "visitor"
	user.IsAdmin = false
	user.IsSuperAdmin = false
	user.IsVerified = true // OAuth users are automatically verified
	user.IsActive = true
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO users (id,
			first_name,
			last_name,
			role,
			is_admin,
			is_super_admin, 
			email, 
			is_verified,
			is_active,
			created_at, 
			updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsAdmin,
		user.IsSuperAdmin,
		user.Email,
		user.IsVerified,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Link OAuth provider
	_, err = tx.ExecContext(ctx,
		`INSERT INTO oauth_providers (user_id, provider, provider_id,created_at) 
		VALUES ($1, $2, $3)`,
		user.ID,
		provider,
		providerID,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &user, tx.Commit()
}
