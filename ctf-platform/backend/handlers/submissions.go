package handlers

import (
	"encoding/json"
	"net/http"

	"ctf-platform/middleware"
	"ctf-platform/services"
)

type SubmissionHandler struct {
	svc *services.SubmissionService
}

func NewSubmissionHandler(svc *services.SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{svc: svc}
}

type submitRequest struct {
	ChallengeID int    `json:"challenge_id"`
	Flag        string `json:"flag"`
}

func (h *SubmissionHandler) Submit(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims.Role == "admin" {
		writeError(w, http.StatusForbidden, "admins cannot submit flags")
		return
	}

	var req submitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ChallengeID == 0 || req.Flag == "" {
		writeError(w, http.StatusBadRequest, "challenge_id and flag are required")
		return
	}

	result, err := h.svc.Submit(claims.UserID, req.ChallengeID, req.Flag)
	if err == services.ErrNotFound {
		writeError(w, http.StatusNotFound, "challenge not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not process submission")
		return
	}

	if result.AlreadySolved {
		writeError(w, http.StatusConflict, "already solved")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"correct": result.Correct,
		"points":  result.Points,
	})
}
