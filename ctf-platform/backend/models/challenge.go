package models

import "time"

type Challenge struct {
	ID          int       `db:"id"          json:"id"`
	Title       string    `db:"title"       json:"title"`
	Description string    `db:"description" json:"description"`
	Category    string    `db:"category"    json:"category"`
	Points      int       `db:"points"      json:"points"`
	FlagHash    string    `db:"flag_hash"   json:"-"`
	IsVisible   bool      `db:"is_visible"  json:"is_visible"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
}
