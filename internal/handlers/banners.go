package handlers

import (
	"app/internal/models"
	"net/http"
)

type banner struct {
	srv models.BannerService
}

func NewBanner(srv models.BannerService) *banner {
	return &banner{srv: srv}
}

func (h *banner) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.srv.Ping(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
