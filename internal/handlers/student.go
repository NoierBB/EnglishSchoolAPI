package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NoierBB/englishSchool/internal/models"
)

type StudentsRepo interface {
	CreateStudent(ctx context.Context, s models.Students) (int, error)
	GetStudents(ctx context.Context) ([]models.Students, error)
	GetStudentById(ctx context.Context, id int) (*models.Students, error)
	UpdateStudent(ctx context.Context, s models.Students) error
	DeleteStudent(ctx context.Context, id int) error
}

type HandlerFacade struct {
	repo StudentsRepo
}

func NewHandlerFacade(r StudentsRepo) *HandlerFacade {
	return &HandlerFacade{repo: r}
}

func (hp *HandlerFacade) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var s models.Students

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := hp.repo.CreateStudent(r.Context(), s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (hp *HandlerFacade) GetStudents(w http.ResponseWriter, r *http.Request) {
	students, err := hp.repo.GetStudents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(students); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (hp *HandlerFacade) GetStudentById(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	student, err := hp.repo.GetStudentById(r.Context(), id)
	if err != nil {
		if err.Error() == "students not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(student); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}

func (hp *HandlerFacade) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	var s models.Students

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := hp.repo.UpdateStudent(r.Context(), s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hp *HandlerFacade) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = hp.repo.DeleteStudent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
