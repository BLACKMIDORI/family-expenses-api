package repositories

import (
	"context"
	"family-expenses-api/models"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strconv"
	"time"
)

type PayerPaymentWeightRepository struct {
	transaction pgx.Tx
}

func CreatePayerPaymentWeightRepository(transaction pgx.Tx) *PayerPaymentWeightRepository {
	return &PayerPaymentWeightRepository{
		transaction: transaction,
	}
}

func (repository *PayerPaymentWeightRepository) Insert(ctx context.Context, appUserId string, entity models.PayerPaymentWeight) (models.PayerPaymentWeight, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"INSERT INTO payer_payment_weight "+
			"(id, creation_unix, weight, fk_payer_id, fk_charge_association_id) "+
			"SELECT $4, $5, $6, payer.id, charge_association.id FROM charge_association "+
			"INNER JOIN charges_model "+
			"ON charges_model.id = charge_association.fk_charges_model_id "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"INNER JOIN payer "+
			"ON payer.fk_workspace_id = workspace.id  "+
			"WHERE workspace.fk_app_user_id = $1 and charge_association.id = $2 and payer.id = $3",
		appUserId, entity.ChargeAssociation.Id, entity.Payer.Id, entity.Id, validTimeToUnixOrNil(entity.CreationDateTime), entity.Weight,
	)
	if err != nil {
		return models.PayerPaymentWeight{}, err
	}
	if result.RowsAffected() == 0 {
		return models.PayerPaymentWeight{}, nil
	}

	return entity, nil
}

func (repository *PayerPaymentWeightRepository) GetByChargeAssociationId(ctx context.Context, appUserId string, chargeAssociationId string) (list []models.PayerPaymentWeight, err error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT payer_payment_weight.id, payer_payment_weight.creation_unix, payer_payment_weight.weight, payer_payment_weight.fk_payer_id, payer_payment_weight.fk_charge_association_id "+
			"FROM payer_payment_weight "+
			"INNER JOIN charge_association "+
			"ON charge_association.id = payer_payment_weight.fk_charge_association_id "+
			"INNER JOIN charges_model "+
			"ON charges_model.id = charge_association.fk_charges_model_id "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and charge_association.id = $2 "+
			"ORDER BY payer_payment_weight.creation_unix",
		appUserId, chargeAssociationId,
	)
	defer rows.Close()
	if err != nil {
		return []models.PayerPaymentWeight{}, err
	}
	if !rows.Next() {
		return []models.PayerPaymentWeight{}, nil
	}
	obj := models.PayerPaymentWeight{}
	creationUnix := time.Time{}.Unix()
	var weight float32
	err = rows.Scan(&obj.Id, &creationUnix, &weight, &obj.Payer.Id, &obj.ChargeAssociation.Id)
	if err != nil {
		return []models.PayerPaymentWeight{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	obj.Weight, _ = strconv.ParseFloat(fmt.Sprint(weight), 64)
	list = append(list, obj)
	for rows.Next() {
		obj := models.PayerPaymentWeight{}
		creationUnix := time.Time{}.Unix()
		var weight float32
		err = rows.Scan(&obj.Id, &creationUnix, &weight, &obj.Payer.Id, &obj.ChargeAssociation.Id)
		if err != nil {
			return []models.PayerPaymentWeight{}, err
		}
		obj.CreationDateTime = time.Unix(creationUnix, 0)
		obj.Weight, _ = strconv.ParseFloat(fmt.Sprint(weight), 64)
		list = append(list, obj)
	}
	return list, nil
}

func (repository *PayerPaymentWeightRepository) GetById(ctx context.Context, appUserId string, entityId string) (models.PayerPaymentWeight, error) {
	rows, err := repository.transaction.Query(
		ctx,
		"SELECT payer_payment_weight.id, payer_payment_weight.creation_unix, payer_payment_weight.weight, payer_payment_weight.fk_payer_id, payer_payment_weight.fk_charge_association_id "+
			"FROM payer_payment_weight "+
			"INNER JOIN charge_association "+
			"ON charge_association.id = payer_payment_weight.fk_charge_association_id "+
			"INNER JOIN charges_model "+
			"ON charges_model.id = charge_association.fk_charges_model_id "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and payer_payment_weight.id = $2",
		appUserId, entityId,
	)
	defer rows.Close()
	if err != nil {
		return models.PayerPaymentWeight{}, err
	}
	if !rows.Next() {
		return models.PayerPaymentWeight{}, nil
	}
	obj := models.PayerPaymentWeight{}
	creationUnix := time.Time{}.Unix()
	var weight float32
	err = rows.Scan(&obj.Id, &creationUnix, &weight, &obj.Payer.Id, &obj.ChargeAssociation.Id)
	if err != nil {
		return models.PayerPaymentWeight{}, err
	}
	obj.CreationDateTime = time.Unix(creationUnix, 0)
	obj.Weight, _ = strconv.ParseFloat(fmt.Sprint(weight), 64)
	return obj, nil
}

func (repository *PayerPaymentWeightRepository) Update(ctx context.Context, appUserId string, entity models.PayerPaymentWeight) (models.PayerPaymentWeight, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"UPDATE payer_payment_weight SET weight = $5, fk_payer_id = $3 "+
			"WHERE payer_payment_weight.id IN ( "+
			"SELECT payer_payment_weight.id FROM payer_payment_weight "+
			"INNER JOIN charge_association "+
			"ON charge_association.id = payer_payment_weight.fk_charge_association_id "+
			"INNER JOIN charges_model "+
			"ON charges_model.id = charge_association.fk_charges_model_id "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"INNER JOIN payer "+
			"ON payer.fk_workspace_id = workspace.id "+
			"WHERE workspace.fk_app_user_id = $1 and charge_association.id = $2 and payer.id = $3 and payer_payment_weight.id = $4 "+
			")",
		appUserId, entity.ChargeAssociation.Id, entity.Payer.Id, entity.Id, entity.Weight,
	)
	if err != nil {
		return models.PayerPaymentWeight{}, err
	}
	if result.RowsAffected() == 0 {
		return models.PayerPaymentWeight{}, nil
	}

	return entity, nil
}

func (repository *PayerPaymentWeightRepository) DeleteById(ctx context.Context, appUserId string, entityId string) (bool, error) {
	result, err := repository.transaction.Exec(
		ctx,
		"DELETE FROM payer_payment_weight "+
			"WHERE payer_payment_weight.id IN ( "+
			"SELECT payer_payment_weight.id FROM payer_payment_weight "+
			"INNER JOIN charge_association "+
			"ON charge_association.id = payer_payment_weight.fk_charge_association_id "+
			"INNER JOIN charges_model "+
			"ON charges_model.id = charge_association.fk_charges_model_id "+
			"INNER JOIN workspace "+
			"ON workspace.id = charges_model.fk_workspace_id "+
			"WHERE workspace.fk_app_user_id = $1 and payer_payment_weight.id = $2 "+
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
