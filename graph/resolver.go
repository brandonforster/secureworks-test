package graph
//go:generate go run github.com/99designs/gqlgen
import (
	"context"
	"fmt"
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

func (r *Resolver) GetAndStore(IP string) (*model.IPDetails, error) {
	// does this IP exist in the DB?
	details, err := r.getFromDB(IP)
	if err != nil {
		return nil, err
	}

	if details == nil {
		// we do not yet have this in the DB, do a Spamhaus lookup
		details, err = newIPLookup(IP)
		if err != nil {
			return nil, err
		}
	} else {
		// this exists, update the response codes and LastUpdated
		responseCode, err := spamhaus.Lookup(IP)
		if err != nil {
			return nil, err
		}

		details.ResponseCode = responseCode
		details.UpdatedAt = time.Now()
	}

	err = r.storeToDB(*details)
	if err != nil {
		return nil, err
	}

	return details, nil
}

func (r *Resolver) Queue(ip string) (<-chan *model.IPDetails,<-chan error) {
	modelChan := make(chan *model.IPDetails)
	errChan := make(chan error)

	go func() {
		details, err := r.GetAndStore(ip)
		if err != nil {
			errChan <- err
		} else {
			modelChan <- details
		}

		close(modelChan)
		close(errChan)
	}()

	return modelChan, errChan
}

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
	details, err := r.getFromDB(IP)
	if err != nil {
		return nil, err
	}
	if details == nil {
		// we do not yet have this in the DB, do a Spamhaus lookup
		details, err = newIPLookup(IP)
		if err != nil {
			return nil, err
		}
	}

	return details, nil
}

// getFromDB returns the model if it exists and nil otherwise.
func (r *Resolver) getFromDB(IP string) (*model.IPDetails, error) {
	var err error
	r.client, err = sqlite.NewClient(FILENAME)
	if err != nil {
		return nil, err
	}

	defer r.client.Close()

	details, err := r.client.GetIPDetailByAddress(IP)
	if err != nil {
		// we do not yet have this in the DB
		if strings.Contains(err.Error(), "no rows in result") {
			return nil, nil
		} else {
			// all other errors besides "not yet in DB"
			return nil, err
		}

	}

	return details, nil
}

// storeToDB stores the model provided
func (r *Resolver) storeToDB(toStore model.IPDetails) error {
	var err error
	r.client, err = sqlite.NewClient(FILENAME)
	if err != nil {
		return err
	}

	defer r.client.Close()

	stored, err := r.client.AddIPDetails(toStore)
	if err != nil {
		return err
	}
	if stored.UUID != toStore.UUID {
		return fmt.Errorf("value stored in database does not match value expected")
	}

	return nil
}

// newIPLookup should be used when an IP is unknown to the system
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
