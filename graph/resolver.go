package graph

//go:generate go run github.com/99designs/gqlgen
import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/brandonforster/resolver/graph/model"
	"github.com/brandonforster/resolver/internal/interfaces"
)

type Resolver struct {
	DBClient        interfaces.DBClient
	BlocklistClient interfaces.BlocklistClient
}

// GetAndStore is the backing of the Enqueue mutation.
// It is designed to get details of an address in the database if they exist, create them if they don't,
// check the blocklist in either case and store the record of the check in the database.
// It returns a pointer to an IPDetails if the check and database store were successful; an error otherwise.
//
// IP is the address that should be checked and stored into memory
//
// Returns the IPDetails if the query executed successfully; an error otherwise.
func (r *Resolver) GetAndStore(IP string) (*model.IPDetails, error) {
	// does this IP exist in the DB?
	details, err := r.getFromDB(IP)
	if err != nil {
		return nil, err
	}

	if details == nil {
		// we do not yet have this in the DB, do a lookup
		details, err = r.newIPLookup(IP)
		if err != nil {
			return nil, err
		}

		err = r.addToDB(*details)
		if err != nil {
			return nil, err
		}
	} else {
		// this exists, update the response codes and LastUpdated
		details, err = r.updateToDB(*details)
		if err != nil {
			return nil, err
		}
	}

	return details, nil
}

// Queue is intended to be used asychronously. It runs GetAndStore and updates the provided channels.
//
// IP is the IP to be checked and stored
// modelChan is the channel to receive full data back from the database
// errChan is the channel to receive any errors
func (r *Resolver) Queue(ip string, modelChan chan *model.IPDetails, errChan chan error) {
	details, err := r.GetAndStore(ip)
	if err != nil {
		errChan <- err
	} else {
		modelChan <- details
	}
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
		// we do not yet have this in the DB, do a lookup
		details, err = r.newIPLookup(IP)
		if err != nil {
			return nil, err
		}
	}

	return details, nil
}

// getFromDB returns the model if it exists and nil otherwise.
func (r *Resolver) getFromDB(IP string) (*model.IPDetails, error) {
	details, err := r.DBClient.GetIPDetailByAddress(IP)
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

// addToDB stores the model provided
func (r *Resolver) addToDB(toStore model.IPDetails) error {
	stored, err := r.DBClient.AddIPDetails(toStore)
	if err != nil {
		return err
	}
	if stored.UUID != toStore.UUID {
		return fmt.Errorf("value stored in database does not match value expected")
	}

	return nil
}

// updateToDB updates the model provided and stores it to the DB
func (r *Resolver) updateToDB(toUpdate model.IPDetails) (*model.IPDetails, error) {
	responseCode, err := r.BlocklistClient.Lookup(toUpdate.IPAddress)
	if err != nil {
		return nil, err
	}

	updated := toUpdate
	updated.ResponseCode = responseCode
	updated.UpdatedAt = time.Now()

	result, err := r.DBClient.UpdateIPDetails(updated)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// newIPLookup should be used when an IP is unknown to the system
func (r *Resolver) newIPLookup(IP string) (*model.IPDetails, error) {
	responseCode, err := r.BlocklistClient.Lookup(IP)
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

// isAuthorized checks the context for the auth flag
func isAuthorized(ctx context.Context) bool {
	isAuth, _ := ctx.Value("isAuth").(bool)

	return isAuth
}
