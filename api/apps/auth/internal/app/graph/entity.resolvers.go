package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.28

import (
	"context"

	"github.com/kloudlite/api/apps/auth/internal/app/graph/generated"
	"github.com/kloudlite/api/apps/auth/internal/app/graph/model"
	"github.com/kloudlite/api/pkg/repos"
)

// FindUserByID is the resolver for the findUserByID field.
func (r *entityResolver) FindUserByID(ctx context.Context, id repos.ID) (*model.User, error) {
	u, err := r.d.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return userModelFromEntity(u), nil
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
