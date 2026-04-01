package services

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type SubmissionService struct{ db *sqlx.DB }

func NewSubmissionService(db *sqlx.DB) *SubmissionService {
	return &SubmissionService{db: db}
}

type SubmitResult struct {
	Correct      bool
	AlreadySolved bool
	Points        int
}

func (s *SubmissionService) Submit(userID int, challengeID int, flag string) (SubmitResult, error) {
	var flagHash string
	var points int
	err := s.db.QueryRow(
		`SELECT flag_hash, points FROM challenges WHERE id = $1 AND is_visible = true`,
		challengeID,
	).Scan(&flagHash, &points)
	if err == sql.ErrNoRows {
		return SubmitResult{}, ErrNotFound
	}
	if err != nil {
		return SubmitResult{}, err
	}

	var alreadySolved bool
	err = s.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM submissions
			WHERE user_id = $1 AND challenge_id = $2 AND is_correct = true
		)`, userID, challengeID,
	).Scan(&alreadySolved)
	if err != nil {
		return SubmitResult{}, err
	}
	if alreadySolved {
		return SubmitResult{AlreadySolved: true}, nil
	}

	correct := bcrypt.CompareHashAndPassword([]byte(flagHash), []byte(flag)) == nil

	_, err = s.db.Exec(
		`INSERT INTO submissions (user_id, challenge_id, is_correct) VALUES ($1, $2, $3)`,
		userID, challengeID, correct,
	)
	if err != nil {
		return SubmitResult{}, err
	}

	if correct {
		return SubmitResult{Correct: true, Points: points}, nil
	}
	return SubmitResult{Correct: false}, nil
}
