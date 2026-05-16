package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NoierBB/englishSchool/internal/models"
)

type UserRepo interface {
	CreateUser(ctx context.Context, u models.User) (int, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, u models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type UserHandlerFacade struct {
	repo UserRepo
}

func NewUserHandlerFacade(r UserRepo) *UserHandlerFacade {
	return &UserHandlerFacade{repo: r}
}

func (hp *UserHandlerFacade) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := hp.repo.CreateUser(r.Context(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (hp *UserHandlerFacade) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := hp.repo.GetUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (hp *UserHandlerFacade) GetUserById(w http.ResponseWriter, r *http.Request) {
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
	user, err := hp.repo.GetUserById(r.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}

func (hp *UserHandlerFacade) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := hp.repo.UpdateUser(r.Context(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hp *UserHandlerFacade) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	err = hp.repo.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
