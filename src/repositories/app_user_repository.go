package repositories

import (
	"context"
	"family-expenses-api/models"
	"github.com/jackc/pgx/v4"
	"time"
)

type AppUserRepository struct {
	transaction pgx.Tx
}

func CreateAppUserRepository(transaction pgx.Tx) *AppUserRepository {
	return &AppUserRepository{
		transaction: transaction,
	}
}

func (repository *AppUserRepository) GetById(ctx context.Context, id string) (models.AppUser, error) {
	rows, err := repository.transaction.Query(ctx, "SELECT id, creation_unix FROM app_user WHERE id = $1", id)
	defer rows.Close()
	if err != nil {
		return models.AppUser{}, err
	}
	if !rows.Next() {
		return models.AppUser{}, nil
	}
	obj := models.AppUser{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix)
	if err != nil {
		return models.AppUser{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	return obj, nil
}

func (repository *AppUserRepository) Insert(ctx context.Context, appUser models.AppUser) (models.AppUser, error) {
	_, err := repository.transaction.Exec(
		ctx,
		"INSERT INTO app_user (id, creation_unix) VALUES ($1, $2)",
		appUser.Id,
		validTimeToUnixOrNil(appUser.CreationDateTime),
	)
	if err != nil {
		return models.AppUser{}, err
	}

	return appUser, nil
}
