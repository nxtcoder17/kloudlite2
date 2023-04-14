package types

import (
	rApi "github.com/kloudlite/operator/pkg/operator"
)

type StatusUpdate struct {
	AuthToken   string         `json:"authToken"`
	AccountName string         `json:"accountName"`
	ClusterName string         `json:"clusterName"`
	Object      map[string]any `json:"object"`
	Status      rApi.Status    `json:"status"`
}
