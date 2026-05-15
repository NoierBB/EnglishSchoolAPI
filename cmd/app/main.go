package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NoierBB/englishSchool/internal/config"
	"github.com/NoierBB/englishSchool/internal/handlers"
	"github.com/NoierBB/englishSchool/internal/repositories"
	"github.com/NoierBB/englishSchool/pkg/db"
)

func main() {
	cfg := config.Load()

	database, err := db.NewDBConnect(
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer database.Close()

	studentsRepo := repositories.NewStudentRepository(database)

	handler := handlers.NewHandlerFacade(studentsRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetStudents(w, r)
		case http.MethodPost:
			handler.CreateStudent(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/student", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetStudentById(w, r)
		case http.MethodPut:
			handler.UpdateStudent(w, r)
		case http.MethodDelete:
			handler.DeleteStudent(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("server started")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
