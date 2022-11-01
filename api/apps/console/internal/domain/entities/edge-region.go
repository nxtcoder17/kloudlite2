package entities

import (
	"kloudlite.io/pkg/repos"
)

type NodePool struct {
	Name   string   `json:"name"`
	Config string   `json:"config"`
	Min    int      `json:"min"`
	Max    int      `json:"max"`
	Nodes  []string `bson:"nodes"`
}

type EdgeRegion struct {
	repos.BaseEntity `bson:",inline"`
	Name             string     `bson:"name"`
	ProviderId       repos.ID   `bson:"provider_id"`
	Region           string     `bson:"region"`
	Pools            []NodePool `bson:"pools"`
}

type CloudProvider struct {
	repos.BaseEntity `bson:",inline"`
	Name             string            `bson:"name"`
	AccountId        *repos.ID         `json:"account_id,omitempty" bson:"account_id"`
	Provider         string            `json:"provider" bson:"provider"`
	Credentials      map[string]string `json:"credentials" bson:"credentials"`
	Status           string            `json:"status" bson:"status"`
}

var CloudProviderIndexes = []repos.IndexField{
	{
		Field: []repos.IndexKey{
			{Key: "id", Value: repos.IndexAsc},
		},
		Unique: true,
	},
	{
		Field: []repos.IndexKey{
			{Key: "account_id", Value: repos.IndexAsc},
		},
	},
	{
		Field: []repos.IndexKey{
			{Key: "provider", Value: repos.IndexAsc},
		},
	},
}

var EdgeRegionIndexes = []repos.IndexField{
	{
		Field: []repos.IndexKey{
			{Key: "id", Value: repos.IndexAsc},
		},
		Unique: true,
	},
	{
		Field: []repos.IndexKey{
			{Key: "region", Value: repos.IndexAsc},
		},
		Unique: true,
	},
	{
		Field: []repos.IndexKey{
			{Key: "provider", Value: repos.IndexAsc},
		},
	},
}
