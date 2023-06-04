package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (*Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	const maxBytes = 1048576
	reader := http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&data)
	fmt.Println(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("should have only one json block")
	}
	return nil
}

func (*Config) writeJson(w http.ResponseWriter, status int, data any, headers ...http.Header) {
	jData, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Write(jData)
}

func (c *Config) errorJson(w http.ResponseWriter, message string, statusCode ...int) {
	status := http.StatusBadRequest
	if len(statusCode) > 0 {
		status = statusCode[0]
	}
	response := jsonResponse{
		Error:   true,
		Message: message,
	}
	c.writeJson(w, status, response)
}
