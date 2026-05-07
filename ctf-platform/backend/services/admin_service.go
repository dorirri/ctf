package services

import (
	"time"

	"ctf-platform/models"
	"ctf-platform/utils"

	"github.com/jmoiron/sqlx"
)

type AdminService struct{ db *sqlx.DB }

func NewAdminService(db *sqlx.DB) *AdminService {
	return &AdminService{db: db}
}

type DefaultAdminInput struct {
	Username string
	Email    string
	Password string
}

type CreateChallengeInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Points      int    `json:"points"`
	Flag        string `json:"flag"`
	IsVisible   bool   `json:"is_visible"`
}

type UpdateChallengeInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Points      int    `json:"points"`
	IsVisible   bool   `json:"is_visible"`
}

type SubmissionView struct {
	Username       string    `db:"username"        json:"username"`
	ChallengeTitle string    `db:"challenge_title" json:"challenge_title"`
	IsCorrect      bool      `db:"is_correct"      json:"is_correct"`
	SubmittedAt    time.Time `db:"submitted_at"    json:"submitted_at"`
}

type UserView struct {
	ID         int       `db:"id"          json:"id"`
	Username   string    `db:"username"    json:"username"`
	Email      string    `db:"email"       json:"email"`
	Role       string    `db:"role"        json:"role"`
	IsDisabled bool      `db:"is_disabled" json:"is_disabled"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
}

type ChallengeAdminView struct {
	ID          int       `db:"id"          json:"id"`
	Title       string    `db:"title"       json:"title"`
	Description string    `db:"description" json:"description"`
	Category    string    `db:"category"    json:"category"`
	Points      int       `db:"points"      json:"points"`
	IsVisible   bool      `db:"is_visible"  json:"is_visible"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
}

func (s *AdminService) ListChallenges() ([]ChallengeAdminView, error) {
	rows, err := s.db.Queryx(`
		SELECT id, title, description, category, points, is_visible, created_at
		FROM challenges
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ChallengeAdminView
	for rows.Next() {
		var ch ChallengeAdminView
		if err := rows.StructScan(&ch); err != nil {
			return nil, err
		}
		out = append(out, ch)
	}
	if out == nil {
		out = []ChallengeAdminView{}
	}
	return out, rows.Err()
}

func (s *AdminService) CreateChallenge(input CreateChallengeInput) (*models.Challenge, error) {
	hash, err := utils.HashPassword(input.Flag)
	if err != nil {
		return nil, err
	}

	var ch models.Challenge
	err = s.db.QueryRowx(`
		INSERT INTO challenges (title, description, category, points, flag_hash, is_visible)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, category, points, flag_hash, is_visible, created_at
	`, input.Title, input.Description, input.Category, input.Points, hash, input.IsVisible,
	).StructScan(&ch)
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *AdminService) UpdateChallenge(id int, input UpdateChallengeInput) error {
	res, err := s.db.Exec(`
		UPDATE challenges
		SET title=$1, description=$2, category=$3, points=$4, is_visible=$5
		WHERE id=$6
	`, input.Title, input.Description, input.Category, input.Points, input.IsVisible, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *AdminService) DeleteChallenge(id int) error {
	_, err := s.db.Exec(`DELETE FROM submissions WHERE challenge_id = $1`, id)
	if err != nil {
		return err
	}
	res, err := s.db.Exec(`DELETE FROM challenges WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *AdminService) ToggleVisibility(id int) error {
	res, err := s.db.Exec(
		`UPDATE challenges SET is_visible = NOT is_visible WHERE id = $1`, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *AdminService) ListSubmissions() ([]SubmissionView, error) {
	rows, err := s.db.Queryx(`
		SELECT u.username,
		       c.title AS challenge_title,
		       s.is_correct,
		       s.submitted_at
		FROM submissions s
		JOIN users      u ON u.id = s.user_id
		JOIN challenges c ON c.id = s.challenge_id
		ORDER BY s.submitted_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SubmissionView
	for rows.Next() {
		var sv SubmissionView
		if err := rows.StructScan(&sv); err != nil {
			return nil, err
		}
		out = append(out, sv)
	}
	if out == nil {
		out = []SubmissionView{}
	}
	return out, rows.Err()
}

func (s *AdminService) ToggleUserDisabled(id int) error {
	res, err := s.db.Exec(
		`UPDATE users SET is_disabled = NOT is_disabled WHERE id = $1`, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *AdminService) ListUsers() ([]UserView, error) {
	rows, err := s.db.Queryx(`
		SELECT id, username, email, role, is_disabled, created_at
		FROM users
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []UserView
	for rows.Next() {
		var uv UserView
		if err := rows.StructScan(&uv); err != nil {
			return nil, err
		}
		out = append(out, uv)
	}
	if out == nil {
		out = []UserView{}
	}
	return out, rows.Err()
}

func (s *AdminService) EnsureDefaultAdmin(input DefaultAdminInput) error {
	if input.Email == "" || input.Password == "" {
		return nil
	}

	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, input.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	username := input.Username
	if username == "" {
		username = "admin"
	}

	_, err = s.db.Exec(`
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, 'admin')
	`, username, input.Email, hash)
	return err
}

func (s *AdminService) Seed() error {
	challenges := []struct {
		title, description, category, flag string
		points                             int
	}{
		{"Hello World", "Find the flag hidden in the response headers.", "Web", "CTF{hello_world}", 50},
		{"Base64 Decode", "Decode this: Q1RGe2Jhc2U2NF9pc19ub3RfZW5jcnlwdGlvbn0=", "Crypto", "CTF{base64_is_not_encryption}", 100},
		{"SQL Basics", "Use SQL injection to bypass the login form.", "Web", "CTF{sql1_byp4ss}", 200},
	}

	for _, ch := range challenges {
		var exists bool
		if err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM challenges WHERE title = $1)`, ch.title).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		hash, err := utils.HashPassword(ch.flag)
		if err != nil {
			return err
		}
		_, err = s.db.Exec(`
			INSERT INTO challenges (title, description, category, points, flag_hash, is_visible)
			VALUES ($1, $2, $3, $4, $5, true)
		`, ch.title, ch.description, ch.category, ch.points, hash)
		if err != nil {
			return err
		}
	}

	return nil
}
