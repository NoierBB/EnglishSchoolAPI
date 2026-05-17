package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NoierBB/englishSchool/internal/models"
	"github.com/NoierBB/englishSchool/internal/services"
)

type GroupHandlerFacade struct {
	service services.GroupService
}

func NewGroupHandlerFacade(service services.GroupService) *GroupHandlerFacade {
	return &GroupHandlerFacade{service: service}
}

func (hp *GroupHandlerFacade) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var g models.Group

	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "invalid request body", http.StatusInternalServerError)
		return
	}

	id, err := hp.service.CreateGroup(r.Context(), g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (hp *GroupHandlerFacade) GetGroup(w http.ResponseWriter, r *http.Request) {
	groups, err := hp.service.GetGroup(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
