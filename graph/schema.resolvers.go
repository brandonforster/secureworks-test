package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/brandonforster/resolver/graph/generated"
	"github.com/brandonforster/resolver/graph/model"
)

func (r *mutationResolver) Enqueue(ctx context.Context, ip []*model.NewIP) ([]*model.IPAddress, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Ips(ctx context.Context) ([]*model.IPAddress, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
