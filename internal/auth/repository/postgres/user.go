package postgres

import (
	"context"
	"log"

	"explorer/internal/apperrors"
	"explorer/internal/auth"
	"explorer/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// PGUserRepository is data/repository implementation
// of service layer UserRepository
type PGUserRepository struct {
	DB *sqlx.DB
}

// NewUserRepository is a factory for initializing User Repositories
func NewUserRepository(db *sqlx.DB) auth.UserRepo {
	return &PGUserRepository{
		DB: db,
	}
}

// Create reaches out to database SQLX api
func (r *PGUserRepository) CreateUser(ctx context.Context, u *models.User) error {
	query := "INSERT INTO users (uid, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING *"

	if _, err := r.DB.ExecContext(ctx, query, u.ID, u.Email, u.Password, u.CreatedAt); err != nil {
		// check unique constraint
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err.Code.Name())
			return apperrors.NewConflict("email", u.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err)
		return apperrors.NewInternal()
	}
	return nil
}

// FindByEmail retrieves user row by email address
func (r *PGUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	query := "SELECT * FROM users WHERE email=$1"

	if err := r.DB.GetContext(ctx, user, query, email); err != nil {
		log.Printf("Unable to get user with email address: %v. Err: %v\n", email, err)
		return user, apperrors.NewNotFound("email", email)
	}

	return user, nil
}

// FindByID fetches user by id
func (r *PGUserRepository) GetByID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	user := &models.User{}

	query := "SELECT * FROM users WHERE uid=$1"

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.GetContext(ctx, user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())
	}

	return user, nil
}
