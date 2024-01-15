package services

import (
	"context"
	"family-expenses-api/core"
	"family-expenses-api/models"
	"family-expenses-api/repositories"
	"github.com/google/uuid"
	"log"
	"time"
)

type ChargesModelService struct {
	chargesModelRepository *repositories.ChargesModelRepository
	user                   core.AuthenticatedUser
}

func CreateChargesModelService(chargesModelRepository *repositories.ChargesModelRepository, user core.AuthenticatedUser) *ChargesModelService {
	return &ChargesModelService{
		chargesModelRepository,
		user,
	}
}

func (service *ChargesModelService) Create(ctx context.Context, newEntity models.ChargesModel) (models.ChargesModel, HttpError) {
	newEntity.Id = uuid.New().String()
	newEntity.CreationDateTime = time.Now()
	entity, err := service.chargesModelRepository.Insert(ctx, service.user.Id, newEntity)
	if err != nil {
		log.Println(err)
		return models.ChargesModel{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.ChargesModel{}, BadRequest{"unknown"}
	}
	return entity, nil
}
func (service *ChargesModelService) GetAll(ctx context.Context, workspaceId string) ([]models.ChargesModel, HttpError) {
	entities, err := service.chargesModelRepository.GetByWorkspaceId(ctx, service.user.Id, workspaceId)
	if err != nil {
		log.Println(err)
		return []models.ChargesModel{}, InternalServerError{"Internal Server Error"}
	}
	return entities, nil
}

func (service *ChargesModelService) GetOne(ctx context.Context, entityId string) (models.ChargesModel, HttpError) {
	entity, err := service.chargesModelRepository.GetById(ctx, service.user.Id, entityId)
	if err != nil {
		log.Println(err)
		return models.ChargesModel{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.ChargesModel{}, NotFound{"Charges Model not found"}
	}
	return entity, nil
}

func (service *ChargesModelService) Update(ctx context.Context, existentEntity models.ChargesModel) (models.ChargesModel, HttpError) {
	entity, err := service.chargesModelRepository.Update(ctx, service.user.Id, existentEntity)
	if err != nil {
		log.Println(err)
		return models.ChargesModel{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.ChargesModel{}, BadRequest{"unknown"}
	}
	return entity, nil
}

func (service *ChargesModelService) Delete(ctx context.Context, id string) HttpError {
	wasSucceeded, err := service.chargesModelRepository.DeleteById(ctx, service.user.Id, id)
	if err != nil {
		log.Println(err)
		return InternalServerError{"Internal Server Error"}
	}
	if !wasSucceeded {
		return BadRequest{"unknown"}
	}
	return nil
}
