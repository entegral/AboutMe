package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/entegral/aboutme/auth"
	e "github.com/entegral/aboutme/errors"
	"github.com/entegral/aboutme/graph/generated"
	"github.com/entegral/aboutme/graph/model"
)

func (r *mutationResolver) UpdateInfo(ctx context.Context, key string, info *model.UpdateMe) (*model.Me, error) {
	if !auth.CheckKey(key) {
		return nil, errors.New("invalid key")
	}
	var user model.Me
	user.FirstName = info.FirstName
	user.LastName = info.LastName
	user.Title = info.Title
	user.Location = info.Location
	user.Save()
	return &user, nil
}

func (r *queryResolver) AboutMe(ctx context.Context) (*model.Me, error) {
	return r.About(ctx, "robby", "bruce")
}

func (r *queryResolver) About(ctx context.Context, firstName string, lastName string) (*model.Me, error) {
	var input model.CompositeKey
	val, err := input.GetMe(lastName, firstName)
	if e.Warn("Error fetching RobbyBruce aboutme data", err) {
		return nil, err
	}
	return val, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
