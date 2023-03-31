package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	router *chi.Mux
}

func (*Config) NewHandler() *Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http//*", "https://"},
		AllowCredentials: false,
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("ping"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) { fmt.Println("Ping authentification service") })
	return &Handler{
		router: r,
	}
}
