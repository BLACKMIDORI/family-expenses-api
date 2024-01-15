package repositories

import (
	"context"
	"family-expenses-api/models"
	"github.com/jackc/pgx/v4"
	"time"
)

type PayerRepository struct {
	transaction pgx.Tx
}

func CreatePayerRepository(transaction pgx.Tx) *PayerRepository {
	return &PayerRepository{
		transaction: transaction,
	}
}

func (repository *PayerRepository) Insert(ctx context.Context, appUserId string, entity models.Payer) (models.Payer, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"INSERT INTO payer "+
			"(id, creation_unix, name, fk_workspace_id) "+
			"SELECT $3, $4, $5, id FROM workspace "+
			"WHERE fk_app_user_id = $1 and id = $2",
		appUserId, entity.Workspace.Id, entity.Id, validTimeToUnixOrNil(entity.CreationDateTime), entity.Name,
	)
	if err != nil {
		return models.Payer{}, err
	}
	if result.RowsAffected() == 0 {
		return models.Payer{}, nil
	}

	return entity, nil
}

func (repository *PayerRepository) GetByWorkspaceId(ctx context.Context, appUserId string, workspaceId string) (list []models.Payer, err error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT payer.id, payer.creation_unix, payer.name, payer.fk_workspace_id "+
			"FROM payer INNER JOIN workspace "+
			"ON workspace.id = payer.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and workspace.id = $2 "+
			"ORDER BY payer.creation_unix",
		appUserId, workspaceId,
	)
	defer rows.Close()
	if err != nil {
		return []models.Payer{}, err
	}
	if !rows.Next() {
		return []models.Payer{}, nil
	}
	obj := models.Payer{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Workspace.Id)
	if err != nil {
		return []models.Payer{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	list = append(list, obj)
	for rows.Next() {
		obj := models.Payer{}
		creationUnix := time.Time{}.Unix()
		err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Workspace.Id)
		if err != nil {
			return []models.Payer{}, err
		}
		obj.CreationDateTime = time.Unix(creationUnix, 0)
		list = append(list, obj)
	}
	return list, nil
}

func (repository *PayerRepository) GetById(ctx context.Context, appUserId string, entityId string) (models.Payer, error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT payer.id, payer.creation_unix, payer.name, payer.fk_workspace_id "+
			"FROM payer INNER JOIN workspace "+
			"ON workspace.id = payer.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and payer.id = $2",
		appUserId, entityId,
	)
	defer rows.Close()
	if err != nil {
		return models.Payer{}, err
	}
	if !rows.Next() {
		return models.Payer{}, nil
	}
	obj := models.Payer{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Workspace.Id)
	if err != nil {
		return models.Payer{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	return obj, nil
}

func (repository *PayerRepository) Update(ctx context.Context, appUserId string, entity models.Payer) (models.Payer, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"UPDATE payer SET name = $4 "+
			"WHERE payer.id IN ("+
			"SELECT payer.id FROM payer "+
			"INNER JOIN workspace "+
			"ON workspace.id = payer.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and workspace.id = $2 and payer.id = $3 "+
			")",
		appUserId, entity.Workspace.Id, entity.Id, entity.Name,
	)
	if err != nil {
		return models.Payer{}, err
	}
	if result.RowsAffected() == 0 {
		return models.Payer{}, nil
	}

	return entity, nil
}

func (repository *PayerRepository) DeleteById(ctx context.Context, appUserId string, entityId string) (bool, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"DELETE FROM payer "+
			"WHERE payer.id IN ("+
			"SELECT payer.id FROM payer "+
			"INNER JOIN workspace "+
			"ON workspace.id = payer.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and payer.id = $2 "+
			")",
		appUserId, entityId,
	)
	if err != nil {
		return false, err
	}
	if result.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}
