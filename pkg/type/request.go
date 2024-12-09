package types

import "kubevulpes/pkg/db/model"

type (
	RBACPolicyRequest struct {
		// user ID or group name is required
		UserId     *int64           `json:"user_id" binding:"required_without=GroupName,excluded_with=GroupName"`
		GroupName  *string          `json:"group_name" binding:"required_without=UserId,excluded_with=UserId"`
		ObjectType model.ObjectType `json:"object_type" binding:"required,rbac_object"`
		SID        string           `json:"sid" binding:"omitempty,rbac_sid"`
		Operation  model.Operation  `json:"operation" binding:"required,rbac_operation"`
	}

	ListRBACPolicyRequest struct {
		UserId       int64             `form:"user_id" binding:"required"`
		ObjectType   *model.ObjectType `form:"object_type" binding:"omitempty,required_with=UserId,rbac_object"`
		SID          *string           `form:"sid" binding:"omitempty,required_with=ObjectType,rbac_sid"`
		Operation    *model.Operation  `form:"operation" binding:"omitempty,required_with=SID,rbac_operation"`
		*PageRequest `json:",inline"`
	}

	GroupBindingRequest struct {
		UserId    int64  `json:"user_id" binding:"required"`
		GroupName string `json:"group_name" binding:"required"`
	}

	ListGroupBindingRequest struct {
		UserId       *int64  `form:"user_id" binding:"omitempty"`
		GroupName    *string `form:"group_name" binding:"omitempty"`
		*PageRequest `json:",inline"`
	}

	// PageRequest 分页配置
	PageRequest struct {
		Page  int `form:"page" json:"page"`   // 页数，表示第几页
		Limit int `form:"limit" json:"limit"` // 每页数量
	}

	// QueryOption 搜索配置
	QueryOption struct {
		LabelSelector string `form:"labelSelector" json:"labelSelector"` // 标签搜索
		NameSelector  string `form:"nameSelector" json:"nameSelector"`   // 名称搜索
	}
)
