package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/entegral/aboutme/auth"
	e "github.com/entegral/aboutme/errors"
	"github.com/entegral/aboutme/graph/generated"
	"github.com/entegral/aboutme/graph/model"
)

func (r *mutationResolver) UpdateInfo(ctx context.Context, key string, info *model.UpdateMeInput) (*model.Me, error) {
	if !auth.CheckKey(key) {
		return nil, errors.New("invalid key")
	}
	user, err := info.Update()
	if e.Warn("error occured in UpdateInfo", err) {
		return nil, err
	}
	return user, nil
}

func (r *queryResolver) AboutMe(ctx context.Context) (*model.Me, error) {
	input := model.GetMeInput{
		FirstName: "robby",
		LastName:  "bruce",
	}
	return r.About(ctx, input)
}

func (r *queryResolver) About(ctx context.Context, input model.GetMeInput) (*model.Me, error) {
	val, err := input.GetMe()
	errMsg := fmt.Sprintf("Error fetching %s %s's general info", input.FirstName, input.LastName)
	if e.Warn(errMsg, err) {
		return nil, err
	}
	return val, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
