package models

import "time"

type AppUserLogin struct {
	Id               string    `json:"id"`
	CreationDateTime time.Time `json:"creationDateTime"`
	IdentityProvider string    `json:"identityProvider"`
	Key              string    `json:"key"`
	UserId           string    `json:"userId"`
}
