package services

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"family-expenses-api/core"
	"family-expenses-api/models"
	"family-expenses-api/repositories"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

type TokenService struct {
	persistentGrantRepository *repositories.PersistentGrantRepository
	issuer                    string
	audiences                 []string
	jwtSigningKeyPem          string
	jwtSigningPublicKeyPem    string
	googleConfig              IdentityProviderConfig
	appleConfig               IdentityProviderConfig
}
type IdentityProviderConfig struct {
	allowedIssuers   []string
	allowedAudiences []string
}

type VerifiedIdentity struct {
	IdentityProvider string
	ClientId         string
	Key              string
}

func CreateTokenService(persistentGrantRepository *repositories.PersistentGrantRepository) *TokenService {
	return &TokenService{
		persistentGrantRepository,
		os.Getenv("JWT_ISSUER"),
		strings.Split(os.Getenv("JWT_AUDIENCES"), ","),
		os.Getenv("JWT_SIGNING_KEY_PEM"),
		os.Getenv("JWT_SIGNING_PUBLIC_KEY_PEM"),
		IdentityProviderConfig{
			strings.Split(os.Getenv("GOOGLE_ALLOWED_ISSUERS"), ","),
			strings.Split(os.Getenv("GOOGLE_ALLOWED_AUDIENCES"), ","),
		},
		IdentityProviderConfig{
			strings.Split(os.Getenv("APPLE_ALLOWED_ISSUERS"), ","),
			strings.Split(os.Getenv("APPLE_ALLOWED_AUDIENCES"), ","),
		},
	}
}
func (service *TokenService) IsClientIdValid(clientId string) bool {
	if !slices.Contains(service.audiences, clientId) {
		return false
	}
	return true
}
func (service *TokenService) VerifyIdToken(idToken string) (VerifiedIdentity, HttpError) {
	// Source: https://stackoverflow.com/a/66574942
	jwksURL, err := getJwksUrlFromJwt(idToken)
	if err != nil {
		log.Println(err)
		return VerifiedIdentity{}, ServiceUnavailable{"Could not obtain identity provider configuration"}
	}

	// Create a context that, when cancelled, ends the JWKS background refresh goroutine.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the keyfunc options. Use an error handler that logs. Refresh the JWKS when a JWT signed by an unknown KID
	// is found or at the specified interval. Rate limit these refreshes. Timeout the initial JWKS refresh request after
	// 10 seconds. This timeout is also used to create the initial context.Context for keyfunc.Get.
	options := keyfunc.Options{
		Ctx: ctx,
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.Get(jwksURL, options)
	defer jwks.EndBackground()
	if err != nil {
		return VerifiedIdentity{}, BadRequest{
			Message: fmt.Sprintf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error()),
		}
	}

	// Parse the JWT.
	token, err := jwt.Parse(idToken, jwks.Keyfunc)
	if err != nil {
		return VerifiedIdentity{}, BadRequest{
			Message: fmt.Sprintf("Failed to parse the JWT.\nError: %s", err.Error()),
		}
	}

	// Check if the token is valid.
	if !token.Valid {
		return VerifiedIdentity{}, BadRequest{
			Message: "The token is not valid.",
		}
	}
	claims := token.Claims.(jwt.MapClaims)
	issuer := claims["iss"].(string)
	audience := claims["aud"].(string)
	subject := claims["sub"].(string)

	if slices.Contains(service.googleConfig.allowedIssuers, issuer) {
		if !slices.Contains(service.googleConfig.allowedAudiences, audience) {
			return VerifiedIdentity{}, BadRequest{
				Message: "Invalid idToken audience",
			}
		}
	} else if slices.Contains(service.appleConfig.allowedIssuers, issuer) {
		if !slices.Contains(service.appleConfig.allowedAudiences, audience) {
			return VerifiedIdentity{}, BadRequest{
				Message: "Invalid idToken audience",
			}
		}
	} else {
		return VerifiedIdentity{}, BadRequest{
			Message: "Invalid Issuer",
		}
	}

	return VerifiedIdentity{issuer, audience, subject}, nil
}

func getJwksUrlFromJwt(jwt string) (string, error) {
	parts := strings.Split(jwt, ".")
	payload := parts[1]
	bytes, err := base64.RawStdEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}
	jsonText := string(bytes)
	var jsonObj map[string]string
	_ = json.NewDecoder(strings.NewReader(jsonText)).Decode(&jsonObj)
	openIdConfigUrl := jsonObj["iss"] + "/.well-known/openid-configuration"
	response, err := http.Get(openIdConfigUrl)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Invalid status %v from %v", response.StatusCode, openIdConfigUrl))
	}
	var openIdConfig map[string]any
	err = json.NewDecoder(response.Body).Decode(&openIdConfig)
	if err != nil {
		return "", err
	}
	return openIdConfig["jwks_uri"].(string), nil
}

