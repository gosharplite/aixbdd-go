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

type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

func (h *CustomerHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	var req service.CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的 JSON 內容")
		return
	}
	comp, err := h.svc.CreateCompany(r.Context(), actor, req)
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
	RespondJSON(w, http.StatusCreated, comp)
}

func (h *CustomerHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	companies, err := h.svc.ListCompanies(r.Context(), actor)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, companies)
}

func (h *CustomerHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的公司識別碼")
		return
	}
	comp, contacts, err := h.svc.GetCompany(r.Context(), actor, id)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", "找不到指定的公司")
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	type resp struct {
		*domain.Company
		Contacts []domain.Contact `json:"contacts"`
	}
	RespondJSON(w, http.StatusOK, resp{Company: comp, Contacts: contacts})
}

func (h *CustomerHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	var req service.CreateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的 JSON 內容")
		return
	}
	contact, err := h.svc.CreateContact(r.Context(), actor, req)
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
	RespondJSON(w, http.StatusCreated, contact)
}

func (h *CustomerHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	var compIDPtr *int64
	if compIDStr := r.URL.Query().Get("company_id"); compIDStr != "" {
		if cid, err := strconv.ParseInt(compIDStr, 10, 64); err == nil {
			compIDPtr = &cid
		}
	}
	contacts, err := h.svc.ListContacts(r.Context(), actor, compIDPtr)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, contacts)
}

func (h *CustomerHandler) GetContact(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的聯絡人識別碼")
		return
	}
	contact, err := h.svc.GetContact(r.Context(), actor, id)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", "找不到指定的客戶聯絡人")
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, contact)
}

func (h *CustomerHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的聯絡人識別碼")
		return
	}
	var req service.UpdateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "無效的 JSON 內容")
		return
	}
	contact, err := h.svc.UpdateContact(r.Context(), actor, id, req)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", "找不到指定的客戶聯絡人")
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, contact)
}
