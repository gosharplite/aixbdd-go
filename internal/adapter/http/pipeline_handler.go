package http

import (
	"errors"
	"net/http"

	"crm/internal/domain"
	"crm/internal/service"
)

type PipelineHandler struct {
	svc *service.PipelineService
}

func NewPipelineHandler(svc *service.PipelineService) *PipelineHandler {
	return &PipelineHandler{svc: svc}
}

func (h *PipelineHandler) GetPipeline(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	summary, err := h.svc.GetPipelineSummary(r.Context(), actor)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, summary)
}
