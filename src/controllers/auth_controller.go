package controllers

import (
	"context"
	"family-expenses-api/core"
	"family-expenses-api/repositories"
	"family-expenses-api/services"
	"log"
	"net/http"
	"time"
)

type AuthController struct {
	basePath string
}

func init() {
	http.Handle("/v1/auth/", &AuthController{"/v1/auth/"})
}

func (controller *AuthController) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	defer sendInternalServerErrorOnPanic(responseWriter)

	switch request.Method {
	case http.MethodPost:
		if request.URL.Path == controller.basePath+"tokensignin" {
			controller.tokenSignIn(responseWriter, request)
			return
		}
		if request.URL.Path == controller.basePath+"renew" {
			controller.renew(responseWriter, request)
			return
		}
	}
	http.NotFound(responseWriter, request)
}

// POST {basePath}/tokensignin
func (_ *AuthController) tokenSignIn(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling", request.Method, request.URL)
	// * Parse
	body := bodyJson(request)

	// * Validate
	validationErrors := make(map[string][]string)
	if prop, ok := body["idToken"]; !ok || prop.(string) == "" {
		validationErrors["idToken"] = []string{"'idToken' is required"}
	}
	if prop, ok := body["clientId"]; !ok || prop.(string) == "" {
		validationErrors["clientId"] = []string{"'clientId' is required"}
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
	tokenService := services.CreateTokenService(
		repositories.CreatePersistentGrantRepository(transaction),
	)
	appUserService := services.CreateAppUserService(
		repositories.CreateAppUserRepository(transaction),
		repositories.CreateAppUserLoginRepository(transaction),
	)

	// Start
	if !tokenService.IsClientIdValid(body["clientId"].(string)) {
		replyAsJson(responseWriter, 400, map[string]any{
			"errors": map[string][]string{
				"clientId": {"Invalid clientId"},
			},
		})
		return
	}

	verifiedIdentity, httpErr := tokenService.VerifyIdToken(body["idToken"].(string))
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}

	appUser, httpErr := appUserService.GetAppUserByLogin(ctx, verifiedIdentity.IdentityProvider, verifiedIdentity.Key)

	if httpErr != nil && httpErr.StatusCode() != 404 {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}

	if appUser.Id == "" {
		appUser, httpErr = appUserService.CreateWithLogin(ctx, verifiedIdentity.IdentityProvider, verifiedIdentity.Key)
	}

	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}
	now := time.Now()

	refreshToken, httpErr := tokenService.CreateRefreshToken(ctx, now, body["clientId"].(string), appUser.Id, "")
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}

	signedJwt, expirationDateTime, httpErr := tokenService.CreateAccessToken(now, body["clientId"].(string), appUser.Id)
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}

	err = transaction.Commit(ctx)
	if err != nil {
		log.Print("failed to save to database", err)
		http.Error(responseWriter, "Error saving changes", 500)
		return
	}
	replyAsJson(responseWriter, 200, map[string]any{
		"appUser": map[string]string{
			"id": appUser.Id,
		},
		"accessTokenExpirationDateTime": expirationDateTime.Format(time.RFC3339),
		"accessToken":                   signedJwt,
		"refreshToken":                  refreshToken,
	})
}

// POST {basePath}/renew
func (_ *AuthController) renew(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling", request.Method, request.URL)
	// * Parse
	body := bodyJson(request)

	// * Validate
	validationErrors := make(map[string][]string)
	if prop, ok := body["refreshToken"]; !ok || prop.(string) == "" {
		validationErrors["refreshToken"] = []string{"'refreshToken' is required"}
	}
	if len(body["refreshToken"].(string)) != 1280 {
		validationErrors["refreshToken"] = []string{"Invalid refresh token format"}
	}
	if prop, ok := body["clientId"]; !ok || prop.(string) == "" {
		validationErrors["clientId"] = []string{"'clientId' is required"}
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
	tokenService := services.CreateTokenService(
		repositories.CreatePersistentGrantRepository(transaction),
	)

	// Start
	now := time.Now()
	persistentGrant, httpErr := tokenService.ConsumeRefreshToken(ctx, body["refreshToken"].(string), now)
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}

	if persistentGrant.Id == "" {
		replyAsJson(responseWriter, 404, map[string]any{
			"error": "Refresh Token Not Found",
		})
		return
	}

	appUserId := persistentGrant.AppUserId
	sessionId := persistentGrant.SessionId

	refreshToken, httpErr := tokenService.CreateRefreshToken(ctx, now, body["clientId"].(string), appUserId, sessionId)
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}

	signedJwt, expirationDateTime, httpErr := tokenService.CreateAccessToken(now, body["clientId"].(string), appUserId)
	if httpErr != nil {
		replyAsJson(responseWriter, httpErr.StatusCode(), map[string]any{
			"error": httpErr.Error(),
		})
		return
	}

	err = transaction.Commit(ctx)
	if err != nil {
		log.Print("failed to save to database", err)
		http.Error(responseWriter, "Error saving changes", 500)
		return
	}

	replyAsJson(responseWriter, 200, map[string]any{
		"appUser": map[string]string{
			"id": appUserId,
		},
		"accessTokenExpirationDateTime": expirationDateTime.Format(time.RFC3339),
		"accessToken":                   signedJwt,
		"refreshToken":                  refreshToken,
	})
}
