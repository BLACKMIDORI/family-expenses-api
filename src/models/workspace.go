package models

import "time"

type Workspace struct {
	Id               string    `json:"id"`
	CreationDateTime time.Time `json:"creationDateTime"`
	Name             string    `json:"name"`
	Owner            Owner     `json:"owner"`
}

type Owner struct {
	Id string `json:"id"`
}
