package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.28

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"kloudlite.io/apps/console/internal/app/graph/generated"
	"kloudlite.io/apps/console/internal/app/graph/model"
	"kloudlite.io/apps/console/internal/domain/entities"
	fn "kloudlite.io/pkg/functions"
)

// CreationTime is the resolver for the creationTime field.
func (r *appResolver) CreationTime(ctx context.Context, obj *entities.App) (string, error) {
	return obj.BaseEntity.CreationTime.Format(time.RFC3339), nil
}

// ID is the resolver for the id field.
func (r *appResolver) ID(ctx context.Context, obj *entities.App) (string, error) {
	return string(obj.Id), nil
}

// Spec is the resolver for the spec field.
func (r *appResolver) Spec(ctx context.Context, obj *entities.App) (*model.GithubComKloudliteOperatorApisCrdsV1AppSpec, error) {
	m := model.GithubComKloudliteOperatorApisCrdsV1AppSpec{}
	if err := fn.JsonConversion(obj.Spec, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// UpdateTime is the resolver for the updateTime field.
func (r *appResolver) UpdateTime(ctx context.Context, obj *entities.App) (string, error) {
	if obj == nil {
		return "", fmt.Errorf("obj is nil")
	}

	return obj.BaseEntity.UpdateTime.Format(time.RFC3339), nil
}

// Metadata is the resolver for the metadata field.
func (r *appInResolver) Metadata(ctx context.Context, obj *entities.App, data *v1.ObjectMeta) error {
	obj.ObjectMeta = *data
	return nil
}

// Spec is the resolver for the spec field.
func (r *appInResolver) Spec(ctx context.Context, obj *entities.App, data *model.GithubComKloudliteOperatorApisCrdsV1AppSpecIn) error {
	return fn.JsonConversion(data, &obj.Spec)
}

// App returns generated.AppResolver implementation.
func (r *Resolver) App() generated.AppResolver { return &appResolver{r} }

// AppIn returns generated.AppInResolver implementation.
func (r *Resolver) AppIn() generated.AppInResolver { return &appInResolver{r} }

type appResolver struct{ *Resolver }
type appInResolver struct{ *Resolver }
