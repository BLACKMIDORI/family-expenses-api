package repositories

import (
	"context"
	"family-expenses-api/models"
	"github.com/jackc/pgx/v4"
	"time"
)

type ChargesModelRepository struct {
	transaction pgx.Tx
}

func CreateChargesModelRepository(transaction pgx.Tx) *ChargesModelRepository {
	return &ChargesModelRepository{
		transaction: transaction,
	}
}

func (repository *ChargesModelRepository) Insert(ctx context.Context, appUserId string, entity models.ChargesModel) (models.ChargesModel, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"INSERT INTO charges_model "+
			"(id, creation_unix, name, fk_workspace_id) "+
			"SELECT $3, $4, $5, id FROM workspace "+
			"WHERE fk_app_user_id = $1 and id = $2",
		appUserId, entity.Workspace.Id, entity.Id, validTimeToUnixOrNil(entity.CreationDateTime), entity.Name,
	)
	if err != nil {
		return models.ChargesModel{}, err
	}
	if result.RowsAffected() == 0 {
		return models.ChargesModel{}, nil
	}

	return entity, nil
}

func (repository *ChargesModelRepository) GetByWorkspaceId(ctx context.Context, appUserId string, workspaceId string) (list []models.ChargesModel, err error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT charges_model.id, charges_model.creation_unix, charges_model.name, charges_model.fk_workspace_id "+
			"FROM charges_model INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and workspace.id = $2 "+
			"ORDER BY charges_model.creation_unix",
		appUserId, workspaceId,
	)
	defer rows.Close()
	if err != nil {
		return []models.ChargesModel{}, err
	}
	if !rows.Next() {
		return []models.ChargesModel{}, nil
	}
	obj := models.ChargesModel{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Workspace.Id)
	if err != nil {
		return []models.ChargesModel{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	list = append(list, obj)
	for rows.Next() {
		obj := models.ChargesModel{}
		creationUnix := time.Time{}.Unix()
		err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Workspace.Id)
		if err != nil {
			return []models.ChargesModel{}, err
		}
		obj.CreationDateTime = time.Unix(creationUnix, 0)
		list = append(list, obj)
	}
	return list, nil
}

func (repository *ChargesModelRepository) GetById(ctx context.Context, appUserId string, entityId string) (models.ChargesModel, error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT charges_model.id, charges_model.creation_unix, charges_model.name, charges_model.fk_workspace_id "+
			"FROM charges_model INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and charges_model.id = $2",
		appUserId, entityId,
	)
	defer rows.Close()
	if err != nil {
		return models.ChargesModel{}, err
	}
	if !rows.Next() {
		return models.ChargesModel{}, nil
	}
	obj := models.ChargesModel{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Workspace.Id)
	if err != nil {
		return models.ChargesModel{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	return obj, nil
}

func (repository *ChargesModelRepository) Update(ctx context.Context, appUserId string, entity models.ChargesModel) (models.ChargesModel, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"UPDATE charges_model SET name = $4 "+
			"WHERE charges_model.id IN ("+
			"SELECT charges_model.id FROM charges_model "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and workspace.id = $2 and charges_model.id = $3 "+
			")",
		appUserId, entity.Workspace.Id, entity.Id, entity.Name,
	)
	if err != nil {
		return models.ChargesModel{}, err
	}
	if result.RowsAffected() == 0 {
		return models.ChargesModel{}, nil
	}

	return entity, nil
}

func (repository *ChargesModelRepository) DeleteById(ctx context.Context, appUserId string, entityId string) (bool, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"DELETE FROM charges_model "+
			"WHERE charges_model.id IN ("+
			"SELECT charges_model.id FROM charges_model "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and charges_model.id = $2 "+
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
