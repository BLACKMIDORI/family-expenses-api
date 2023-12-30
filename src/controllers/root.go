package controllers

import (
	"net/http"
)

type Root struct{}

func init() {
	http.Handle("/", &Root{})
}

func (controller *Root) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" { // Check path here
		http.NotFound(responseWriter, request)
		return
	}
	http.Redirect(responseWriter, request, "/healthcheck", 307)
}
