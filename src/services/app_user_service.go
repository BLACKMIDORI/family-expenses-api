package services

import (
	"context"
	"family-expenses-api/models"
	"family-expenses-api/repositories"
	"github.com/google/uuid"
	"log"
	"time"
)

type AppUserService struct {
	appUserRepository      *repositories.AppUserRepository
	appUserLoginRepository *repositories.AppUserLoginRepository
}

func CreateAppUserService(appUserRepository *repositories.AppUserRepository, appUserLoginRepository *repositories.AppUserLoginRepository) *AppUserService {
	return &AppUserService{
		appUserRepository,
		appUserLoginRepository,
	}
}

func (service *AppUserService) GetAppUserByLogin(ctx context.Context, identityProvider string, key string) (models.AppUser, HttpError) {
	appUserLogin, err := service.appUserLoginRepository.GetByIdentityProviderAndKey(ctx, identityProvider, key)
	if err != nil {
		log.Println(err)
		return models.AppUser{}, InternalServerError{"Internal Server Error"}
	}
	if appUserLogin.Id == "" {
		return models.AppUser{}, NotFound{"AppUserLogin not found"}
	}
	appUser, err := service.appUserRepository.GetById(ctx, appUserLogin.UserId)
	if err != nil {
		log.Println(err)
		return models.AppUser{}, InternalServerError{"Internal Server Error"}
	}

	return appUser, nil
}

func (service *AppUserService) CreateWithLogin(ctx context.Context, identityProvider string, key string) (models.AppUser, HttpError) {
	now := time.Now()
	newAppUser := models.AppUser{
		Id:               uuid.New().String(),
		CreationDateTime: now,
	}
	newAppUserLogin := models.AppUserLogin{
		Id:               uuid.New().String(),
		CreationDateTime: now,
		IdentityProvider: identityProvider,
		Key:              key,
		UserId:           newAppUser.Id,
	}
	appUser, err := service.appUserRepository.Insert(ctx, newAppUser)
	if err != nil {
		log.Println(err)
		return models.AppUser{}, InternalServerError{"Internal Server Error"}
	}
	_, err = service.appUserLoginRepository.Insert(ctx, newAppUserLogin)
	if err != nil {
		log.Println(err)
		return models.AppUser{}, InternalServerError{"Internal Server Error"}
	}

	return appUser, nil
}
