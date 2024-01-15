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

type PayerPaymentWeightService struct {
	payerPaymentWeightRepository *repositories.PayerPaymentWeightRepository
	user                         core.AuthenticatedUser
}

func CreatePayerPaymentWeightService(payerPaymentWeightRepository *repositories.PayerPaymentWeightRepository, user core.AuthenticatedUser) *PayerPaymentWeightService {
	return &PayerPaymentWeightService{
		payerPaymentWeightRepository,
		user,
	}
}

func (service *PayerPaymentWeightService) Create(ctx context.Context, newEntity models.PayerPaymentWeight) (models.PayerPaymentWeight, HttpError) {
	newEntity.Id = uuid.New().String()
	newEntity.CreationDateTime = time.Now()
	entity, err := service.payerPaymentWeightRepository.Insert(ctx, service.user.Id, newEntity)
	if err != nil {
		log.Println(err)
		return models.PayerPaymentWeight{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.PayerPaymentWeight{}, BadRequest{"unknown"}
	}
	return entity, nil
}
func (service *PayerPaymentWeightService) GetAll(ctx context.Context, chargeAssociationId string) ([]models.PayerPaymentWeight, HttpError) {
	entities, err := service.payerPaymentWeightRepository.GetByChargeAssociationId(ctx, service.user.Id, chargeAssociationId)
	if err != nil {
		log.Println(err)
		return []models.PayerPaymentWeight{}, InternalServerError{"Internal Server Error"}
	}
	return entities, nil
}

func (service *PayerPaymentWeightService) GetOne(ctx context.Context, entityId string) (models.PayerPaymentWeight, HttpError) {
	entity, err := service.payerPaymentWeightRepository.GetById(ctx, service.user.Id, entityId)
	if err != nil {
		log.Println(err)
		return models.PayerPaymentWeight{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.PayerPaymentWeight{}, NotFound{"Payer Payment Weight not found"}
	}
	return entity, nil
}

func (service *PayerPaymentWeightService) Update(ctx context.Context, existentEntity models.PayerPaymentWeight) (models.PayerPaymentWeight, HttpError) {
	entity, err := service.payerPaymentWeightRepository.Update(ctx, service.user.Id, existentEntity)
	if err != nil {
		log.Println(err)
		return models.PayerPaymentWeight{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.PayerPaymentWeight{}, BadRequest{"unknown"}
	}
	return entity, nil
}

func (service *PayerPaymentWeightService) Delete(ctx context.Context, id string) HttpError {
	wasSucceeded, err := service.payerPaymentWeightRepository.DeleteById(ctx, service.user.Id, id)
	if err != nil {
		log.Println(err)
		return InternalServerError{"Internal Server Error"}
	}
	if !wasSucceeded {
		return BadRequest{"unknown"}
	}
	return nil
}
