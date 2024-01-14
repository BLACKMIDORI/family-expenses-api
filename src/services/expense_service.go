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

type ExpenseService struct {
	expenseRepository *repositories.ExpenseRepository
	user              core.AuthenticatedUser
}

func CreateExpenseService(expenseRepository *repositories.ExpenseRepository, user core.AuthenticatedUser) *ExpenseService {
	return &ExpenseService{
		expenseRepository,
		user,
	}
}

func (service *ExpenseService) Create(ctx context.Context, newEntity models.Expense) (models.Expense, HttpError) {
	newEntity.Id = uuid.New().String()
	newEntity.CreationDateTime = time.Now()
	entity, err := service.expenseRepository.Insert(ctx, service.user.Id, newEntity)
	if err != nil {
		log.Println(err)
		return models.Expense{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.Expense{}, BadRequest{"unknown"}
	}
	return entity, nil
}
func (service *ExpenseService) GetAll(ctx context.Context, workspaceId string) ([]models.Expense, HttpError) {
	entities, err := service.expenseRepository.GetByWorkspaceId(ctx, service.user.Id, workspaceId)
	if err != nil {
		log.Println(err)
		return []models.Expense{}, InternalServerError{"Internal Server Error"}
	}
	return entities, nil
}

func (service *ExpenseService) GetOne(ctx context.Context, entityId string) (models.Expense, HttpError) {
	entity, err := service.expenseRepository.GetById(ctx, service.user.Id, entityId)
	if err != nil {
		log.Println(err)
		return models.Expense{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.Expense{}, NotFound{"Expense not found"}
	}
	return entity, nil
}

func (service *ExpenseService) Update(ctx context.Context, existentEntity models.Expense) (models.Expense, HttpError) {
	entity, err := service.expenseRepository.Update(ctx, service.user.Id, existentEntity)
	if err != nil {
		log.Println(err)
		return models.Expense{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.Expense{}, BadRequest{"unknown"}
	}
	return entity, nil
}

func (service *ExpenseService) Delete(ctx context.Context, id string) HttpError {
	wasSucceeded, err := service.expenseRepository.DeleteById(ctx, service.user.Id, id)
	if err != nil {
		log.Println(err)
		return InternalServerError{"Internal Server Error"}
	}
	if !wasSucceeded {
		return BadRequest{"unknown"}
	}
	return nil
}
