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

type ChargesModelController struct {
	basePath string
	user     core.AuthenticatedUser
}

func init() {
	http.Handle("/v1/charges-models/", &ChargesModelController{"/v1/charges-models/", core.AuthenticatedUser{}})
}

func (controller *ChargesModelController) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
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

// POST /charges-models/
func (controller *ChargesModelController) create(responseWriter http.ResponseWriter, request *http.Request) {
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
	if value, ok := body["name"]; !ok || value.(string) == "" {
		validationErrors["name"] = []string{"'name' is required"}
	}
	{
		value, ok := body["workspace"]
		if !ok {
			validationErrors["workspace.id"] = []string{"'workspace.id' is required"}
		} else if value2, ok2 := value.(map[string]any)["id"]; !ok2 || value2.(string) == "" {
			validationErrors["workspace.id"] = []string{"'workspace.id' is required"}
		}
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
	chargesModelService := services.CreateChargesModelService(
		repositories.CreateChargesModelRepository(transaction),
		controller.user,
	)
	chargesModel, httpErr := chargesModelService.Create(ctx, models.ChargesModel{
		Name: body["name"].(string),
		Workspace: models.ForeignKeyHolder{
			Id: body["workspace"].(map[string]any)["id"].(string),
		},
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
	_ = json.NewEncoder(responseWriter).Encode(chargesModel)
}

// GET /charges-models/
func (controller *ChargesModelController) readPage(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling", request.Method, request.URL)
	// * Auth
	if controller.user.Id == "" {
		replyAsJson(responseWriter, 401, map[string]any{
			"error": "Unauthenticated",
		})
		return
	}

	// * Parse
	filters := filtersFromQuery(request)

	// * Validate
	validationErrors := make(map[string][]string)
	if value := filters.Get("workspace.id"); value == "" {
		validationErrors["filter"] = []string{"query 'filter=workspace.id__*' is required"}
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
	chargesModelService := services.CreateChargesModelService(
		repositories.CreateChargesModelRepository(transaction),
		controller.user,
	)
	chargesModels, httpErr := chargesModelService.GetAll(ctx, filters.Get("workspace.id"))
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}
	_ = json.NewEncoder(responseWriter).Encode(core.PagedList[models.ChargesModel]{
		Size:    999,
		From:    0,
		Results: chargesModels,
	})
}

// GET /charges-models/{id}
func (controller *ChargesModelController) read(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling", request.Method, request.URL)
	// * Auth
	if controller.user.Id == "" {
		replyAsJson(responseWriter, 401, map[string]any{
			"error": "Unauthenticated",
		})
		return
	}

	// * Parse
	params := routeParams(request, controller.basePath+"{id}")

	// * Validate
	validationErrors := make(map[string][]string)
	if prop, ok := params["id"]; !ok || prop == "" {
		validationErrors["id"] = []string{"'id' is required"}
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
	chargesModelService := services.CreateChargesModelService(
		repositories.CreateChargesModelRepository(transaction),
		controller.user,
	)
	chargesModel, httpErr := chargesModelService.GetOne(ctx, params["id"])
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}
	_ = json.NewEncoder(responseWriter).Encode(chargesModel)
}

// PUT /charges-models/{id}
func (controller *ChargesModelController) update(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling", request.Method, request.URL)
	// * Auth
	if controller.user.Id == "" {
		replyAsJson(responseWriter, 401, map[string]any{
			"error": "Unauthenticated",
		})
		return
	}

	// * Parse
	params := routeParams(request, controller.basePath+"{id}")
	body := bodyJson(request)

	// * Validate
	validationErrors := make(map[string][]string)
	if prop, ok := params["id"]; !ok || prop == "" {
		validationErrors["id"] = []string{"'id' is required"}
	}
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
	chargesModelService := services.CreateChargesModelService(
		repositories.CreateChargesModelRepository(transaction),
		controller.user,
	)
	chargesModel, httpErr := chargesModelService.GetOne(ctx, params["id"])
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}

	// Update fields
	chargesModel.Name = body["name"].(string)

	chargesModel, httpErr = chargesModelService.Update(ctx, chargesModel)
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
	_ = json.NewEncoder(responseWriter).Encode(chargesModel)
}

// DELETE /charges-models/{id}
func (controller *ChargesModelController) delete(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling", request.Method, request.URL)
	// * Auth
	if controller.user.Id == "" {
		replyAsJson(responseWriter, 401, map[string]any{
			"error": "Unauthenticated",
		})
		return
	}

	// * Parse
	params := routeParams(request, controller.basePath+"{id}")

	// * Validate
	validationErrors := make(map[string][]string)
	if prop, ok := params["id"]; !ok || prop == "" {
		validationErrors["id"] = []string{"'id' is required"}
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
	chargesModelService := services.CreateChargesModelService(
		repositories.CreateChargesModelRepository(transaction),
		controller.user,
	)
	httpErr := chargesModelService.Delete(ctx, params["id"])
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
	_ = json.NewEncoder(responseWriter).Encode(map[string]any{
		"ok": true,
	})
}
