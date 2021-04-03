package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/entegral/aboutme/graph/generated"
	"github.com/entegral/aboutme/graph/model"
)

func (r *mutationResolver) AddEducation(ctx context.Context, key string, info *model.EducationInput) (*model.Education, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveEducation(ctx context.Context, firstName string, lastName string) (*model.Education, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
