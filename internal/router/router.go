package router

import (
	"net/http"

	"github.com/NoierBB/englishSchool/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	studentHandler *handlers.HandlerFacade,
	userHandler *handlers.UserHandlerFacade,
	groupHandler *handlers.GroupHandlerFacade,
) http.Handler {
	r := chi.NewRouter()

	r.Route("/students", func(r chi.Router) {
		r.Get("/", studentHandler.GetStudents)
		r.Post("/", studentHandler.CreateStudent)

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
	})
	return r
}
