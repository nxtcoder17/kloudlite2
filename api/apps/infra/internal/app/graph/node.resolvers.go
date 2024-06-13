package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"github.com/kloudlite/api/pkg/errors"
	"time"

	"github.com/kloudlite/api/apps/infra/internal/app/graph/generated"
	"github.com/kloudlite/api/apps/infra/internal/app/graph/model"
	"github.com/kloudlite/api/apps/infra/internal/entities"
	fn "github.com/kloudlite/api/pkg/functions"
	"github.com/kloudlite/api/pkg/repos"
)

// CreationTime is the resolver for the creationTime field.
func (r *nodeResolver) CreationTime(ctx context.Context, obj *entities.Node) (string, error) {
	if obj == nil || obj.CreationTime.IsZero() {
		return "", errors.Newf("node is nil")
	}
	return obj.CreationTime.Format(time.RFC3339), nil
}

// ID is the resolver for the id field.
func (r *nodeResolver) ID(ctx context.Context, obj *entities.Node) (repos.ID, error) {
	if obj == nil {
		return "", errors.Newf("node is nil")
	}
	return obj.Id, nil
}

// Spec is the resolver for the spec field.
func (r *nodeResolver) Spec(ctx context.Context, obj *entities.Node) (*model.GithubComKloudliteOperatorApisClustersV1NodeSpec, error) {
	var m model.GithubComKloudliteOperatorApisClustersV1NodeSpec
	if err := fn.JsonConversion(obj.Spec, &m); err != nil {
		return nil, errors.NewE(err)
	}
	return &m, nil
}

// UpdateTime is the resolver for the updateTime field.
func (r *nodeResolver) UpdateTime(ctx context.Context, obj *entities.Node) (string, error) {
	if obj == nil || obj.UpdateTime.IsZero() {
		return "", errors.Newf("node is nil")
	}
	return obj.UpdateTime.Format(time.RFC3339), nil
}

// Node returns generated.NodeResolver implementation.
func (r *Resolver) Node() generated.NodeResolver { return &nodeResolver{r} }

type nodeResolver struct{ *Resolver }
