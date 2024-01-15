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

type ChargeAssociationService struct {
	chargeAssociationRepository *repositories.ChargeAssociationRepository
	user                        core.AuthenticatedUser
}

func CreateChargeAssociationService(chargeAssociationRepository *repositories.ChargeAssociationRepository, user core.AuthenticatedUser) *ChargeAssociationService {
	return &ChargeAssociationService{
		chargeAssociationRepository,
		user,
	}
}

func (service *ChargeAssociationService) Create(ctx context.Context, newEntity models.ChargeAssociation) (models.ChargeAssociation, HttpError) {
	newEntity.Id = uuid.New().String()
	newEntity.CreationDateTime = time.Now()
	entity, err := service.chargeAssociationRepository.Insert(ctx, service.user.Id, newEntity)
	if err != nil {
		log.Println(err)
		return models.ChargeAssociation{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.ChargeAssociation{}, BadRequest{"unknown"}
	}
	return entity, nil
}
func (service *ChargeAssociationService) GetAll(ctx context.Context, chargesModelId string) ([]models.ChargeAssociation, HttpError) {
	entities, err := service.chargeAssociationRepository.GetByChargesModelId(ctx, service.user.Id, chargesModelId)
	if err != nil {
		log.Println(err)
		return []models.ChargeAssociation{}, InternalServerError{"Internal Server Error"}
	}
	return entities, nil
}

func (service *ChargeAssociationService) GetOne(ctx context.Context, entityId string) (models.ChargeAssociation, HttpError) {
	entity, err := service.chargeAssociationRepository.GetById(ctx, service.user.Id, entityId)
	if err != nil {
		log.Println(err)
		return models.ChargeAssociation{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.ChargeAssociation{}, NotFound{"Charge Association not found"}
	}
	return entity, nil
}

func (service *ChargeAssociationService) Update(ctx context.Context, existentEntity models.ChargeAssociation) (models.ChargeAssociation, HttpError) {
	entity, err := service.chargeAssociationRepository.Update(ctx, service.user.Id, existentEntity)
	if err != nil {
		log.Println(err)
		return models.ChargeAssociation{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.ChargeAssociation{}, BadRequest{"unknown"}
	}
	return entity, nil
}

func (service *ChargeAssociationService) Delete(ctx context.Context, id string) HttpError {
	wasSucceeded, err := service.chargeAssociationRepository.DeleteById(ctx, service.user.Id, id)
	if err != nil {
		log.Println(err)
		return InternalServerError{"Internal Server Error"}
	}
	if !wasSucceeded {
		return BadRequest{"unknown"}
	}
	return nil
}
