package entity

import "time"

// Post docs
type Post struct {
	Message string    `json:"message"`
	Created time.Time `json:"created,omitempty"`
	Keys map[string]string `json:"keys,omitempty"`
}
