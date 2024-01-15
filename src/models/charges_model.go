package models

import "time"

type ChargesModel struct {
	Id               string           `json:"id"`
	CreationDateTime time.Time        `json:"creationDateTime"`
	Name             string           `json:"name"`
	Workspace        ForeignKeyHolder `json:"workspace"`
}
