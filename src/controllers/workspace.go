package controllers

import (
	"context"
	"encoding/json"
	"family-expenses-api/core"
	"family-expenses-api/models"
	"family-expenses-api/repositories"
	"family-expenses-api/services"
	"log"
	"net/http"
)

type Workspace struct {
	basePath string
	user     core.AuthenticatedUser
}

func init() {
	http.Handle("/v1/workspaces/", &Workspace{"/v1/workspaces/", core.AuthenticatedUser{}})
}

func (controller *Workspace) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	defer sendInternalServerErrorOnPanic(responseWriter)
	userChan := parseUserAsync(request)

	basePath := controller.basePath
	switch request.Method {
	case http.MethodPost:
		if request.URL.Path == basePath {
			controller.user = <-userChan
			controller.create(responseWriter, request)
			return
		}
	case http.MethodGet:
		if request.URL.Path == basePath {
			controller.user = <-userChan
			controller.readPage(responseWriter, request)
			return
		}
		if match(basePath+"[\\w-]+", request.URL.Path) {
			controller.user = <-userChan
			controller.read(responseWriter, request)
			return
		}
	case http.MethodPut:
		if match(basePath+"[\\w-]+", request.URL.Path) {
			controller.user = <-userChan
			controller.update(responseWriter, request)
			return
		}
	case http.MethodDelete:
		if match(basePath+"[\\w-]+", request.URL.Path) {
			controller.user = <-userChan
			controller.delete(responseWriter, request)
			return
		}
	}
	http.NotFound(responseWriter, request)
}

// POST {basePath}
func (controller *Workspace) create(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling", request.Method, request.URL)
	// * Auth
	if controller.user.Id == "" {
		replyAsJson(responseWriter, 401, map[string]any{
			"error": "Unauthenticated",
		})
		return
	}

	// * Parse
	body := bodyJson(request)

	// * Validate
	validationErrors := make(map[string][]string)
	if prop, ok := body["name"]; !ok || prop.(string) == "" {
		validationErrors["name"] = []string{"'name' is required"}
	}
	if len(validationErrors) != 0 {
		replyAsJson(responseWriter, 400, map[string]any{
			"errors": validationErrors,
		})
		return
	}

	// Initialize a new DB transaction
	ctx := context.TODO()
	transaction, err := core.BeginTransaction(ctx)
	if err != nil {
		log.Print("failed to connect database", err)
		http.Error(responseWriter, "Could not connect to database", 503)
		return
	}

	// Construct Dependencies
	workspaceService := services.CreateWorkspaceService(
		repositories.CreateWorkspaceRepository(transaction),
		controller.user,
	)
	workspace, httpErr := workspaceService.Create(ctx, models.Workspace{
		Name: body["name"].(string),
	})
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}
	err = transaction.Commit(ctx)
	if err != nil {
		log.Print("failed to save to database", err)
		replyAsJson(responseWriter, 500, map[string]any{
			"error": "Error saving changes",
		})
		return
	}
	_ = json.NewEncoder(responseWriter).Encode(workspace)
}

// GET /workspace
func (controller *Workspace) readPage(responseWriter http.ResponseWriter, request *http.Request) {
	// Validate

	params := routeParams(request, controller.basePath+"{id}")
	_ = params["id"]
}

// GET /workspace/{id}
func (controller *Workspace) read(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling", request.Method, request.URL)
	// Validate
	params := routeParams(request, controller.basePath+"{id}")
	_ = params["id"]
}

// PUT /workspace/{id}
func (controller *Workspace) update(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling", request.Method, request.URL)
	// Validate
	params := routeParams(request, controller.basePath+"{id}")
	_ = params["id"]
}

// DELETE /workspace/{id}
func (controller *Workspace) delete(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling", request.Method, request.URL)
	// Validate
	params := routeParams(request, controller.basePath+"{id}")
	_ = params["id"]
}
