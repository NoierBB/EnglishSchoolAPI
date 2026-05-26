package router

import (
	"log"
	"net/http"
	"time"

	"github.com/NoierBB/englishSchool/internal/handlers"
	"github.com/NoierBB/englishSchool/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func PrintRoutes(r chi.Router) {
	chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s", method, route)
		return nil
	})
}

func NewRouter(
	studentHandler *handlers.HandlerFacade,
	userHandler *handlers.UserHandlerFacade,
	groupHandler *handlers.GroupHandlerFacade,
) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recover)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))
	r.Use(middleware.Logger)

	r.Use(middleware.RequestId)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Route("/students", func(r chi.Router) {
		r.Get("/", studentHandler.GetStudents)
		// r.Post("/register", studentHandler.CreateStudent)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", studentHandler.GetStudentById)
			r.Put("/", studentHandler.UpdateStudent)
			r.Delete("/", studentHandler.DeleteStudent)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.GetUsers)
		r.Post("/", userHandler.CreateUser)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", userHandler.GetUserById)
			r.Put("/", userHandler.UpdateUser)
			r.Delete("/", userHandler.DeleteUser)
		})
	})

	r.Route("/groups", func(r chi.Router) {
		r.Get("/", groupHandler.GetGroup)
		r.Post("/", groupHandler.CreateGroup)

		r.Route("/{id}", func(r chi.Router) {
			r.Post("/students", groupHandler.AddStudent)
			r.Get("/students", groupHandler.GetStudentGroup)
		})
	})
	r.Route("/auth", func(r chi.Router) {
		r.Route("/register", func(r chi.Router) {
			r.Post("/student", studentHandler.RegisterStudent)
		})
		r.Post("/login", userHandler.Login)
	})
	return r
}
