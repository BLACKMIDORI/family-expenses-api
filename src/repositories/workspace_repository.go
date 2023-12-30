package repositories

import (
	"context"
	"family-expenses-api/models"
	"github.com/jackc/pgx/v4"
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
	rows, err := repository.transaction.Query(ctx, "SELECT id, name, fk_app_user_id FROM workspace WHERE id = $1", id)
	defer rows.Close()
	if err != nil {
		return models.Workspace{}, err
	}
	if !rows.Next() {
		return models.Workspace{}, nil
	}
	obj := models.Workspace{}
	err = rows.Scan(&obj.Id, &obj.Name, &obj.OwnerId)
	if err != nil {
		return models.Workspace{}, err
	}
	return obj, nil
}
