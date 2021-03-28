package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/entegral/aboutme/errors"
	"github.com/entegral/aboutme/graph/generated"
	"github.com/entegral/aboutme/graph/model"
)

func (r *queryResolver) AboutMe(ctx context.Context) (*model.Me, error) {
	var input model.CompositeKey
	val, err := input.GetMe("bruce", "robby")
	if errors.Warn("Error fetching RobbyBruce aboutme data", err) {
		return nil, err
	}
	return val, nil}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
