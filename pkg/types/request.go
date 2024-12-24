package types

import "kubevulpes/pkg/db/model"

type (
	// LoginRequest 登录请求
	LoginRequest struct {
		Name     string `json:"name" binding:"required"`     // required
		Password string `json:"password" binding:"required"` // required
	}

	// RBACPolicyRequest RBAC策略请求
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
		Page    int    `form:"page" json:"page"`       // 页数，表示第几页
		Limit   int    `form:"limit" json:"limit"`     // 每页数量
		Keyword string `form:"keyword" json:"keyword"` // 排序 升序 asc，降序 desc，默认desc
	}

	// QueryOption 搜索配置
	QueryOption struct {
		LabelSelector string `form:"labelSelector" json:"labelSelector"` // 标签搜索
		NameSelector  string `form:"nameSelector" json:"nameSelector"`   // 名称搜索
	}

	// CreateUserRequest 创建用户请求
	CreateUserRequest struct {
		Name        string           `json:"name" binding:"required"`              // required
		Password    string           `json:"password" binding:"required,password"` // required
		Role        model.UserRole   `json:"role" binding:"omitempty,oneof=0 1 2"` // optional
		Status      model.UserStatus `json:"status" binding:"omitempty"`
		Email       string           `json:"email" binding:"omitempty,email"` // optional
		Description string           `json:"description" binding:"omitempty"` // optional
	}

	UpdateUserRequest struct {
		Role            model.UserRole   `json:"role" binding:"omitempty,oneof=0 1 2"`   // required
		Status          model.UserStatus `json:"status" binding:"omitempty,oneof=0 1 2"` // required
		Email           string           `json:"email" binding:"omitempty,email"`        // optional
		Description     string           `json:"description" binding:"omitempty"`        // optional
		ResourceVersion *int64           `json:"resource_version" binding:"required"`    // required
	}

	// CreateRoleRequest 创建角色请求
	CreateRoleRequest struct {
		Name        string `json:"name" binding:"required"` // required
		Description string `json:"description" binding:"omitempty"`
	}

	UpdateRoleRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"omitempty"` // optional
	}

	// CreateClusterRequest 创建集群请求
	CreateClusterRequest struct {
		Name      string `json:"name" binding:"omitempty"`       // optional
		AliasName string `json:"alias_name" binding:"omitempty"` // optional
		//Type        model.ClusterType `json:"cluster_type" binding:"omitempty,oneof=0 1"` // optional
		KubeConfig  string `json:"kube_config" binding:"required"`  // required
		Description string `json:"description" binding:"omitempty"` // optional
		Protected   bool   `json:"protected" binding:"omitempty"`   // optional
	}

	UpdateClusterRequest struct {
		AliasName       *string `json:"alias_name" binding:"omitempty"`      // optional
		Description     *string `json:"description" binding:"omitempty"`     // optional
		ResourceVersion *int64  `json:"resource_version" binding:"required"` // required
	}

	ProtectClusterRequest struct {
		ResourceVersion *int64 `json:"resource_version" binding:"required"` // required
		Protected       bool   `json:"protected" binding:"omitempty"`       // optional
	}
)

type (
	LoginResponse struct {
		UserId      int64          `json:"user_id"`
		UserName    string         `json:"user_name"`
		Token       string         `json:"token"`
		Role        model.UserRole `json:"role"`
		*model.User `json:"-"`
	}

	// PageResponse 分页查询返回值
	PageResponse struct {
		PageRequest `json:",inline"` // 分页请求属性

		Total int         `json:"total"` // 分页总数
		Items interface{} `json:"items"` // 指定页的元素列表
	}
)
