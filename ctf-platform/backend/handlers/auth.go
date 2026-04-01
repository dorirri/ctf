package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"ctf-platform/middleware"
	"ctf-platform/models"
	"ctf-platform/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type AuthHandler struct {
	db        *sqlx.DB
	jwtSecret string
}

func NewAuthHandler(db *sqlx.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{db: db, jwtSecret: jwtSecret}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username, email, and password are required")
		return
	}
	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	var user models.User
	err = h.db.QueryRowx(
		`INSERT INTO users (username, email, password_hash, role)
		 VALUES ($1, $2, $3, 'player')
		 RETURNING id, username, email, password_hash, role, is_disabled, created_at`,
		req.Username, req.Email, hash,
	).StructScan(&user)
	if err != nil {
		// unique constraint violations surface as pq error code 23505
		writeError(w, http.StatusConflict, "username or email already taken")
		return
	}

	token, err := h.issueToken(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not issue token")
		return
	}

	writeJSON(w, http.StatusCreated, tokenResponse{Token: token, User: user})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	var user models.User
	err := h.db.QueryRowx(
		`SELECT id, username, email, password_hash, role, is_disabled, created_at
		 FROM users WHERE email = $1`, req.Email,
	).StructScan(&user)
	if err != nil {
		// deliberately vague to avoid user enumeration
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if user.IsDisabled {
		writeError(w, http.StatusForbidden, "account is disabled")
		return
	}

	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := h.issueToken(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not issue token")
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{Token: token, User: user})
}

func (h *AuthHandler) issueToken(u models.User) (string, error) {
	claims := middleware.Claims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.jwtSecret))
}
