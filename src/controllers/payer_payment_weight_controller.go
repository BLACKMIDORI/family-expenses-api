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

type PayerPaymentWeightController struct {
	basePath string
	user     core.AuthenticatedUser
}

func init() {
	http.Handle("/v1/payer-payment-weights/", &PayerPaymentWeightController{"/v1/payer-payment-weights/", core.AuthenticatedUser{}})
}

func (controller *PayerPaymentWeightController) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
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

// POST /payer-payment-weights/
func (controller *PayerPaymentWeightController) create(responseWriter http.ResponseWriter, request *http.Request) {
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
	if value, ok := body["weight"]; !ok || value.(float64) == 0 {
		validationErrors["weight"] = []string{"'weight' is required"}
	}
	{
		value, ok := body["payer"]
		if !ok {
			validationErrors["payer.id"] = []string{"'payer.id' is required"}
		} else if value2, ok2 := value.(map[string]any)["id"]; !ok2 || value2.(string) == "" {
			validationErrors["payer.id"] = []string{"'payer.id' is required"}
		}
	}
	{
		value, ok := body["chargeAssociation"]
		if !ok {
			validationErrors["chargeAssociation.id"] = []string{"'chargeAssociation.id' is required"}
		} else if value2, ok2 := value.(map[string]any)["id"]; !ok2 || value2.(string) == "" {
			validationErrors["chargeAssociation.id"] = []string{"'chargeAssociation.id' is required"}
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
	payerPaymentWeightService := services.CreatePayerPaymentWeightService(
		repositories.CreatePayerPaymentWeightRepository(transaction),
		controller.user,
	)
	payerPaymentWeight, httpErr := payerPaymentWeightService.Create(ctx, models.PayerPaymentWeight{
		Weight: body["weight"].(float64),
		Payer: models.ForeignKeyHolder{
			Id: body["payer"].(map[string]any)["id"].(string),
		},
		ChargeAssociation: models.ForeignKeyHolder{
			Id: body["chargeAssociation"].(map[string]any)["id"].(string),
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
	_ = json.NewEncoder(responseWriter).Encode(payerPaymentWeight)
}

// GET /payer-payment-weights/
func (controller *PayerPaymentWeightController) readPage(responseWriter http.ResponseWriter, request *http.Request) {
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
	if value := filters.Get("chargeAssociation.id"); value == "" {
		validationErrors["filter"] = []string{"query 'filter=chargeAssociation.id__*' is required"}
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
	payerPaymentWeightService := services.CreatePayerPaymentWeightService(
		repositories.CreatePayerPaymentWeightRepository(transaction),
		controller.user,
	)
	payerPaymentWeights, httpErr := payerPaymentWeightService.GetAll(ctx, filters.Get("chargeAssociation.id"))
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}
	_ = json.NewEncoder(responseWriter).Encode(core.PagedList[models.PayerPaymentWeight]{
		Size:    999,
		From:    0,
		Results: payerPaymentWeights,
	})
}

// GET /payer-payment-weights/{id}
func (controller *PayerPaymentWeightController) read(responseWriter http.ResponseWriter, request *http.Request) {
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
	payerPaymentWeightService := services.CreatePayerPaymentWeightService(
		repositories.CreatePayerPaymentWeightRepository(transaction),
		controller.user,
	)
	payerPaymentWeight, httpErr := payerPaymentWeightService.GetOne(ctx, params["id"])
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}
	_ = json.NewEncoder(responseWriter).Encode(payerPaymentWeight)
}

// PUT /payer-payment-weights/{id}
func (controller *PayerPaymentWeightController) update(responseWriter http.ResponseWriter, request *http.Request) {
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
	if prop, ok := body["weight"]; !ok || prop.(float64) == 0 {
		validationErrors["weight"] = []string{"'weight' is required"}
	}
	{
		value, ok := body["payer"]
		if !ok {
			validationErrors["payer.id"] = []string{"'payer.id' is required"}
		} else if value2, ok2 := value.(map[string]any)["id"]; !ok2 || value2.(string) == "" {
			validationErrors["payer.id"] = []string{"'payer.id' is required"}
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
	payerPaymentWeightService := services.CreatePayerPaymentWeightService(
		repositories.CreatePayerPaymentWeightRepository(transaction),
		controller.user,
	)
	payerPaymentWeight, httpErr := payerPaymentWeightService.GetOne(ctx, params["id"])
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}

	// Update fields
	payerPaymentWeight.Weight = body["weight"].(float64)
	payerPaymentWeight.Payer.Id = body["payer"].(map[string]any)["id"].(string)

	payerPaymentWeight, httpErr = payerPaymentWeightService.Update(ctx, payerPaymentWeight)
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
	_ = json.NewEncoder(responseWriter).Encode(payerPaymentWeight)
}

// DELETE /payer-payment-weights/{id}
func (controller *PayerPaymentWeightController) delete(responseWriter http.ResponseWriter, request *http.Request) {
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
	payerPaymentWeightService := services.CreatePayerPaymentWeightService(
		repositories.CreatePayerPaymentWeightRepository(transaction),
		controller.user,
	)
	httpErr := payerPaymentWeightService.Delete(ctx, params["id"])
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
