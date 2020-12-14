package graph
//go:generate go run github.com/99designs/gqlgen
import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/brandonforster/resolver/graph/model"
	"github.com/brandonforster/resolver/internal/spamhaus"
)

// TODO: have this backed by a SQLite DB and not an array
type Resolver struct{
	IPs []*model.IPDetails
}

// getFromDB returns the IP details if they exist in the DB. Otherwise, it returns nil.
func (r Resolver) getFromDB(IP string) *model.IPDetails {
	for _, detail := range r.IPs {
		if IP == detail.IPAddress {
			return detail
		}
	}

	return nil
}

// TODO: this does not even begin to do what enqueue should
func (r *Resolver) Store(IP string) bool {
	details := r.getFromDB(IP)
	r.IPs = append(r.IPs, details)

	return true
}

// Get will return the IPDetails of a given IP address. If that IP is known, it returns details from the DB, otherwise
// it will do a lookup via the Internet.
//
// IP is an IPv4 formatted address to be queried.
//
// Returns the IPDetails if the query executed successfully; an error otherwise.
func (r *Resolver) Get(IP string) (*model.IPDetails, error) {
	value := r.getFromDB(IP)
	if value == nil {
		var err error
		value, err = newIPLookup(IP)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

func newIPLookup(IP string) (*model.IPDetails, error) {
	responseCode, err := spamhaus.Lookup(IP)
	if err != nil {
		return nil, err
	}

	details := model.IPDetails{
		UUID:         uuid.New().String(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ResponseCode: responseCode,
		IPAddress:    IP,
	}

	return &details, nil
}

func isAuthorized(ctx context.Context) bool {
	isAuth, _ := ctx.Value("isAuth").(bool)

	return isAuth
}
