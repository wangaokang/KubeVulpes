package types

import "kubevulpes/pkg/db/model"

type RBACPolicy struct {
	UserName   string           `json:"username,omitempty"`
	GroupName  string           `json:"groupname,omitempty"`
	ObjectType model.ObjectType `json:"resource_type,omitempty"`
	StringID   string           `json:"sid,omitempty"`
	Operation  model.Operation  `json:"operation,omitempty"`
}
