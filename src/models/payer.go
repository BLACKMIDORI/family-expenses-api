package models

import "time"

type Payer struct {
	Id               string    `json:"id"`
	CreationDateTime time.Time `json:"creationDateTime"`
	Name             string    `json:"name"`
	Workspace        Workspace `json:"workspace"`
}
