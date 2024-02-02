package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (cfg *apiConfig) routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	api := chi.NewRouter()
	api.Get("/readiness", readinessHandler)
	api.Get("/err", errHandler)

	api.Post("/users", cfg.createUserHandler)
	api.Get("/users", cfg.middlewareAuth(cfg.getUserHandler))

	api.Post("/feeds", cfg.middlewareAuth(cfg.createFeedHandler))

	router.Mount("/v1", api)
	return router
}
