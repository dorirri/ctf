package models

import "time"

type User struct {
	ID           int       `db:"id"            json:"id"`
	Username     string    `db:"username"      json:"username"`
	Email        string    `db:"email"         json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role"          json:"role"`
	IsDisabled   bool      `db:"is_disabled"   json:"is_disabled"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
}
