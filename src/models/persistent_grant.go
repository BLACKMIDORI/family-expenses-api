package models

import "time"

type PersistedGrant struct {
	Id                 string
	CreationDateTime   time.Time
	KeyDigest          string
	ClientId           string
	AppUserId          string
	SessionId          string
	ExpirationDateTime time.Time
	ConsumedDateTime   time.Time
}
