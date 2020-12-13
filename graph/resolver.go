package graph
//go:generate go run github.com/99designs/gqlgen
import (
	"time"

	"github.com/google/uuid"

	"github.com/brandonforster/resolver/graph/model"
)

// TODO: have this backed by a SQLite DB and not an array
type Resolver struct{
	IPs []*model.IPDetails
}

func (r Resolver) lookup(IP string) model.IPDetails {
	var details model.IPDetails
	exists := false
	for _, detail := range r.IPs {
		if IP == detail.IPAddress {
			details = *detail
			exists = true
			break
		}
	}

	if !exists {
		details = model.IPDetails{
			UUID:         uuid.New().String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			ResponseCode: "",
			IPAddress:    IP,
		}
	}

	return details
}

// TODO: this does not even begin to do what enqueue should
func (r *Resolver) Store(IP string) bool {
	details := r.lookup(IP)
	r.IPs = append(r.IPs, &details)

	return true
}

func (r *Resolver) Get(IP string) (*model.IPDetails, error) {
	retval := r.lookup(IP)
	return &retval, nil
}
