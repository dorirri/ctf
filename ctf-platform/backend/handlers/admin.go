package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ctf-platform/services"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	svc *services.AdminService
}

func NewAdminHandler(svc *services.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) CreateChallenge(w http.ResponseWriter, r *http.Request) {
	var input services.CreateChallengeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Title == "" || input.Flag == "" {
		writeError(w, http.StatusBadRequest, "title and flag are required")
		return
	}

	ch, err := h.svc.CreateChallenge(input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create challenge")
		return
	}
	writeJSON(w, http.StatusCreated, ch)
}

func (h *AdminHandler) UpdateChallenge(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid challenge id")
		return
	}

	var input services.UpdateChallengeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.UpdateChallenge(id, input); err == services.ErrNotFound {
		writeError(w, http.StatusNotFound, "challenge not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update challenge")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) DeleteChallenge(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid challenge id")
		return
	}

	if err := h.svc.DeleteChallenge(id); err == services.ErrNotFound {
		writeError(w, http.StatusNotFound, "challenge not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete challenge")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) ToggleVisibility(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid challenge id")
		return
	}

	if err := h.svc.ToggleVisibility(id); err == services.ErrNotFound {
		writeError(w, http.StatusNotFound, "challenge not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "could not toggle visibility")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.svc.ListSubmissions()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch submissions")
		return
	}
	writeJSON(w, http.StatusOK, subs)
}

func (h *AdminHandler) DisableUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.svc.ToggleUserDisabled(id); err == services.ErrNotFound {
		writeError(w, http.StatusNotFound, "user not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "could not toggle user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch users")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *AdminHandler) Seed(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Seed(); err != nil {
		writeError(w, http.StatusInternalServerError, "seed failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "seeded"})
}
