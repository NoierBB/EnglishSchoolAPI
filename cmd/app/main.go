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

	studentsRepo := repositories.NewStudentRepository(database.DB)
	userRepo := repositories.NewUserRepository(database.DB)
	groupRepo := repositories.NewGropRepo(database.DB)

	jwtSecret := cfg.JWTSecret
	userService := services.NewAuthService(userRepo, jwtSecret)

	handler := handlers.NewHandlerFacade(studentsRepo)
	handlerUser := handlers.NewUserHandlerFacade(userService)
	handlerGroup := handlers.NewGroupHandlerFacade(groupRepo)

	r := router.NewRouter(handler, handlerUser, handlerGroup)
	log.Println("Available routes:")
	router.PrintRoutes(r)
	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("server started on port", cfg.HTTPPort)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}

}
