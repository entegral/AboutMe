package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/entegral/aboutme/graph/model"
)

func (r *mutationResolver) UpdateProject(ctx context.Context, key string, info *model.ProjectInput) (*model.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveProject(ctx context.Context, key string, company *string, title *string) (*model.Project, error) {
	panic(fmt.Errorf("not implemented"))
}
