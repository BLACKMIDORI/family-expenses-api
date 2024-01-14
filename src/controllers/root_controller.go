package controllers

import (
	"net/http"
)

type RootController struct{}

func init() {
	http.Handle("/", &RootController{})
}

func (controller *RootController) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" { // Check path here
		http.NotFound(responseWriter, request)
		return
	}
	http.Redirect(responseWriter, request, "/healthcheck", 307)
}
