package models

import "time"

type Workspace struct {
	Id               string    `json:"id"`
	CreationDateTime time.Time `json:"creationDateTime"`
	Name             string    `json:"name"`
	OwnerId          string    `json:"ownerId"`
}