func (service *TokenService) CreateRefreshToken(ctx context.Context, creationDateTime time.Time, clientId string, appUserId string, sessionId string) (string, HttpError) {
	if sessionId == "" {
		sessionId = uuid.New().String()
	}
	refreshToken := ""
	for i := 0; i < 10; i++ {
		hash := sha512.New()
		hash.Write([]byte(uuid.New().String()))
		refreshToken += hex.EncodeToString(hash.Sum(nil))
	}
	hash := sha512.New()
	hash.Write([]byte(refreshToken))
	keyDigest := hex.EncodeToString(hash.Sum(nil))
	newEntity := models.PersistedGrant{
		Id:                 uuid.New().String(),
		CreationDateTime:   creationDateTime,
		KeyDigest:          keyDigest,
		ClientId:           clientId,
		AppUserId:          appUserId,
		SessionId:          sessionId,
		ExpirationDateTime: creationDateTime.Add(30 * 24 * time.Hour),
		ConsumedDateTime:   time.Time{},
	}
	_, err := service.persistentGrantRepository.Insert(ctx, newEntity)
	if err != nil {
		log.Println(err)
		return "", InternalServerError{"Internal Server Error"}
	}
	return refreshToken, nil
}

func (service *TokenService) CreateAccessToken(creationDateTime time.Time, clientId string, userId string) (string, time.Time, HttpError) {
	keyPem := service.jwtSigningKeyPem
	keyPem = strings.ReplaceAll(keyPem, "\\n", "\n")
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyPem))
	if err != nil {
		return "", time.Time{}, InternalServerError{err.Error()}
	}
	expirationDateTime := creationDateTime.Add(10 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": service.issuer,
		"aud": clientId,
		"sub": userId,
		"iat": creationDateTime.Unix(),
		"exp": expirationDateTime.Unix(),
	})
	signedJwt, err := token.SignedString(key)
	if err != nil {
		return "", time.Time{}, InternalServerError{err.Error()}
	}
	return signedJwt, expirationDateTime, nil
}

func (service *TokenService) VerifyAccessToken(accessToken string) (core.AuthenticatedUser, HttpError) {
	// Parse the JWT.
	token, err := jwt.Parse(
		accessToken,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
			}
			keyPem := service.jwtSigningPublicKeyPem
			keyPem = strings.ReplaceAll(keyPem, "\\n", "\n")
			return jwt.ParseRSAPublicKeyFromPEM([]byte(keyPem))
		},
	)
	if err != nil {
		return core.AuthenticatedUser{}, BadRequest{
			Message: fmt.Sprintf("Failed to parse the JWT.\nError: %s", err.Error()),
		}
	}

	// Check if the token is valid.
	if !token.Valid {
		return core.AuthenticatedUser{}, BadRequest{
			Message: "The token is not valid.",
		}
	}
	claims := token.Claims.(jwt.MapClaims)

	if claims["iss"] != service.issuer {
		return core.AuthenticatedUser{}, BadRequest{
			Message: "Invalid Issuer: " + claims["iss"].(string),
		}
	}
	if !slices.Contains(service.audiences, claims["aud"].(string)) {
		return core.AuthenticatedUser{}, BadRequest{
			Message: "Invalid Audience: " + claims["aud"].(string),
		}
	}

	return core.AuthenticatedUser{claims["sub"].(string)}, nil
}

func (service *TokenService) ConsumeRefreshToken(ctx context.Context, refreshToken string, consumedDateTime time.Time) (models.PersistedGrant, HttpError) {
	hash := sha512.New()
	hash.Write([]byte(refreshToken))
	keyDigest := hex.EncodeToString(hash.Sum(nil))
	entity, err := service.persistentGrantRepository.GetByKeyDigest(ctx, keyDigest)
	if err != nil {
		log.Println(err)
		return models.PersistedGrant{}, InternalServerError{"Internal Server Error"}
	}
	err = service.persistentGrantRepository.UpdateConsumed(ctx, entity.Id, consumedDateTime)
	if err != nil {
		log.Println(err)
		return models.PersistedGrant{}, InternalServerError{"Internal Server Error"}
	}
	return entity, nil
}
