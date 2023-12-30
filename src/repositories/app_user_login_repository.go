package repositories

import (
	"context"
	"family-expenses-api/models"
	"github.com/jackc/pgx/v4"
	"time"
)

type AppUserLoginRepository struct {
	transaction pgx.Tx
}

func CreateAppUserLoginRepository(transaction pgx.Tx) *AppUserLoginRepository {
	return &AppUserLoginRepository{
		transaction: transaction,
	}
}

func (repository *AppUserLoginRepository) GetByIdentityProviderAndKey(ctx context.Context, identityProvider string, key string) (models.AppUserLogin, error) {
	rows, err := repository.transaction.Query(ctx, "SELECT id, creation_unix, identity_provider, key, fk_app_user_id FROM app_user_login WHERE identity_provider = $1 AND key = $2 ", identityProvider, key)
	defer rows.Close()
	if err != nil {
		return models.AppUserLogin{}, err
	}
	if !rows.Next() {
		return models.AppUserLogin{}, nil
	}
	obj := models.AppUserLogin{}
	creationUnix := time.Time{}.Unix()
	err = rows.Scan(&obj.Id, &creationUnix, &obj.IdentityProvider, &obj.Key, &obj.UserId)
	if err != nil {
		return models.AppUserLogin{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	return obj, nil
}

func (repository *AppUserLoginRepository) Insert(ctx context.Context, appUserLogin models.AppUserLogin) (models.AppUserLogin, error) {
	_, err := repository.transaction.Exec(
		ctx,
		"INSERT INTO app_user_login (id, creation_unix, identity_provider, key, fk_app_user_id) VALUES ($1, $2,$3, $4, $5)",
		appUserLogin.Id, validTimeToUnixOrNil(appUserLogin.CreationDateTime), appUserLogin.IdentityProvider, appUserLogin.Key, appUserLogin.UserId,
	)
	if err != nil {
		return models.AppUserLogin{}, err
	}
	return appUserLogin, nil
}
