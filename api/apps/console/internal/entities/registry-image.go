package entities

import (
	fc "github.com/kloudlite/api/apps/console/internal/entities/field-constants"
	"github.com/kloudlite/api/common/fields"
	"github.com/kloudlite/api/pkg/repos"
)

type RegistryImage struct {
	repos.BaseEntity `json:",inline" graphql:"noinput"`
	AccountName      string         `json:"accountName"`
	ImageName        string         `json:"imageName"`
	ImageTag         string         `json:"imageTag"`
	Meta             map[string]any `json:"meta"`
	MetaData         string         `json:"metadata" graphql:"ignore"`
}

type RegistryImageURL struct {
	URL        []string `json:"url"`
	URLExample []string `json:"urlExample"`

	KlWebhookAuthToken string `json:"klWebhookAuthToken"`

	ScriptURL        []string `json:"scriptUrl"`
	ScriptURLExample []string `json:"scriptUrlExample"`
}

var RegistryImageIndexes = []repos.IndexField{
	{
		Field: []repos.IndexKey{
			{Key: fields.Id, Value: repos.IndexAsc},
		},
		Unique: true,
	},
	{
		Field: []repos.IndexKey{
			{Key: fields.AccountName, Value: repos.IndexAsc},
			{Key: fc.RegistryImageImageName, Value: repos.IndexAsc},
			{Key: fc.RegistryImageImageTag, Value: repos.IndexAsc},
		},
		Unique: true,
	},
	{
		Field: []repos.IndexKey{
			{Key: fields.AccountName, Value: repos.IndexAsc},
			{Key: fc.Metadata, Value: repos.IndexAsc, IsText: true},
			{Key: fc.RegistryImageImageName, Value: repos.IndexAsc, IsText: true},
		},
	},
}
