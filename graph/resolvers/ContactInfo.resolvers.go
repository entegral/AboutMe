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

func (r *mutationResolver) UpdateContactInfo(ctx context.Context, key string, input model.ContactInfoInput) (*model.ContactInfo, error) {
	if !auth.CheckKey(key) {
		return nil, errors.New("invalid key")
	}
	info, err := input.Update()
	if e.Warn("error occurred in UpdateInfo", err) {
		return nil, err
	}
	return info, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
