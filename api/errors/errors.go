package errors

import (
	"net/http"

	"kubevulpes/pkg/util/errors"
)

type Error struct {
	Code int
	Err  error
}

func (e Error) Error() string {
	return e.Err.Error()
}

func NewError(err error, code int) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

var (
	ErrUnauthorized = Error{
		Code: http.StatusUnauthorized,
		Err:  errors.NoUserIdError,
	}
	ErrForbidden = Error{
		Code: http.StatusForbidden,
		Err:  errors.NoPermission,
	}
	ErrInvalidRequest = Error{
		Code: http.StatusBadRequest,
		Err:  errors.ErrReqParams,
	}
	ErrServerInternal = Error{
		Code: http.StatusInternalServerError,
		Err:  errors.ErrInternal,
	}
	ErrGroupBindingNotFound = Error{
		Code: http.StatusNotFound,
		Err:  errors.PolicyNotExistError,
	}
	ErrGroupBindingExists = Error{
		Code: http.StatusConflict,
		Err:  errors.PolicyExistError,
	}
	ErrRBACPolicyNotFound = Error{
		Code: http.StatusNotFound,
		Err:  errors.PolicyNotExistError,
	}
	ErrRBACPolicyExists = Error{
		Code: http.StatusConflict,
		Err:  errors.PolicyExistError,
	}
	ErrUserNotFound = Error{
		Code: http.StatusNotFound,
		Err:  errors.ErrUserNotFound,
	}
	ErrNotAcceptable = Error{
		Code: http.StatusNotAcceptable,
		Err:  errors.ErrNotAcceptable,
	}
	ErrUserExists = Error{
		Code: http.StatusConflict,
		Err:  errors.UserExistError,
	}
	ErrInvalidPassword = Error{
		Code: http.StatusConflict,
		Err:  errors.ErrUserPassword,
	}
	ErrDuplicatedPassword = Error{
		Code: http.StatusConflict,
		Err:  errors.ErrDuplicatedPassword,
	}
	ErrLicenseNotFound = Error{
		Code: http.StatusNotFound,
		Err:  errors.ErrClusterNotFound,
	}
	ErrProjectExists = Error{
		Code: http.StatusConflict,
		Err:  errors.ErrProjectDuplicatedName,
	}
	ErrProjectNotFound = Error{
		Code: http.StatusNotFound,
		Err:  errors.ErrTenantNotFound,
	}
	ErrAuditNotFound = Error{
		Code: http.StatusNotFound,
		Err:  errors.ErrAuditNotFound,
	}
	ErrAuditExists = Error{
		Code: http.StatusConflict,
		Err:  errors.ErrAuditExists,
	}
)
