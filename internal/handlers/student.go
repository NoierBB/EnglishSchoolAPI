package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/NoierBB/englishSchool/internal/models"
)

type StudentsRepo interface {
	GetStudents(ctx context.Context) ([]models.Students, error)
}

type HandlerFacade struct {
	repo StudentsRepo
}

func NewHandlerFacade(r StudentsRepo) *HandlerFacade {
	return &HandlerFacade{repo: r}
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
