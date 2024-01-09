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

func (repository *WorkspaceRepository) GetAll(ctx context.Context, appUserId string) (list []models.Workspace, err error) {
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

func (repository *WorkspaceRepository) GetById(ctx context.Context, appUserId string, entityId string) (models.Workspace, error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT id,creation_unix, name, fk_app_user_id FROM workspace WHERE fk_app_user_id = $1 and id = $2",
		appUserId, entityId,
	)
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

func (repository *WorkspaceRepository) Update(ctx context.Context, appUserId string, entity models.Workspace) (models.Workspace, error) {
	_, err := repository.transaction.Exec(
		ctx,
		"UPDATE workspace SET name = $3 WHERE fk_app_user_id = $1 and id = $2",
		appUserId, entity.Id, entity.Name,
	)
	if err != nil {
		return models.Workspace{}, err
	}

	return entity, nil
}

func (repository *WorkspaceRepository) DeleteById(ctx context.Context, appUserId string, entityId string) error {
	_, err := repository.transaction.Exec(
		ctx,
		"DELETE FROM workspace WHERE fk_app_user_id = $1 and id = $2",
		appUserId, entityId,
	)
	if err != nil {
		return err
	}

	return nil
}
