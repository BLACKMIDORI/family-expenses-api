package repositories

import (
	"context"
	"family-expenses-api/models"
	"github.com/jackc/pgx/v4"
	"time"
)

type PersistentGrantRepository struct {
	transaction pgx.Tx
}

func CreatePersistentGrantRepository(transaction pgx.Tx) *PersistentGrantRepository {
	return &PersistentGrantRepository{
		transaction: transaction,
	}
}

func (repository *PersistentGrantRepository) Insert(ctx context.Context, entity models.PersistedGrant) (models.PersistedGrant, error) {
	_, err := repository.transaction.Exec(
		ctx,
		"INSERT INTO persistent_grant (id, creation_unix, key_digest, client_id, fk_app_user_id, session_id,expiration_unix, consumed_unix) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		entity.Id, validTimeToUnixOrNil(entity.CreationDateTime), entity.KeyDigest, entity.ClientId, entity.AppUserId, entity.SessionId, validTimeToUnixOrNil(entity.ExpirationDateTime), validTimeToUnixOrNil(entity.ConsumedDateTime),
	)
	if err != nil {
		return models.PersistedGrant{}, err
	}

	return entity, nil
}

func (repository *PersistentGrantRepository) GetByKeyDigest(ctx context.Context, keyDigest string) (models.PersistedGrant, error) {
	rows, err := repository.transaction.Query(ctx, "SELECT id, creation_unix, key_digest, client_id, fk_app_user_id, session_id,expiration_unix, consumed_unix FROM persistent_grant WHERE key_digest = $1", keyDigest)
	defer rows.Close()
	if err != nil {
		return models.PersistedGrant{}, err
	}
	if !rows.Next() {
		return models.PersistedGrant{}, nil
	}
	obj := models.PersistedGrant{}
	creationUnix := time.Time{}.Unix()
	expirationUnix := time.Time{}.Unix()
	consumedUnix := time.Time{}.Unix()
	nullableConsumedUnix := &consumedUnix
	err = rows.Scan(&obj.Id, &creationUnix, &obj.KeyDigest, &obj.ClientId, &obj.AppUserId, &obj.SessionId, &expirationUnix, &nullableConsumedUnix)
	if err != nil {
		return models.PersistedGrant{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	obj.ExpirationDateTime = time.Unix(expirationUnix, 0)
	obj.ConsumedDateTime = time.Unix(consumedUnix, 0)
	return obj, nil
}

func (repository *PersistentGrantRepository) UpdateConsumed(ctx context.Context, persistentGrantId string, consumedDateTime time.Time) error {
	_, err := repository.transaction.Exec(
		ctx,
		"UPDATE persistent_grant SET consumed_unix = $1 WHERE id = $2",
		consumedDateTime.Unix(), persistentGrantId,
	)
	if err != nil {
		return err
	}

	return nil
}
