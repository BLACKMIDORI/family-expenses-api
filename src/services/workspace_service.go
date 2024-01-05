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

type WorkspaceService struct {
	workspaceRepository *repositories.WorkspaceRepository
	user                core.AuthenticatedUser
}

func CreateWorkspaceService(workspaceRepository *repositories.WorkspaceRepository, user core.AuthenticatedUser) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepository,
		user,
	}
}

func (service *WorkspaceService) Create(ctx context.Context, newEntity models.Workspace) (models.Workspace, HttpError) {
	newEntity.Id = uuid.New().String()
	newEntity.CreationDateTime = time.Now()
	newEntity.OwnerId = service.user.Id
	entity, err := service.workspaceRepository.Insert(ctx, newEntity)
	if err != nil {
		log.Println(err)
		return models.Workspace{}, InternalServerError{"Internal Server Error"}
	}
	return entity, nil
}

func (service *WorkspaceService) GetById(ctx context.Context, id string) (models.Workspace, HttpError) {
	entity, err := service.workspaceRepository.GetById(ctx, id)
	if err != nil {
		log.Println(err)
		return models.Workspace{}, InternalServerError{"Internal Server Error"}
	}
	if entity.Id == "" {
		return models.Workspace{}, NotFound{"Workspace not found"}
	}
	return entity, nil
}

func (service *WorkspaceService) GetAllByUser(ctx context.Context, appUserId string) ([]models.Workspace, HttpError) {
	entities, err := service.workspaceRepository.GetAllByUserId(ctx, appUserId)
	if err != nil {
		log.Println(err)
		return []models.Workspace{}, InternalServerError{"Internal Server Error"}
	}
	return entities, nil
}
