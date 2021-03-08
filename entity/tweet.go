package entity

import "time"

//Tweet docs
type Tweet struct {
	Message string    `json:"message"`
	Created time.Time `json:"created,omitempty"`
}
