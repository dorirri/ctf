package models

import "time"

type Team struct {
	ID         int       `db:"id"          json:"id"`
	Name       string    `db:"name"        json:"name"`
	InviteCode string    `db:"invite_code" json:"invite_code"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
}

type TeamMember struct {
	UserID int `db:"user_id" json:"user_id"`
	TeamID int `db:"team_id" json:"team_id"`
}
