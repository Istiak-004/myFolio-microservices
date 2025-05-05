package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/istiak-004/myFolio-microservices/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type BaseRepository struct {
	// DB is the database connection
	db     *sqlx.DB
	logger *logger.Logger
}

// NewBaseRepository creates a new BaseRepository
func NewBaseRepository(db *sqlx.DB, logger *logger.Logger) *BaseRepository {
	return &BaseRepository{
		db:     db,
		logger: logger,
	}
}

// WithTransaction executes a function within a transaction
// It begins a transaction, executes the function, and commits or rolls back the transaction
func (r *BaseRepository) WithTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil) // Begin a new transaction
	if err != nil {
		r.logger.WithError(err).Error("Failed to begin transaction!")
		return err
	}

	// Ensure that the transaction is rolled back in case of panic
	defer func() {
		if p := recover(); p != nil {
			r.logger.WithError(err).Error("Panic occurred in transaction!")
			if err := tx.Rollback(); err != nil {
				r.logger.WithError(err).Error("Failed to rollback transaction!")
			}
			panic(p)
		}
	}()

	// Execute the function with the transaction
	// If the function returns an error, rollback the transaction
	if err := fn(tx); err != nil {
		r.logger.WithError(err).Error("Transaction failed!")
		if rbErr := tx.Rollback(); rbErr != nil {
			r.logger.WithError(err).Error("Failed to rollback transaction!")
			return fmt.Errorf("transaction failed: %w, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	// If the function succeeds, commit the transaction
	if err := tx.Commit(); err != nil {
		r.logger.WithError(err).Error("Failed to commit transaction!")
		return err
	}
	return nil
}

// Get executes a query and scans the result into destination
// It is used for SELECT queries that return a single row
func (r *BaseRepository) Get(ctx context.Context, dest any, query string, args ...any) error {
	err := r.db.GetContext(ctx, dest, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to execute query!")
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

// Select executes a query and scans the result into destination
// It is used for SELECT queries that return multiple rows
func (r *BaseRepository) Select(ctx context.Context, dest any, query string, args ...any) error {
	err := r.db.SelectContext(ctx, dest, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to execute query!")
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

// Exec executes a query without returning any rows
// It is used for INSERT, UPDATE, DELETE queries
func (r *BaseRepository) Exec(ctx context.Context, query string, args ...any) error {
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to execute query!")
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

// QueryRow executes a query and scans the result into destination
// It is used for SELECT queries that return a single row
func (r *BaseRepository) QueryRow(ctx context.Context, dest any, query string, args ...any) error {
	err := r.db.QueryRowContext(ctx, query, args...).Scan(dest)
	if err != nil {
		r.logger.WithError(err).Error("Failed to execute query!")
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

// Query executes a query and returns the rows
// It is used for SELECT queries that return multiple rows
// Note: You need to close the rows after using them
func (r *BaseRepository) Query(ctx context.Context, dest any, query string, args ...any) (*sql.Rows, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to execute query!")
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	return rows, nil
}
