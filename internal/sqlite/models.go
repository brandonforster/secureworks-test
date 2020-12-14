package sqlite

import (
	"strings"
	"time"

	"github.com/brandonforster/resolver/graph/model"
)

const SEPARATOR = ","

// IPDetails is a struct designed to make interacting with the database easier.
// The struct used by most of the system is close, but not quite what we need for interacting with the DB.
type IPDetails struct {
	ID           string
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	ResponseCode string    `db:"response_code"`
	IPAddress    string    `db:"ip_address"`
}

// use this method any time you have a system-wide representation of an object and need a database specific data model
func (d *IPDetails) fromContract(cd model.IPDetails) {
	d.ID = cd.UUID
	d.CreatedAt = cd.CreatedAt
	d.UpdatedAt = cd.UpdatedAt
	// storing the array as CSV is easy and elegant
	d.ResponseCode = strings.Join(cd.ResponseCode, SEPARATOR)
	d.IPAddress = cd.IPAddress
}

// use this method any time you have a database specific representation of an object and need a system-wide data model
func (d *IPDetails) toContract() model.IPDetails {
	return model.IPDetails{
		UUID:      d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		// retrieving the array as CSV is easy and elegant
		ResponseCode: strings.Split(d.ResponseCode, SEPARATOR),
		IPAddress:    d.IPAddress,
	}
}
