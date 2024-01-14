package models

import "time"

type Expense struct {
	Id               string    `json:"id"`
	CreationDateTime time.Time `json:"creationDateTime"`
	Name             string    `json:"name"`
	Workspace        Workspace `json:"workspace"`
}
