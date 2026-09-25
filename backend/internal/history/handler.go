package history

import (
	"context"
	"encoding/json"
	"net/http"

	"calculator-backend/internal/calc"
)

// Store is the persistence dependency Handler needs. *Repository satisfies
// it against a real database; tests can supply a fake instead so handler
// logic is testable without a live Postgres.
type Store interface {
	Create(ctx context.Context, expression, result string) (History, error)
	List(ctx context.Context) ([]History, error)
}

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

type calculateRequest struct {
	Expression string `json:"expression"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	var req calculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := calc.Evaluate(req.Expression)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	entry, err := h.store.Create(r.Context(), req.Expression, result)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save calculation")
		return
	}

	writeJSON(w, http.StatusCreated, entry)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load history")
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}
