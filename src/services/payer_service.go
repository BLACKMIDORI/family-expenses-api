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

type PayerService struct {
	payerRepository *repositories.PayerRepository
	user            core.AuthenticatedUser
}

func CreatePayerService(payerRepository *repositories.PayerRepository, user core.AuthenticatedUser) *PayerService {
	return &PayerService{
		payerRepository,
		user,
	}
}

func (service *PayerService) Create(ctx context.Context, newEntity models.Payer) (models.Payer, HttpError) {
	newEntity.Id = uuid.New().String()
	newEntity.CreationDateTime = time.Now()
	entity, err := service.payerRepository.Insert(ctx, service.user.Id, newEntity)
	if err != nil {
		log.Println(err)
		return models.Payer{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.Payer{}, BadRequest{"unknown"}
	}
	return entity, nil
}
func (service *PayerService) GetAll(ctx context.Context, workspaceId string) ([]models.Payer, HttpError) {
	entities, err := service.payerRepository.GetByWorkspaceId(ctx, service.user.Id, workspaceId)
	if err != nil {
		log.Println(err)
		return []models.Payer{}, InternalServerError{"Internal Server Error"}
	}
	return entities, nil
}

func (service *PayerService) GetOne(ctx context.Context, entityId string) (models.Payer, HttpError) {
	entity, err := service.payerRepository.GetById(ctx, service.user.Id, entityId)
	if err != nil {
		log.Println(err)
		return models.Payer{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.Payer{}, NotFound{"Payer not found"}
	}
	return entity, nil
}

func (service *PayerService) Update(ctx context.Context, existentEntity models.Payer) (models.Payer, HttpError) {
	entity, err := service.payerRepository.Update(ctx, service.user.Id, existentEntity)
	if err != nil {
		log.Println(err)
		return models.Payer{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.Payer{}, BadRequest{"unknown"}
	}
	return entity, nil
}

func (service *PayerService) Delete(ctx context.Context, id string) HttpError {
	wasSucceeded, err := service.payerRepository.DeleteById(ctx, service.user.Id, id)
	if err != nil {
		log.Println(err)
		return InternalServerError{"Internal Server Error"}
	}
	if !wasSucceeded {
		return BadRequest{"unknown"}
	}
	return nil
}
