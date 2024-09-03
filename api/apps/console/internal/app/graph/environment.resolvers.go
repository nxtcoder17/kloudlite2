package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"github.com/kloudlite/api/pkg/errors"
	"time"

	"github.com/kloudlite/api/apps/console/internal/app/graph/generated"
	"github.com/kloudlite/api/apps/console/internal/app/graph/model"
	"github.com/kloudlite/api/apps/console/internal/entities"
	fn "github.com/kloudlite/api/pkg/functions"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreationTime is the resolver for the creationTime field.
func (r *environmentResolver) CreationTime(ctx context.Context, obj *entities.Environment) (string, error) {
	if obj == nil {
		return "", errNilEnvironment
	}
	return obj.BaseEntity.CreationTime.Format(time.RFC3339), nil
}

// Spec is the resolver for the spec field.
func (r *environmentResolver) Spec(ctx context.Context, obj *entities.Environment) (*model.GithubComKloudliteOperatorApisCrdsV1EnvironmentSpec, error) {
	if obj == nil {
		return nil, errNilEnvironment
	}
	m := &model.GithubComKloudliteOperatorApisCrdsV1EnvironmentSpec{}
	if err := fn.JsonConversion(obj.Spec, &m); err != nil {
		return nil, errors.NewE(err)
	}

	if m.Suspend == nil {
		m.Suspend = fn.New(false)
	}

	return m, nil
}

// UpdateTime is the resolver for the updateTime field.
func (r *environmentResolver) UpdateTime(ctx context.Context, obj *entities.Environment) (string, error) {
	if obj == nil {
		return "", errNilEnvironment
	}
	return obj.BaseEntity.UpdateTime.Format(time.RFC3339), nil
}

// Metadata is the resolver for the metadata field.
func (r *environmentInResolver) Metadata(ctx context.Context, obj *entities.Environment, data *v1.ObjectMeta) error {
	if obj == nil {
		return errNilEnvironment
	}
	if data != nil {
		obj.ObjectMeta = *data
	}
	return nil
}

// Spec is the resolver for the spec field.
func (r *environmentInResolver) Spec(ctx context.Context, obj *entities.Environment, data *model.GithubComKloudliteOperatorApisCrdsV1EnvironmentSpecIn) error {
	if obj == nil {
		return errNilEnvironment
	}
	return fn.JsonConversion(data, &obj.Spec)
}

// Environment returns generated.EnvironmentResolver implementation.
func (r *Resolver) Environment() generated.EnvironmentResolver { return &environmentResolver{r} }

// EnvironmentIn returns generated.EnvironmentInResolver implementation.
func (r *Resolver) EnvironmentIn() generated.EnvironmentInResolver { return &environmentInResolver{r} }

type environmentResolver struct{ *Resolver }
type environmentInResolver struct{ *Resolver }
