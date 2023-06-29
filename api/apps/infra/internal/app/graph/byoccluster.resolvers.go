package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.28

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"kloudlite.io/apps/infra/internal/app/graph/generated"
	"kloudlite.io/apps/infra/internal/app/graph/model"
	"kloudlite.io/apps/infra/internal/domain/entities"
	fn "kloudlite.io/pkg/functions"
)

// CreationTime is the resolver for the creationTime field.
func (r *bYOCClusterResolver) CreationTime(ctx context.Context, obj *entities.BYOCCluster) (string, error) {
	if obj == nil || obj.CreationTime.IsZero() {
		return "", fmt.Errorf("byocCluster/creation-time is nil")
	}
	return obj.CreationTime.Format(time.RFC3339), nil
}

// HelmStatus is the resolver for the helmStatus field.
func (r *bYOCClusterResolver) HelmStatus(ctx context.Context, obj *entities.BYOCCluster) (map[string]interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("byocCluster is nil")
	}
	var m map[string]any
	if err := fn.JsonConversion(obj.HelmStatus, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ID is the resolver for the id field.
func (r *bYOCClusterResolver) ID(ctx context.Context, obj *entities.BYOCCluster) (string, error) {
	if obj == nil {
		return "", fmt.Errorf("byocCluster is nil")
	}
	return string(obj.Id), nil
}

// Spec is the resolver for the spec field.
func (r *bYOCClusterResolver) Spec(ctx context.Context, obj *entities.BYOCCluster) (*model.GithubComKloudliteOperatorApisClustersV1BYOCSpec, error) {
	var m model.GithubComKloudliteOperatorApisClustersV1BYOCSpec
	if err := fn.JsonConversion(obj.Spec, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// UpdateTime is the resolver for the updateTime field.
func (r *bYOCClusterResolver) UpdateTime(ctx context.Context, obj *entities.BYOCCluster) (string, error) {
	if obj == nil || obj.UpdateTime.IsZero() {
		return "", fmt.Errorf("byocCluster/update-time is nil")
	}
	return obj.UpdateTime.Format(time.RFC3339), nil
}

// Metadata is the resolver for the metadata field.
func (r *bYOCClusterInResolver) Metadata(ctx context.Context, obj *entities.BYOCCluster, data *v1.ObjectMeta) error {
	if obj == nil {
		return nil
	}

	return fn.JsonConversion(data, &obj.ObjectMeta)
}

// Spec is the resolver for the spec field.
func (r *bYOCClusterInResolver) Spec(ctx context.Context, obj *entities.BYOCCluster, data *model.GithubComKloudliteOperatorApisClustersV1BYOCSpecIn) error {
	if obj == nil {
		return nil
	}

	return fn.JsonConversion(data, &obj.Spec)
}

// BYOCCluster returns generated.BYOCClusterResolver implementation.
func (r *Resolver) BYOCCluster() generated.BYOCClusterResolver { return &bYOCClusterResolver{r} }

// BYOCClusterIn returns generated.BYOCClusterInResolver implementation.
func (r *Resolver) BYOCClusterIn() generated.BYOCClusterInResolver { return &bYOCClusterInResolver{r} }

type bYOCClusterResolver struct{ *Resolver }
type bYOCClusterInResolver struct{ *Resolver }
