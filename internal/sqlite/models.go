package sqlite

import (
	"strings"
	"time"

	"github.com/brandonforster/resolver/graph/model"
)

const SEPARATOR = ","

type IPDetails struct {
	ID         string
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	ResponseCode string `db:"response_code"`
	IPAddress    string `db:"ip_address"`
}

func (d *IPDetails) fromContract(cd model.IPDetails) {
	d.ID = cd.UUID
	d.CreatedAt = cd.CreatedAt
	d.UpdatedAt = cd.UpdatedAt
	// storing the array as CSV is easy and elegant
	d.ResponseCode = strings.Join(cd.ResponseCode, SEPARATOR)
	d.IPAddress = cd.IPAddress
}

func (d *IPDetails) toContract() model.IPDetails {
	return model.IPDetails{
		UUID:         d.ID,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		// retrieving the array as CSV is easy and elegant
		ResponseCode: strings.Split(d.ResponseCode, SEPARATOR),
		IPAddress:    d.IPAddress,
	}
}
