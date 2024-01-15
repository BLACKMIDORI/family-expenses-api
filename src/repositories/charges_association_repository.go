package repositories

import (
	"context"
	"family-expenses-api/models"
	"github.com/jackc/pgx/v4"
	"time"
)

type ChargeAssociationRepository struct {
	transaction pgx.Tx
}

func CreateChargeAssociationRepository(transaction pgx.Tx) *ChargeAssociationRepository {
	return &ChargeAssociationRepository{
		transaction: transaction,
	}
}

func (repository *ChargeAssociationRepository) Insert(ctx context.Context, appUserId string, entity models.ChargeAssociation) (models.ChargeAssociation, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"INSERT INTO charge_association "+
			"(id, creation_unix, name, fk_expense_id, fk_charges_model_id) "+
			"SELECT $4, $5, $6, expense.id, charges_model.id FROM charges_model "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"INNER JOIN expense "+
			"ON expense.fk_workspace_id = workspace.id  "+
			"WHERE workspace.fk_app_user_id = $1 and charges_model.id = $2 and expense.id = $3",
		appUserId, entity.ChargesModel.Id, entity.Expense.Id, entity.Id, validTimeToUnixOrNil(entity.CreationDateTime), entity.Name,
	)
	if err != nil {
		return models.ChargeAssociation{}, err
	}
	if result.RowsAffected() == 0 {
		return models.ChargeAssociation{}, nil
	}

	return entity, nil
}

func (repository *ChargeAssociationRepository) GetByChargesModelId(ctx context.Context, appUserId string, chargesModelId string) (list []models.ChargeAssociation, err error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT charge_association.id, charge_association.creation_unix, charge_association.name, charge_association.fk_expense_id, charge_association.fk_charges_model_id "+
			"FROM charge_association "+
			"INNER JOIN charges_model "+
			"ON charges_model.id = charge_association.fk_charges_model_id "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and charges_model.id = $2 "+
			"ORDER BY charge_association.creation_unix",
		appUserId, chargesModelId,
	)
	defer rows.Close()
	if err != nil {
		return []models.ChargeAssociation{}, err
	}
	if !rows.Next() {
		return []models.ChargeAssociation{}, nil
	}
	obj := models.ChargeAssociation{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Expense.Id, &obj.ChargesModel.Id)
	if err != nil {
		return []models.ChargeAssociation{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	list = append(list, obj)
	for rows.Next() {
		obj := models.ChargeAssociation{}
		creationUnix := time.Time{}.Unix()
		err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Expense.Id, &obj.ChargesModel.Id)
		if err != nil {
			return []models.ChargeAssociation{}, err
		}
		obj.CreationDateTime = time.Unix(creationUnix, 0)
		list = append(list, obj)
	}
	return list, nil
}

func (repository *ChargeAssociationRepository) GetById(ctx context.Context, appUserId string, entityId string) (models.ChargeAssociation, error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT charge_association.id, charge_association.creation_unix, charge_association.name, charge_association.fk_expense_id, charge_association.fk_charges_model_id "+
			"FROM charge_association "+
			"INNER JOIN charges_model "+
			"ON charges_model.id = charge_association.fk_charges_model_id "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and charge_association.id = $2",
		appUserId, entityId,
	)
	defer rows.Close()
	if err != nil {
		return models.ChargeAssociation{}, err
	}
	if !rows.Next() {
		return models.ChargeAssociation{}, nil
	}
	obj := models.ChargeAssociation{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Expense.Id, &obj.ChargesModel.Id)
	if err != nil {
		return models.ChargeAssociation{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	return obj, nil
}

func (repository *ChargeAssociationRepository) Update(ctx context.Context, appUserId string, entity models.ChargeAssociation) (models.ChargeAssociation, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"UPDATE charge_association SET name = $5, fk_expense_id = $3 "+
			"WHERE charge_association.id IN ( "+
			"SELECT charge_association.id FROM charge_association "+
			"INNER JOIN charges_model "+
			"ON charges_model.id = charge_association.fk_charges_model_id "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"INNER JOIN expense "+
			"ON expense.fk_workspace_id = workspace.id "+
			"WHERE workspace.fk_app_user_id = $1 and charges_model.id = $2 and expense.id = $3 and charge_association.id = $4 "+
			")",
		appUserId, entity.ChargesModel.Id, entity.Expense.Id, entity.Id, entity.Name,
	)
	if err != nil {
		return models.ChargeAssociation{}, err
	}
	if result.RowsAffected() == 0 {
		return models.ChargeAssociation{}, nil
	}

	return entity, nil
}

func (repository *ChargeAssociationRepository) DeleteById(ctx context.Context, appUserId string, entityId string) (bool, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"DELETE FROM charge_association "+
			"WHERE charge_association.id IN ( "+
			"SELECT charge_association.id FROM charge_association "+
			"INNER JOIN charges_model "+
			"ON charges_model.id = charge_association.fk_charges_model_id "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and charge_association.id = $2 "+
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
