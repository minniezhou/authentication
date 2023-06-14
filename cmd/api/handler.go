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
		err := jsonToolBox.WriteJson(w, http.StatusAccepted, jsonToolBox.JsonResponse{Error: false, Message: "Ping Authortication Service"})
		if err != nil {
			fmt.Println("writing json error")
		}
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
		err = jsonToolBox.ErrorJson(w, "Wrong username or password!")
		if err != nil {
			fmt.Println("writing json error")
		}
		return
	}
	user, err := c.userInterface.GetInfoByEmail(input.Email)
	if err != nil || user == nil {
		err = jsonToolBox.ErrorJson(w, "Wrong username or password!")
		if err != nil {
			fmt.Println("writing json error")
		}
		return
	}
	match := c.userInterface.MatchPassword(input.Password)

	if match {
		response := jsonToolBox.JsonResponse{Error: false, Message: "User Authorized!"}
		err = jsonToolBox.WriteJson(w, http.StatusAccepted, response)
		if err != nil {
			fmt.Println("writing json error")
		}
	} else {
		err = jsonToolBox.ErrorJson(w, "Wrong username or password!")
		if err != nil {
			fmt.Println("writing json error")
		}
	}
	//c.handleLog(input.Email, match)
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
