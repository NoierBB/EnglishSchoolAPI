package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NoierBB/englishSchool/internal/config"
	"github.com/NoierBB/englishSchool/internal/handlers"
	"github.com/NoierBB/englishSchool/internal/repositories"
	"github.com/NoierBB/englishSchool/internal/router"
	"github.com/NoierBB/englishSchool/internal/services"
	"github.com/NoierBB/englishSchool/pkg/db"
)

func main() {
	log.Println("Starting server...")

	cfg := config.Load()
	log.Println("Config loaded")

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
	log.Println("Database connected")

	studentsRepo := repositories.NewStudentRepository(database.DB)
	userRepo := repositories.NewUserRepository(database.DB)

	groupRepo := repositories.NewGropRepo(database.DB)

	log.Println("Repositories created")

	jwtSecret := cfg.JWTSecret
	log.Printf("JWT secret length: %d", len(jwtSecret))

	authService := services.NewAuthService(userRepo, studentsRepo, database.DB, jwtSecret)

	handler := handlers.NewHandlerFacade(studentsRepo, *authService)
	handlerUser := handlers.NewUserHandlerFacade(authService)
	handlerGroup := handlers.NewGroupHandlerFacade(groupRepo)

	log.Println("Handlers created")

	r := router.NewRouter(handler, handlerUser, handlerGroup)
	log.Println("Router created")

	log.Println("Available routes:")
	router.PrintRoutes(r)

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server started on port %s", cfg.HTTPPort)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
