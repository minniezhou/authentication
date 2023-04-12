package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	router *chi.Mux
}

func (c *Config) NewHandler() *Handler {
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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		c.writeJson(w, http.StatusAccepted, jsonResponse{Error: false, Message: "Ping Authortication Service"})
	})
	r.Post("/auth", c.CheckUser)
	return &Handler{
		router: r,
	}
}

func (c *Config) CheckUser(w http.ResponseWriter, r *http.Request) {
	type userType struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input userType
	err := c.readJson(w, r, &input)
	if err != nil {
		c.errorJson(w, "Wrong username or password!")
		return
	}
	m := NewModel(c.DB)
	user, err := m.user.GetInfoByEmail(input.Email)
	if err != nil || user == nil {
		c.errorJson(w, "Wrong username or password!")
		return
	}
	match := user.MatchPassword(input.Password)

	if match {
		response := jsonResponse{Error: false, Message: "User Authorized!"}
		c.writeJson(w, http.StatusAccepted, response)
	} else {
		c.errorJson(w, "Wrong username or password!")
	}
	c.handleLog(input.Email, match)
}

func (c *Config) handleLog(email string, authed bool) error {
	type logType struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	message := email + "is not authorized!"
	if authed {
		message = email + "is authorized!"
	}
	log := logType{
		Name:    "log from authentitcation service",
		Message: message,
	}

	logPayload, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return err
	}
	requesBody := bytes.NewBuffer(logPayload)

	_, err = http.Post("http://localhost:4321/log", "application/json", requesBody)
	if err != nil {
		return err
	}
	return nil
}
