package models

import "time"

type AppUser struct {
	Id               string    `json:"id"`
	CreationDateTime time.Time `json:"creationDateTime"`
}
