package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *PaymentHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/payment", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.GetByPerson)
		r.Patch("/", handler.Update)
		r.Delete("/", handler.Delete)
	})

	return r
}
