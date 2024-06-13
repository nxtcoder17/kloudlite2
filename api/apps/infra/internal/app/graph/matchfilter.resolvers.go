package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/kloudlite/api/apps/infra/internal/app/graph/generated"
	"github.com/kloudlite/api/apps/infra/internal/app/graph/model"
	"github.com/kloudlite/api/pkg/repos"
)

// MatchType is the resolver for the matchType field.
func (r *matchFilterResolver) MatchType(ctx context.Context, obj *repos.MatchFilter) (model.GithubComKloudliteAPIPkgReposMatchType, error) {
	panic(fmt.Errorf("not implemented: MatchType - matchType"))
}

// MatchType is the resolver for the matchType field.
func (r *matchFilterInResolver) MatchType(ctx context.Context, obj *repos.MatchFilter, data model.GithubComKloudliteAPIPkgReposMatchType) error {
	panic(fmt.Errorf("not implemented: MatchType - matchType"))
}

// MatchFilter returns generated.MatchFilterResolver implementation.
func (r *Resolver) MatchFilter() generated.MatchFilterResolver { return &matchFilterResolver{r} }

// MatchFilterIn returns generated.MatchFilterInResolver implementation.
func (r *Resolver) MatchFilterIn() generated.MatchFilterInResolver { return &matchFilterInResolver{r} }

type matchFilterResolver struct{ *Resolver }
type matchFilterInResolver struct{ *Resolver }
