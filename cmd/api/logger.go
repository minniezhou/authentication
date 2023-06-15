package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/minniezhou/jsonToolBox"
)

type LoggerInterface interface {
	HandleLog(string, bool)
}

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (*Logger) HandleLog(email string, authed bool) {
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

	logging_host := jsonToolBox.GetEnv("LOGGING_SERVICE", "localhost")
	_, err = http.Post("http://"+logging_host+":4321/log", "application/json", requesBody)
	if err != nil {
		fmt.Println("http post to log failed")
	}
}
