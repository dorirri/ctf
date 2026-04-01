package models

import "time"

type Submission struct {
	ID          int       `db:"id"           json:"id"`
	UserID      int       `db:"user_id"      json:"user_id"`
	ChallengeID int       `db:"challenge_id" json:"challenge_id"`
	IsCorrect   bool      `db:"is_correct"   json:"is_correct"`
	SubmittedAt time.Time `db:"submitted_at" json:"submitted_at"`
}
