package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NoierBB/englishSchool/internal/models"
	"github.com/NoierBB/englishSchool/internal/services"
	"github.com/go-chi/chi/v5"
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

func (hp *GroupHandlerFacade) AddStudent(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid group id", 400)
		return
	}
	var req struct {
		StudentID int `json:"student_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	err = hp.service.AddStudent(r.Context(), groupId, req.StudentID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (hp *GroupHandlerFacade) GetStudentGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupId, _ := strconv.Atoi(idStr)

	students, err := hp.service.GetStudentGroup(r.Context(), groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(students)
}
