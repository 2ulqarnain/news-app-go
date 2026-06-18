package routes

import (
	"GoNewsScrapper/internals/handlers"

	"github.com/go-chi/chi/v5"
)

func NewsRouter(h *handlers.News) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/news", func(r chi.Router) {
		r.Get("/", h.GetAllNews)
	})
	return r
}
