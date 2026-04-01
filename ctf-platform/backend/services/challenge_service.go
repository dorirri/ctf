package services

import "github.com/jmoiron/sqlx"

type ChallengeService struct{ db *sqlx.DB }

func NewChallengeService(db *sqlx.DB) *ChallengeService {
	return &ChallengeService{db: db}
}

type ChallengeDTO struct {
	ID       int    `db:"id"       json:"id"`
	Title    string `db:"title"    json:"title"`
	Category string `db:"category" json:"category"`
	Points   int    `db:"points"   json:"points"`
	IsSolved bool   `db:"is_solved" json:"is_solved"`
}

type ChallengeDTOFull struct {
	ChallengeDTO
	Description string `db:"description" json:"description"`
}

func (s *ChallengeService) ListVisible(userID int) ([]ChallengeDTO, error) {
	rows, err := s.db.Queryx(`
		SELECT c.id, c.title, c.category, c.points,
		       EXISTS (
		           SELECT 1 FROM submissions sub
		           WHERE sub.challenge_id = c.id
		             AND sub.user_id = $1
		             AND sub.is_correct = true
		       ) AS is_solved
		FROM challenges c
		WHERE c.is_visible = true
		ORDER BY c.points ASC, c.id ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ChallengeDTO
	for rows.Next() {
		var ch ChallengeDTO
		if err := rows.StructScan(&ch); err != nil {
			return nil, err
		}
		out = append(out, ch)
	}
	if out == nil {
		out = []ChallengeDTO{}
	}
	return out, rows.Err()
}

func (s *ChallengeService) GetByID(id int, userID int) (*ChallengeDTOFull, error) {
	var ch ChallengeDTOFull
	err := s.db.QueryRowx(`
		SELECT c.id, c.title, c.description, c.category, c.points,
		       EXISTS (
		           SELECT 1 FROM submissions sub
		           WHERE sub.challenge_id = c.id
		             AND sub.user_id = $2
		             AND sub.is_correct = true
		       ) AS is_solved
		FROM challenges c
		WHERE c.id = $1 AND c.is_visible = true
	`, id, userID).StructScan(&ch)
	if err != nil {
		return nil, ErrNotFound
	}
	return &ch, nil
}
