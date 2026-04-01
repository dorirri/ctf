package handlers

import (
	"net/http"
	"strconv"

	"ctf-platform/middleware"
	"ctf-platform/services"

	"github.com/go-chi/chi/v5"
)

type ChallengeHandler struct {
	svc *services.ChallengeService
}

func NewChallengeHandler(svc *services.ChallengeService) *ChallengeHandler {
	return &ChallengeHandler{svc: svc}
}

func (h *ChallengeHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	challenges, err := h.svc.ListVisible(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch challenges")
		return
	}
	writeJSON(w, http.StatusOK, challenges)
}

func (h *ChallengeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid challenge id")
		return
	}

	claims := middleware.GetClaims(r)
	ch, err := h.svc.GetByID(id, claims.UserID)
	if err == services.ErrNotFound {
		writeError(w, http.StatusNotFound, "challenge not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch challenge")
		return
	}
	writeJSON(w, http.StatusOK, ch)
}
