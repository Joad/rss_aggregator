package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func routes() *chi.Mux {
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

	router.Mount("/v1", api)
	return router
}
