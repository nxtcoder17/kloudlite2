package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.28

import (
	"context"
	"time"

	"github.com/kloudlite/api/pkg/errors"

	"github.com/kloudlite/api/apps/infra/internal/app/graph/generated"
	"github.com/kloudlite/api/apps/infra/internal/app/graph/model"
	"github.com/kloudlite/api/apps/infra/internal/entities"
	fn "github.com/kloudlite/api/pkg/functions"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreationTime is the resolver for the creationTime field.
func (r *clusterResolver) CreationTime(ctx context.Context, obj *entities.Cluster) (string, error) {
	if obj == nil {
		return "", errors.Newf("cluster obj is nil")
	}
	return obj.CreationTime.Format(time.RFC3339), nil
}

// ID is the resolver for the id field.
func (r *clusterResolver) ID(ctx context.Context, obj *entities.Cluster) (string, error) {
	if obj == nil {
		return "", errors.Newf("cluster obj is nil")
	}
	return string(obj.Id), nil
}

// Spec is the resolver for the spec field.
func (r *clusterResolver) Spec(ctx context.Context, obj *entities.Cluster) (*model.GithubComKloudliteOperatorApisClustersV1ClusterSpec, error) {
	if obj == nil {
		return nil, errors.Newf("cluster is nil")
	}

	var spec model.GithubComKloudliteOperatorApisClustersV1ClusterSpec
	if err := fn.JsonConversion(&obj.Spec, &spec); err != nil {
		return nil, errors.NewE(err)
	}
	return &spec, nil
}

// UpdateTime is the resolver for the updateTime field.
func (r *clusterResolver) UpdateTime(ctx context.Context, obj *entities.Cluster) (string, error) {
	if obj == nil {
		return "", errors.Newf("cluster is nil")
	}
	return obj.UpdateTime.Format(time.RFC3339), nil
}

// Metadata is the resolver for the metadata field.
func (r *clusterInResolver) Metadata(ctx context.Context, obj *entities.Cluster, data *v1.ObjectMeta) error {
	if obj == nil {
		return errors.Newf("cluster is nil")
	}
	return fn.JsonConversion(data, &obj.ObjectMeta)
}

// Spec is the resolver for the spec field.
func (r *clusterInResolver) Spec(ctx context.Context, obj *entities.Cluster, data *model.GithubComKloudliteOperatorApisClustersV1ClusterSpecIn) error {
	if obj == nil {
		return errors.Newf("cluster is nil")
	}
	return fn.JsonConversion(data, &obj.Spec)
}

// Cluster returns generated.ClusterResolver implementation.
func (r *Resolver) Cluster() generated.ClusterResolver { return &clusterResolver{r} }

// ClusterIn returns generated.ClusterInResolver implementation.
func (r *Resolver) ClusterIn() generated.ClusterInResolver { return &clusterInResolver{r} }

type (
	clusterResolver   struct{ *Resolver }
	clusterInResolver struct{ *Resolver }
)
