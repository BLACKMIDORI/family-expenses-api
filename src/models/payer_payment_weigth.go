package models

import "time"

type PayerPaymentWeight struct {
	Id                string           `json:"id"`
	CreationDateTime  time.Time        `json:"creationDateTime"`
	Weight            float64          `json:"weight"`
	Payer             ForeignKeyHolder `json:"payer"`
	ChargeAssociation ForeignKeyHolder `json:"chargeAssociation"`
}
