package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/minniezhou/jsonToolBox"
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
		jsonToolBox.WriteJson(w, http.StatusAccepted, jsonResponse{Error: false, Message: "Ping Authortication Service"})
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
	err := jsonToolBox.ReadJson(w, r, &input)
	if err != nil {
		jsonToolBox.ErrorJson(w, "Wrong username or password!")
		return
	}
	m := NewModel(c.DB)
	user, err := m.user.GetInfoByEmail(input.Email)
	if err != nil || user == nil {
		jsonToolBox.ErrorJson(w, "Wrong username or password!")
		return
	}
	match := user.MatchPassword(input.Password)

	if match {
		response := jsonToolBox.JsonResponse{Error: false, Message: "User Authorized!"}
		jsonToolBox.WriteJson(w, http.StatusAccepted, response)
	} else {
		jsonToolBox.ErrorJson(w, "Wrong username or password!")
	}
	c.handleLog(input.Email, match)
}

func (c *Config) handleLog(email string, authed bool) {
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
		return
	}
	requesBody := bytes.NewBuffer(logPayload)

	logging_host := getEnv("LOGGING_SERVICE", "localhost")
	_, err = http.Post("http://"+logging_host+":4321/log", "application/json", requesBody)
	if err != nil {
		fmt.Println("http post to log failed")
	}
}

func getEnv(key, default_value string) string {
	value := os.Getenv(key)
	if len(value) != 0 {
		return value
	}
	return default_value
}
