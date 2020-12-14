package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/brandonforster/resolver/graph/generated"
	"github.com/brandonforster/resolver/graph/model"
)

func (r *mutationResolver) Enqueue(ctx context.Context, ip []string) ([]*model.IPDetails, error) {
		if !isAuthorized(ctx) {
		return nil, fmt.Errorf("access denied")
	}

	outputModels := make([]*model.IPDetails, len(ip))

	for i, address := range ip {
		modelChan, errChan := r.Queue(address)

		err := <- errChan
		if err != nil {
			return nil, err
		}
		outputModels[i] = <-modelChan
	}

	return outputModels, nil
}

func (r *queryResolver) GetIPDetails(ctx context.Context, ip string) (*model.IPDetails, error) {
	if !isAuthorized(ctx) {
		return nil, fmt.Errorf("access denied")
	}

	return r.Get(ip)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
