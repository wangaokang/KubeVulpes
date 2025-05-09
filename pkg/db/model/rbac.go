/*
Copyright 2024 The Vuples Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package model

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"kubevulpes/pkg/db/model/base"
	utils "kubevulpes/pkg/util"
)

const (
	AdminGroup           = "root"
	ReadOnlyGroup        = "readonly"
	ReadWriteGroup       = "readwrite"
	ReadWriteUpdateGroup = "readwriteupdate"
	SidAll               = "*"
)

type Operation string

const (
	OpRead   Operation = "read"
	OpCreate Operation = "create"
	OpUpdate Operation = "update"
	OpDelete Operation = "delete"
	OpAll    Operation = "*"
)

func (o Operation) String() string {
	return string(o)
}

func buildOperation(ops ...Operation) Operation {
	if len(ops) == 0 {
		return ""
	}

	s := `\b(`
	for i, op := range ops {
		if op == "*" {
			return "*"
		}

		s += op.String()
		if i < len(ops)-1 {
			s += "|"
			continue
		}
		s += `)\b`
	}
	return Operation(s)
}

var OperationMap = map[Operation]struct{}{
	OpRead:   {},
	OpCreate: {},
	OpUpdate: {},
	OpDelete: {},
	OpAll:    {},
}

type ObjectType string

const (
	ObjectUser    ObjectType = "users"
	ObjectCluster ObjectType = "clusters"
	ObjectAuth    ObjectType = "auth"
	ObjectAll     ObjectType = "*"
)

func (o ObjectType) String() string {
	return string(o)
}

var ObjectTypeMap = map[ObjectType]struct{}{
	ObjectUser:    {},
	ObjectCluster: {},
	//ObjectAuth:    {},
	ObjectAll: {},
}

// TODO:
type RBACInterface interface{}

// Casbin RBAC model
// ref: https://github.com/casbin/casbin/blob/master/examples/rbac_with_domains_model.conf
const RBACModel = `
[request_definition]
r = sub, obj, id, op

[policy_definition]
p = sub, obj, id, op

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && keyMatch2(r.id, p.id) && keyMatch2(r.op, p.op)`

// TODO:
type CasbinRBACImpl struct{}

type Policy interface {
	Raw() []string
}

// UserPolicy is a RBAC policy for user.
// e.g. ["foo", "clusters", "*", "read"]
type UserPolicy [4]string

// NewUserPolicy returns a policy slice for user.
// e.g. ["foo", "clusters", "*", "read"]: foo is a user name
func NewUserPolicy(userName string, obj ObjectType, sid string, op Operation) UserPolicy {
	return UserPolicy{userName, obj.String(), sid, op.String()}
}

func (p UserPolicy) Raw() []string {
	return p[:]
}

func (p UserPolicy) GetUserName() string {
	return p[0]
}

func (p UserPolicy) GetObjectType() ObjectType {
	return ObjectType(p[1])
}

func (p UserPolicy) GetSID() string {
	return p[2]
}

func (p UserPolicy) GetOperation() Operation {
	return Operation(p[3])
}

// GroupPolicy is a RBAC policy for group.
// e.g. ["master", "clusters", "*", "*"]
type GroupPolicy [4]string

// NewGroupPolicy returns a policy slice for group.
// e.g. ["master", "clusters", "*", "*"]: master is a group name
func NewGroupPolicy(groupName string, obj ObjectType, sid string, op Operation) GroupPolicy {
	return GroupPolicy{groupName, obj.String(), sid, op.String()}
}

func (p GroupPolicy) Raw() []string {
	return p[:]
}

func (p GroupPolicy) GetGroupName() string {
	return p[0]
}

func (p GroupPolicy) GetObjectType() ObjectType {
	return ObjectType(p[1])
}

func (p GroupPolicy) GetSID() string {
	return p[2]
}

func (p GroupPolicy) GetOperation() Operation {
	return Operation(p[3])
}

// GroupBinding binds a user to a group.
// e.g. ["foo", "master"]: user foo belongs to group master
type GroupBinding [2]string

// NewGroupBinding returns a binding slice for relationship between user and group.
func NewGroupBinding(userName, groupName string) GroupBinding {
	return GroupBinding{userName, groupName}
}

func (p GroupBinding) Raw() []string {
	return p[:]
}

func (p GroupBinding) GetUserName() string {
	return p[0]
}

func (p GroupBinding) GetGroupName() string {
	return p[1]
}

// AdminPolicy is the specific policy for admin/root user.
// \b(read|write|update)\b
var (
	AdminPolicy           = NewGroupPolicy(AdminGroup, ObjectAll, SidAll, OpAll)
	ReadOnlyPolicy        = NewGroupPolicy(ReadOnlyGroup, ObjectAll, SidAll, buildOperation(OpRead))
	ReadWritePolicy       = NewGroupPolicy(ReadWriteGroup, ObjectAll, SidAll, buildOperation(OpRead, OpCreate))
	ReadWriteUpdatePolicy = NewGroupPolicy(ReadWriteUpdateGroup, ObjectAll, SidAll, buildOperation(OpRead, OpCreate, OpUpdate))
)

// IsAdminPolicy returns true if the policy is the admin policy.
func IsAdminPolicy(policy Policy) bool {
	switch p := policy.(type) {
	case GroupPolicy:
		return reflect.DeepEqual(p.Raw(), AdminPolicy.Raw())
	default:
		return false
	}
}

// BindingToAdmin returns true if policy binding to admin group exists.
func BindingToAdmin(policies []GroupBinding) bool {
	for _, policy := range policies {
		if policy.GetGroupName() == AdminGroup {
			return true
		}
	}
	return false
}

// NewPolicyFromModels returns a policy slice.
// e.g. ["foo", "clusters", "*", "*"]
func NewPolicyFromModels(user *User, obj ObjectType, model base.Model, op Operation) Policy {
	return NewUserPolicy(user.Name, obj, model.GetSID(), op)
}

// NOTE: GetIdRangeFromPolicy is only used for listing API request.
// GetIdRangeFromPolicy returns true and an empty list when policy with all operation(*) are allowed exists,
// otherwise it returns false and a list of object IDs.
func GetIdRangeFromPolicy(policies []Policy) (all bool, ids []int64) {
	ids = make([]int64, 0)
	for _, policy := range policies {
		if _, ok := policy.(GroupBinding); ok {
			continue
		}

		raw := policy.Raw() // e.g. ["foo", "clusters", "*", "read"]
		sid := raw[2]
		// TODO: GetIdRangeFromPolicy 当前不准确
		switch sid {
		case "":
			continue
		case SidAll:
			// permit to read all
			return true, []int64{}
		}

		// operation
		if !(raw[3] == OpRead.String() || raw[3] == OpAll.String()) {
			continue
		}

		id, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			// invalid sid
			continue
		}
		ids = append(ids, id)
	}
	return false, utils.DeduplicateIntSlice(ids)
}

// 自定义 keyMatch 函数，支持正则表达式匹配
func CustomKeyMatch(parameters ...interface{}) (interface{}, error) {
	if len(parameters) != 2 {
		return false, fmt.Errorf("expected 2 parameters, got %d", len(parameters))
	}

	key1, ok1 := parameters[0].(string)
	key2, ok2 := parameters[1].(string)

	if !ok1 || !ok2 {
		return false, fmt.Errorf("parameters must be strings")
	}

	matched, _ := regexp.MatchString(key2, key1)
	return matched, nil
}
