package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/NoierBB/englishSchool/internal/dto"
	"github.com/NoierBB/englishSchool/internal/models"
	"github.com/NoierBB/englishSchool/internal/services"
	"github.com/go-chi/chi/v5"
)

type HandlerFacade struct {
	service     services.StudentService
	userService services.UserService
	authService *services.AuthService
}

func NewHandlerFacade(service services.StudentService, authService services.AuthService) *HandlerFacade {
	return &HandlerFacade{service: service, authService: &authService}
}

func (h *HandlerFacade) RegisterStudent(w http.ResponseWriter, r *http.Request) {
	log.Println("=== RegisterStudent handler called ===")

	var req dto.RegisterStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Request: Email=%s, Name=%s, Age=%d, Level=%s", req.Email, req.Name, req.Age, req.Level)

	if h.authService == nil {
		log.Println("ERROR: authService is nil in HandlerFacade")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID, studentID, err := h.authService.RegisterStudent(r.Context(), req)
	if err != nil {
		log.Printf("RegisterStudent error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Success: userID=%d, studentID=%d", userID, studentID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":    userID,
		"student_id": studentID,
		"message":    "Student registered successfully",
	})
}
func (hp *HandlerFacade) GetStudents(w http.ResponseWriter, r *http.Request) {
	students, err := hp.service.GetStudents(r.Context())
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

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	student, err := hp.service.GetStudentById(r.Context(), id)
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

	err := hp.service.UpdateStudent(r.Context(), s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hp *HandlerFacade) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = hp.service.DeleteStudent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
