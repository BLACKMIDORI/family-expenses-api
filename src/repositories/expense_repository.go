package repositories

import (
	"context"
	"family-expenses-api/models"
	"github.com/jackc/pgx/v4"
	"time"
)

type ExpenseRepository struct {
	transaction pgx.Tx
}

func CreateExpenseRepository(transaction pgx.Tx) *ExpenseRepository {
	return &ExpenseRepository{
		transaction: transaction,
	}
}

func (repository *ExpenseRepository) Insert(ctx context.Context, appUserId string, entity models.Expense) (models.Expense, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"INSERT INTO expense "+
			"(id, creation_unix, name, fk_workspace_id) "+
			"SELECT $3, $4, $5, id FROM workspace "+
			"where fk_app_user_id = $1 and id = $2",
		appUserId, entity.Workspace.Id, entity.Id, validTimeToUnixOrNil(entity.CreationDateTime), entity.Name,
	)
	if err != nil {
		return models.Expense{}, err
	}
	if result.RowsAffected() == 0 {
		return models.Expense{}, nil
	}

	return entity, nil
}

func (repository *ExpenseRepository) GetByWorkspaceId(ctx context.Context, appUserId string, workspaceId string) (list []models.Expense, err error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT expense.id, expense.creation_unix, expense.name, expense.fk_workspace_id "+
			"FROM expense INNER JOIN workspace "+
			"ON workspace.id = expense.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and workspace.id = $2 "+
			"ORDER BY expense.creation_unix",
		appUserId, workspaceId,
	)
	defer rows.Close()
	if err != nil {
		return []models.Expense{}, err
	}
	if !rows.Next() {
		return []models.Expense{}, nil
	}
	obj := models.Expense{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Workspace.Id)
	if err != nil {
		return []models.Expense{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	list = append(list, obj)
	for rows.Next() {
		obj := models.Expense{}
		creationUnix := time.Time{}.Unix()
		err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Workspace.Id)
		if err != nil {
			return []models.Expense{}, err
		}
		obj.CreationDateTime = time.Unix(creationUnix, 0)
		list = append(list, obj)
	}
	return list, nil
}

func (repository *ExpenseRepository) GetById(ctx context.Context, appUserId string, entityId string) (models.Expense, error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT expense.id, expense.creation_unix, expense.name, expense.fk_workspace_id "+
			"FROM expense INNER JOIN workspace "+
			"ON workspace.id = expense.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and expense.id = $2",
		appUserId, entityId,
	)
	defer rows.Close()
	if err != nil {
		return models.Expense{}, err
	}
	if !rows.Next() {
		return models.Expense{}, nil
	}
	obj := models.Expense{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.Name, &obj.Workspace.Id)
	if err != nil {
		return models.Expense{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	return obj, nil
}

func (repository *ExpenseRepository) Update(ctx context.Context, appUserId string, entity models.Expense) (models.Expense, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"UPDATE expense SET name = $4 "+
			"WHERE expense.id IN ("+
			"SELECT expense.id FROM expense "+
			"INNER JOIN workspace "+
			"ON workspace.id = expense.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and workspace.id = $2 and expense.id = $3 "+
			")",
		appUserId, entity.Workspace.Id, entity.Id, entity.Name,
	)
	if err != nil {
		return models.Expense{}, err
	}
	if result.RowsAffected() == 0 {
		return models.Expense{}, nil
	}

	return entity, nil
}

func (repository *ExpenseRepository) DeleteById(ctx context.Context, appUserId string, entityId string) (bool, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"DELETE FROM expense "+
			"WHERE expense.id IN ("+
			"SELECT expense.id FROM expense "+
			"INNER JOIN workspace "+
			"ON workspace.id = expense.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and expense.id = $2 "+
			")",
		appUserId, entityId,
	)
	if err != nil {
		return true, err
	}
	if result.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}
