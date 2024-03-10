package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.28

import (
	"context"
	"github.com/kloudlite/api/pkg/errors"
	"time"

	"github.com/kloudlite/api/apps/infra/internal/app/graph/generated"
	"github.com/kloudlite/api/apps/infra/internal/app/graph/model"
	"github.com/kloudlite/api/apps/infra/internal/entities"
	fn "github.com/kloudlite/api/pkg/functions"
	ct "github.com/kloudlite/operator/apis/common-types"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Aws is the resolver for the aws field.
func (r *cloudProviderSecretResolver) Aws(ctx context.Context, obj *entities.CloudProviderSecret) (*model.GithubComKloudliteAPIAppsInfraInternalEntitiesAWSSecretCredentials, error) {
	if obj == nil || obj.CreationTime.IsZero() {
		return nil, errors.Newf("CloudProviderSecret object is nil")
	}

	return fn.JsonConvertP[model.GithubComKloudliteAPIAppsInfraInternalEntitiesAWSSecretCredentials](obj.AWS)
}

// CloudProviderName is the resolver for the cloudProviderName field.
func (r *cloudProviderSecretResolver) CloudProviderName(ctx context.Context, obj *entities.CloudProviderSecret) (model.GithubComKloudliteOperatorApisCommonTypesCloudProvider, error) {
	return model.GithubComKloudliteOperatorApisCommonTypesCloudProvider(obj.CloudProviderName), nil
}

// CreationTime is the resolver for the creationTime field.
func (r *cloudProviderSecretResolver) CreationTime(ctx context.Context, obj *entities.CloudProviderSecret) (string, error) {
	if obj == nil || obj.CreationTime.IsZero() {
		return "", errors.Newf("CloudProviderSecret object is nil")
	}
	return obj.CreationTime.Format(time.RFC3339), nil
}

// ID is the resolver for the id field.
func (r *cloudProviderSecretResolver) ID(ctx context.Context, obj *entities.CloudProviderSecret) (string, error) {
	if obj == nil {
		return "", errors.Newf("CloudProviderSecret object is nil")
	}

	return string(obj.Id), nil
}

// UpdateTime is the resolver for the updateTime field.
func (r *cloudProviderSecretResolver) UpdateTime(ctx context.Context, obj *entities.CloudProviderSecret) (string, error) {
	if obj == nil || obj.UpdateTime.IsZero() {
		return "", errors.Newf("CloudProviderSecret object is nil")
	}

	return obj.UpdateTime.Format(time.RFC3339), nil
}

// Aws is the resolver for the aws field.
func (r *cloudProviderSecretInResolver) Aws(ctx context.Context, obj *entities.CloudProviderSecret, data *model.GithubComKloudliteAPIAppsInfraInternalEntitiesAWSSecretCredentialsIn) error {
	return fn.JsonConversion(data, &obj.AWS)
}

// CloudProviderName is the resolver for the cloudProviderName field.
func (r *cloudProviderSecretInResolver) CloudProviderName(ctx context.Context, obj *entities.CloudProviderSecret, data model.GithubComKloudliteOperatorApisCommonTypesCloudProvider) error {
	if !data.IsValid() {
		return errors.Newf("invalid cloud provider name")
	}
	// obj.CloudProviderName = ct.CloudProvider(parser.RestoreSanitizedPackagePath(data.String()))
	obj.CloudProviderName = ct.CloudProvider(data)
	return nil
}

// Metadata is the resolver for the metadata field.
func (r *cloudProviderSecretInResolver) Metadata(ctx context.Context, obj *entities.CloudProviderSecret, data *v1.ObjectMeta) error {
	return fn.JsonConversion(data, &obj.ObjectMeta)
}

// CloudProviderSecret returns generated.CloudProviderSecretResolver implementation.
func (r *Resolver) CloudProviderSecret() generated.CloudProviderSecretResolver {
	return &cloudProviderSecretResolver{r}
}

// CloudProviderSecretIn returns generated.CloudProviderSecretInResolver implementation.
func (r *Resolver) CloudProviderSecretIn() generated.CloudProviderSecretInResolver {
	return &cloudProviderSecretInResolver{r}
}

type cloudProviderSecretResolver struct{ *Resolver }
type cloudProviderSecretInResolver struct{ *Resolver }
