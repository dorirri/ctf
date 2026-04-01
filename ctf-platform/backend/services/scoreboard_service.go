package services

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type ScoreboardService struct{ db *sqlx.DB }

func NewScoreboardService(db *sqlx.DB) *ScoreboardService {
	return &ScoreboardService{db: db}
}

type ScoreboardEntry struct {
	Rank          int        `json:"rank"`
	Username      string     `json:"username"`
	TotalPoints   int        `json:"total_points"`
	SolvesCount   int        `json:"solves_count"`
	LastSolveTime *time.Time `json:"last_solve_time"`
}

type scoreboardRow struct {
	Username      string     `db:"username"`
	TotalPoints   int        `db:"total_points"`
	SolvesCount   int        `db:"solves_count"`
	LastSolveTime *time.Time `db:"last_solve_time"`
}

func (s *ScoreboardService) GetTop(limit int) ([]ScoreboardEntry, error) {
	rows, err := s.db.Queryx(`
		SELECT u.username,
		       COALESCE(SUM(c.points) FILTER (WHERE s.is_correct), 0)  AS total_points,
		       COUNT(*)              FILTER (WHERE s.is_correct)        AS solves_count,
		       MAX(s.submitted_at)   FILTER (WHERE s.is_correct)        AS last_solve_time
		FROM users u
		LEFT JOIN submissions s ON s.user_id = u.id
		LEFT JOIN challenges c  ON c.id = s.challenge_id
		WHERE u.role = 'player' AND u.is_disabled = false
		GROUP BY u.id, u.username
		ORDER BY total_points DESC, last_solve_time ASC NULLS LAST
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ScoreboardEntry
	rank := 1
	for rows.Next() {
		var row scoreboardRow
		if err := rows.StructScan(&row); err != nil {
			return nil, err
		}
		out = append(out, ScoreboardEntry{
			Rank:          rank,
			Username:      row.Username,
			TotalPoints:   row.TotalPoints,
			SolvesCount:   row.SolvesCount,
			LastSolveTime: row.LastSolveTime,
		})
		rank++
	}
	if out == nil {
		out = []ScoreboardEntry{}
	}
	return out, rows.Err()
}
