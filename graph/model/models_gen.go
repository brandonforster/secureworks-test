// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type IPDetails struct {
	UUID         string    `json:"uuid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ResponseCode string    `json:"response_code"`
	IPAddress    string    `json:"ip_address"`
}
