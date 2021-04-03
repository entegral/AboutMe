package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/entegral/aboutme/auth"
	e "github.com/entegral/aboutme/errors"
	"github.com/entegral/aboutme/graph/model"
)

func (r *mutationResolver) UpdateExperience(ctx context.Context, key string, input model.ExperienceInput) (*model.Experience, error) {
	if !auth.CheckKey(key) {
		return nil, errors.New("invalid key")
	}
	ex, err := input.Update()
	if e.Warn("error occurred in UpdateInfo", err) {
		return nil, err
	}
	return ex, nil
}

func (r *mutationResolver) RemoveExperience(ctx context.Context, key string, input model.ExperienceKeyInput) (*model.Experience, error) {
	if !auth.CheckKey(key) {
		return nil, errors.New("invalid key")
	}
	return &model.Experience{}, nil
}
