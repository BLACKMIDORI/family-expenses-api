package models

import "time"

type ChargeAssociation struct {
	Id               string           `json:"id"`
	CreationDateTime time.Time        `json:"creationDateTime"`
	Name             string           `json:"name"`
	Expense          ForeignKeyHolder `json:"expense"`
	ChargesModel     ForeignKeyHolder `json:"chargesModel"`
}
