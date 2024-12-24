package model

import (
	"fmt"
	"kubevulpes/pkg/db/model/base"
)

func init() {
	register(&Audit{})
}

type AuditOperationStatus uint8

func (s AuditOperationStatus) String() string {
	switch s {
	case 0:
		return "failed"
	case 1:
		return "succeed"
	default:
		return "unknown"
	}
}

type Audit struct {
	base.Model

	RequestId  string               `gorm:"column:request_id;types:varchar(32);index" json:"request_id"`  // 请求 ID
	Ip         string               `gorm:"types:varchar(128)" json:"ip"`                                 // 客户端 IP
	Action     string               `gorm:"types:varchar(255)" json:"action"`                             // HTTP 方法 [POST/DELETE/PUT/GET]
	Operator   string               `gorm:"types:varchar(255)" json:"operator"`                           // 操作人 ID
	Path       string               `gorm:"types:varchar(255)" json:"path"`                               // HTTP 路径
	ObjectType ObjectType           `gorm:"column:resource_type;types:varchar(128)" json:"resource_type"` // 操作资源类型 [cluster/plan...]
	Status     AuditOperationStatus `gorm:"types:tinyint" json:"status"`                                  // 记录操作运行结果[OperationStatus]
}

func (a *Audit) TableName() string {
	return "audits"
}

func (a *Audit) String() string {
	return fmt.Sprintf("user %s(ip addr: %s) access %s with %s then %s", a.Operator, a.Ip,
		a.Path, a.Action, a.Status.String())
}
