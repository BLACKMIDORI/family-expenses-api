package controllers

import (
	"encoding/json"
	"net/http"
)

type HealthCheck struct{}

func init() {
	http.Handle("/healthcheck", &HealthCheck{})
}
func (controller *HealthCheck) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		_ = json.NewEncoder(responseWriter).Encode(map[string]any{
			"healthy": true,
			"message": "Healthy!",
		})
	default:
		http.NotFound(responseWriter, request)
	}
}
