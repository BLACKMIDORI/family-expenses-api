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

type ChargeAssociationController struct {
	basePath string
	user     core.AuthenticatedUser
}

func init() {
	http.Handle("/v1/charge-associations/", &ChargeAssociationController{"/v1/charge-associations/", core.AuthenticatedUser{}})
}

func (controller *ChargeAssociationController) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
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

// POST /charge-associations/
func (controller *ChargeAssociationController) create(responseWriter http.ResponseWriter, request *http.Request) {
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
		value, ok := body["expense"]
		if !ok {
			validationErrors["expense.id"] = []string{"'expense.id' is required"}
		} else if value2, ok2 := value.(map[string]any)["id"]; !ok2 || value2.(string) == "" {
			validationErrors["expense.id"] = []string{"'expense.id' is required"}
		}
	}
	{
		value, ok := body["actualPayer"]
		if !ok {
			validationErrors["actualPayer.id"] = []string{"'actualPayer.id' is required"}
		} else if value2, ok2 := value.(map[string]any)["id"]; !ok2 || value2.(string) == "" {
			validationErrors["actualPayer.id"] = []string{"'actualPayer.id' is required"}
		}
	}
	{
		value, ok := body["chargesModel"]
		if !ok {
			validationErrors["chargesModel.id"] = []string{"'chargesModel.id' is required"}
		} else if value2, ok2 := value.(map[string]any)["id"]; !ok2 || value2.(string) == "" {
			validationErrors["chargesModel.id"] = []string{"'chargesModel.id' is required"}
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
	chargeAssociationService := services.CreateChargeAssociationService(
		repositories.CreateChargeAssociationRepository(transaction),
		controller.user,
	)
	chargeAssociation, httpErr := chargeAssociationService.Create(ctx, models.ChargeAssociation{
		Name: body["name"].(string),
		Expense: models.ForeignKeyHolder{
			Id: body["expense"].(map[string]any)["id"].(string),
		},
		ActualPayer: models.ForeignKeyHolder{
			Id: body["actualPayer"].(map[string]any)["id"].(string),
		},
		ChargesModel: models.ForeignKeyHolder{
			Id: body["chargesModel"].(map[string]any)["id"].(string),
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
	_ = json.NewEncoder(responseWriter).Encode(chargeAssociation)
}

// GET /charge-associations/
func (controller *ChargeAssociationController) readPage(responseWriter http.ResponseWriter, request *http.Request) {
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
	if value := filters.Get("chargesModel.id"); value == "" {
		validationErrors["filter"] = []string{"query 'filter=chargesModel.id__*' is required"}
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
	chargeAssociationService := services.CreateChargeAssociationService(
		repositories.CreateChargeAssociationRepository(transaction),
		controller.user,
	)
	chargeAssociations, httpErr := chargeAssociationService.GetAll(ctx, filters.Get("chargesModel.id"))
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}
	_ = json.NewEncoder(responseWriter).Encode(core.PagedList[models.ChargeAssociation]{
		Size:    999,
		From:    0,
		Results: chargeAssociations,
	})
}

// GET /charge-associations/{id}
func (controller *ChargeAssociationController) read(responseWriter http.ResponseWriter, request *http.Request) {
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
	chargeAssociationService := services.CreateChargeAssociationService(
		repositories.CreateChargeAssociationRepository(transaction),
		controller.user,
	)
	chargeAssociation, httpErr := chargeAssociationService.GetOne(ctx, params["id"])
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}
	_ = json.NewEncoder(responseWriter).Encode(chargeAssociation)
}

// PUT /charge-associations/{id}
func (controller *ChargeAssociationController) update(responseWriter http.ResponseWriter, request *http.Request) {
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
	{
		value, ok := body["expense"]
		if !ok {
			validationErrors["expense.id"] = []string{"'expense.id' is required"}
		} else if value2, ok2 := value.(map[string]any)["id"]; !ok2 || value2.(string) == "" {
			validationErrors["expense.id"] = []string{"'expense.id' is required"}
		}
	}
	{
		value, ok := body["actualPayer"]
		if !ok {
			validationErrors["actualPayer.id"] = []string{"'actualPayer.id' is required"}
		} else if value2, ok2 := value.(map[string]any)["id"]; !ok2 || value2.(string) == "" {
			validationErrors["actualPayer.id"] = []string{"'actualPayer.id' is required"}
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
	chargeAssociationService := services.CreateChargeAssociationService(
		repositories.CreateChargeAssociationRepository(transaction),
		controller.user,
	)
	chargeAssociation, httpErr := chargeAssociationService.GetOne(ctx, params["id"])
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}

	// Update fields
	chargeAssociation.Name = body["name"].(string)
	chargeAssociation.Expense.Id = body["expense"].(map[string]any)["id"].(string)
	chargeAssociation.ActualPayer.Id = body["actualPayer"].(map[string]any)["id"].(string)

	chargeAssociation, httpErr = chargeAssociationService.Update(ctx, chargeAssociation)
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
	_ = json.NewEncoder(responseWriter).Encode(chargeAssociation)
}

// DELETE /charge-associations/{id}
func (controller *ChargeAssociationController) delete(responseWriter http.ResponseWriter, request *http.Request) {
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
	chargeAssociationService := services.CreateChargeAssociationService(
		repositories.CreateChargeAssociationRepository(transaction),
		controller.user,
	)
	httpErr := chargeAssociationService.Delete(ctx, params["id"])
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
