package repositories

import (
	"context"
	"family-expenses-api/models"
	"github.com/jackc/pgx/v4"
	"time"
)

type WorkspaceRepository struct {
	transaction pgx.Tx
}

func CreateWorkspaceRepository(transaction pgx.Tx) *WorkspaceRepository {
	return &WorkspaceRepository{
		transaction: transaction,
	}
}

func (repository *WorkspaceRepository) Insert(ctx context.Context, entity models.Workspace) (models.Workspace, error) {
	_, err := repository.transaction.Exec(
		ctx,
		"INSERT INTO workspace (id, creation_unix, name, fk_app_user_id) VALUES ($1, $2, $3, $4)",
		entity.Id, validTimeToUnixOrNil(entity.CreationDateTime), entity.Name, entity.OwnerId,
	)
	if err != nil {
		return models.Workspace{}, err
	}

	return entity, nil
}

func (repository *WorkspaceRepository) GetById(ctx context.Context, id string) (models.Workspace, error) {
	rows, err := repository.transaction.Query(ctx, "SELECT id,creation_unix, name, fk_app_user_id FROM workspace WHERE id = $1", id)
	defer rows.Close()
	if err != nil {
		return models.Workspace{}, err
	}
	if !rows.Next() {
		return models.Workspace{}, nil
	}
	obj := models.Workspace{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.OwnerId)
	if err != nil {
		return models.Workspace{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	return obj, nil
}

func (repository *WorkspaceRepository) GetAllByUserId(ctx context.Context, appUserId string) (list []models.Workspace, err error) {
	rows, err := repository.transaction.Query(ctx, "SELECT id, creation_unix, name, fk_app_user_id FROM workspace WHERE fk_app_user_id = $1 ORDER BY creation_unix", appUserId)
	defer rows.Close()
	if err != nil {
		return []models.Workspace{}, err
	}
	if !rows.Next() {
		return []models.Workspace{}, nil
	}
	obj := models.Workspace{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.OwnerId)
	if err != nil {
		return []models.Workspace{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	list = append(list, obj)
	for rows.Next() {
		obj := models.Workspace{}
		creationUnix := time.Time{}.Unix()
		err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.OwnerId)
		if err != nil {
			return []models.Workspace{}, err
		}
		obj.CreationDateTime = time.Unix(creationUnix, 0)
		list = append(list, obj)
	}
	return list, nil
}
