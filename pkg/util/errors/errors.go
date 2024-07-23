package errors

import (
	"errors"
	"gorm.io/gorm"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrRecordNotFound   = gorm.ErrRecordNotFound
	ErrRecordNotUpdate  = errors.New("record not updated")
	ErrBusySystem       = errors.New("系统繁忙，请稍后再试")
	ErrReqParams        = errors.New("请求参数错误")
	ErrCloudNotRegister = errors.New("cloud 集群未注册")
	ErrUserNotFound     = errors.New("用户不存在")
	ErrUserPassword     = errors.New("密码错误")

	ParamsError        = errors.New("参数错误")
	OperateFailed      = errors.New("操作失败")
	NoPermission       = errors.New("无权限")
	InnerError         = errors.New("内部错误")
	NoUserIdError      = errors.New("请登录")
	RoleExistError     = errors.New("角色已存在")
	RoleNotExistError  = errors.New("角色不存在")
	MenusExistError    = errors.New("权限已存在")
	MenusNtoExistError = errors.New("权限不存在")
)

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsNotUpdated(err error) bool {
	return errors.Is(err, ErrRecordNotUpdate)
}

func IsUniqueConstraintError(err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	// 数据库的 1062 错误码为固定的主键冲突号
	return mysqlErr.Number == 1062
}