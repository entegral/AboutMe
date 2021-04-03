package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/entegral/aboutme/graph/model"
)

func (r *mutationResolver) AddGoSkill(ctx context.Context, key string, info *model.GoSkillsInput) (*model.GoSkills, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddJSSkill(ctx context.Context, key string, info *model.JSSkillsInput) (*model.JSSkills, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddPythonSkill(ctx context.Context, key string, info *model.PythonSkillsInput) (*model.PythonSkills, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveGoSkill(ctx context.Context, key string, firstName *string, lastName *string) (*model.GoSkills, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveJSSkill(ctx context.Context, key string, firstName *string, lastName *string) (*model.JSSkills, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemovePythonSkill(ctx context.Context, key string, firstName *string, lastName *string) (*model.PythonSkills, error) {
	panic(fmt.Errorf("not implemented"))
}
