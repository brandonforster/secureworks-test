package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/brandonforster/resolver/graph/generated"
	"github.com/brandonforster/resolver/graph/model"
)

// TODO: need to do... *waves hands*
func (r *mutationResolver) Enqueue(ctx context.Context, ip []string) (*bool, error) {
	retval := true
	for _, address := range ip {
		retval = r.Store(address) && retval
	}

	return &retval, nil
}

func (r *queryResolver) GetIPDetails(ctx context.Context, ip string) (*model.IPDetails, error) {
	return r.Get(ip)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
