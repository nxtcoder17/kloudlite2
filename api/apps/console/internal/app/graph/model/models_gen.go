// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"kloudlite.io/pkg/repos"
)

type Cluster struct {
	ID         repos.ID  `json:"id"`
	Name       string    `json:"name"`
	Provider   string    `json:"provider"`
	Region     string    `json:"region"`
	IP         *string   `json:"ip"`
	Devices    []*Device `json:"devices"`
	NodesCount int       `json:"nodesCount"`
	Status     string    `json:"status"`
}

func (Cluster) IsEntity() {}

type Device struct {
	ID            repos.ID `json:"id"`
	User          *User    `json:"user"`
	Name          string   `json:"name"`
	Cluster       *Cluster `json:"cluster"`
	Configuration string   `json:"configuration"`
	DeviceIndex int `json:"device_index"`
	Ip string `json:"ip"`
}

func (Device) IsEntity() {}

type User struct {
	ID      repos.ID  `json:"id"`
	Devices []*Device `json:"devices"`
}

func (User) IsEntity() {}
