package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/brandonforster/resolver/graph/generated"
	"github.com/brandonforster/resolver/graph/model"
	"github.com/brandonforster/resolver/internal/spamhaus"
	"github.com/brandonforster/resolver/internal/sqlite"
)

const FILENAME = `resolver.db`

// Enqueue is designed to kick off a background job to do the DNS lookup and store it in the database for each IP
// passed in for future lookups. If the lookup has already happened, this will queue it up again and update the
// record in the database. It returns an array of records stored in the database if successful, an error otherwise.
//
// ctx is the context of the running app, used in this method to check authentication
// ip is a list of IPs to check and ultimately store
//
// Returns a slice of IPDetails if the query executed successfully; an error otherwise.
func (r *mutationResolver) Enqueue(ctx context.Context, ip []string) ([]*model.IPDetails, error) {
	if !isAuthorized(ctx) {
		return nil, fmt.Errorf("access denied")
	}

	var err error
	r.DBClient, err = sqlite.NewClient(FILENAME)
	if err != nil {
		return nil, err
	}

	defer r.DBClient.Close()

	r.BlocklistClient = spamhaus.Client{}

	outputModels := make([]*model.IPDetails, len(ip))
	modelChan := make(chan *model.IPDetails, len(ip))
	errChan := make(chan error, len(ip))
	defer close(modelChan)
	defer close(errChan)

	for _, address := range ip {
		go r.Queue(address, modelChan, errChan)
	}

	for i := range ip {
		select {
		case err := <-errChan:
			return nil, err
		case result := <-modelChan:
			outputModels[i] = result
		}
	}

	return outputModels, nil
}

// GetIPDetails will look up the IP provided in the database, and if it doesn't exist it will do a one-off lookup.
// This function will never write to the database.
//
// ctx is the context of the running app, used in this method to check authentication
// ip is the IP to check
//
// Returns the IPDetails if the query executed successfully; an error otherwise.
func (r *queryResolver) GetIPDetails(ctx context.Context, ip string) (*model.IPDetails, error) {
	if !isAuthorized(ctx) {
		return nil, fmt.Errorf("access denied")
	}

	var err error
	r.DBClient, err = sqlite.NewClient(FILENAME)
	if err != nil {
		return nil, err
	}

	defer r.DBClient.Close()

	r.BlocklistClient = spamhaus.Client{}

	return r.Get(ip)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
