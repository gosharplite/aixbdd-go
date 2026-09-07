package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"crm/internal/domain"
	"crm/internal/service"
	"github.com/go-chi/chi/v5"
)

type InteractionHandler struct {
	svc *service.InteractionService
}

func NewInteractionHandler(svc *service.InteractionService) *InteractionHandler {
	return &InteractionHandler{svc: svc}
}

func (h *InteractionHandler) CreateInteraction(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	idStr := chi.URLParam(r, "id")
	contactID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的聯絡人識別碼")
		return
	}

	var req service.CreateInteractionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的 JSON 內容")
		return
	}

	inter, err := h.svc.CreateInteraction(r.Context(), actor, contactID, req)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			RespondError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, inter)
}

func (h *InteractionHandler) ListInteractions(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	idStr := chi.URLParam(r, "id")
	contactID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的聯絡人識別碼")
		return
	}

	list, err := h.svc.ListInteractions(r.Context(), actor, contactID)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", "找不到指定的聯絡人")
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if list == nil {
		list = []domain.Interaction{}
	}
	RespondJSON(w, http.StatusOK, list)
}
