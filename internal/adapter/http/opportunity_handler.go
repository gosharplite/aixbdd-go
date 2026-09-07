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

type OpportunityHandler struct {
	svc *service.OpportunityService
}

func NewOpportunityHandler(svc *service.OpportunityService) *OpportunityHandler {
	return &OpportunityHandler{svc: svc}
}

func (h *OpportunityHandler) CreateOpportunity(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	var req service.CreateOpportunityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的 JSON 內容")
		return
	}

	opp, err := h.svc.CreateOpportunity(r.Context(), actor, req)
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
	RespondJSON(w, http.StatusCreated, opp)
}

func (h *OpportunityHandler) GetOpportunity(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的機會識別碼")
		return
	}

	opp, err := h.svc.GetOpportunity(r.Context(), actor, id)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", "找不到指定的銷售機會")
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, opp)
}

func (h *OpportunityHandler) UpdateStage(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的機會識別碼")
		return
	}

	var req service.UpdateStageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的 JSON 內容")
		return
	}

	opp, err := h.svc.UpdateStage(r.Context(), actor, id, req)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			RespondError(w, http.StatusForbidden, "FORBIDDEN", "只有銷售機會負責人或主管可以推進階段")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", "找不到指定的銷售機會")
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, opp)
}

func (h *OpportunityHandler) ListOpportunities(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	var stagePtr *domain.SalesStage
	if st := r.URL.Query().Get("stage"); st != "" {
		s := parseSalesStage(st)
		stagePtr = &s
	}
	var ownerIDPtr *int64
	if oidStr := r.URL.Query().Get("owner_id"); oidStr != "" {
		if oid, err := strconv.ParseInt(oidStr, 10, 64); err == nil {
			ownerIDPtr = &oid
		}
	}

	list, err := h.svc.ListOpportunities(r.Context(), actor, stagePtr, ownerIDPtr)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if list == nil {
		list = []domain.Opportunity{}
	}
	RespondJSON(w, http.StatusOK, list)
}

func parseSalesStage(st string) domain.SalesStage {
	switch st {
	case "潛在", "Prospecting":
		return domain.StageProspecting
	case "已聯繫", "Contacted":
		return domain.StageContacted
	case "提案中", "Proposal":
		return domain.StageProposal
	case "已成交", "ClosedWon":
		return domain.StageClosedWon
	case "已失敗", "ClosedLost":
		return domain.StageClosedLost
	default:
		return domain.SalesStage(st)
	}
}
