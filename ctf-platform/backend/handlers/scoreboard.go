package handlers

import (
	"net/http"

	"ctf-platform/services"
)

type ScoreboardHandler struct {
	svc *services.ScoreboardService
}

func NewScoreboardHandler(svc *services.ScoreboardService) *ScoreboardHandler {
	return &ScoreboardHandler{svc: svc}
}

func (h *ScoreboardHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	entries, err := h.svc.GetTop(20)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch scoreboard")
		return
	}
	writeJSON(w, http.StatusOK, entries)
}
