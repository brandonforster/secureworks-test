package graph

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/brandonforster/resolver/graph/model"
	"github.com/brandonforster/resolver/internal/spamhaus"
	"github.com/brandonforster/resolver/internal/sqlite"
)

const FILENAME= `C:\Users\brandon\Desktop\resolver.db`

type Resolver struct {
	client *sqlite.Client
}

// TODO: this does not even begin to do what enqueue should
func (r *Resolver) Store(IP string) bool {
	return true
}

// Get will return the IPDetails of a given IP address. If that IP is known, it returns details from the DB, otherwise
// it will do a lookup via the Internet.
//
// IP is an IPv4 formatted address to be queried.
//
// Returns the IPDetails if the query executed successfully; an error otherwise.
func (r *Resolver) Get(IP string) (*model.IPDetails, error) {
	var err error
	r.client, err = sqlite.NewClient(FILENAME)
	if err != nil {
		return nil, err
	}

	details, err := r.client.GetIPDetailByAddress(IP)
	if err != nil {
		// we do not yet have this in the DB, do a Spamhaus lookup
		if strings.Contains(err.Error(), "no rows in result") {
			details, err = newIPLookup(IP)
			if err != nil {
				return nil, err
			}
		} else {
			// all other errors besides "not yet in DB"
			return nil, err
		}

	}

	return details, nil
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
