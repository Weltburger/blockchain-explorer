package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"uid" json:"uid"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-" `
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
}
