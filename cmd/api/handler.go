package main

import (
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

type HandlerInterface interface {
	NewHandler() *Handler
	CheckUser(http.ResponseWriter, *http.Request)
	HandleLog(string, bool)
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
	c.loggerInterface.HandleLog(input.Email, match)
}

func getEnv(key, default_value string) string {
	value := os.Getenv(key)
	if len(value) != 0 {
		return value
	}
	return default_value
}
