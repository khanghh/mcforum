package service

import (
	"bbs-go/internal/locale"
	"errors"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

func IsDatabaseError(err error) bool {
	if err == nil {
		return false
	}

	gormErrors := []error{
		gorm.ErrRecordNotFound,
		gorm.ErrInvalidTransaction,
		gorm.ErrNotImplemented,
		gorm.ErrMissingWhereClause,
		gorm.ErrUnsupportedRelation,
		gorm.ErrPrimaryKeyRequired,
		gorm.ErrModelValueRequired,
		gorm.ErrModelAccessibleFieldsRequired,
		gorm.ErrSubQueryRequired,
		gorm.ErrInvalidData,
		gorm.ErrUnsupportedDriver,
		gorm.ErrRegistered,
		gorm.ErrInvalidField,
		gorm.ErrEmptySlice,
		gorm.ErrDryRunModeUnsupported,
		gorm.ErrInvalidDB,
		gorm.ErrInvalidValue,
		gorm.ErrInvalidValueOfLength,
		gorm.ErrPreloadNotAllowed,
		gorm.ErrDuplicatedKey,
		gorm.ErrForeignKeyViolated,
		gorm.ErrCheckConstraintViolated,
	}

	for _, gErr := range gormErrors {
		if errors.Is(err, gErr) {
			return true
		}
	}

	return false
}

type ResponseError struct {
	code    int
	message string
}

func (e ResponseError) Code() int {
	return e.code
}

func (e ResponseError) Error() string {
	return e.message
}

var (
	ErrTopicNotFound  = ResponseError{iris.StatusNotFound, locale.T("topic.not_found")}
	ErrForbidden      = ResponseError{iris.StatusForbidden, locale.T("system.message.permission_denied")}
	ErrUnauthorized   = ResponseError{iris.StatusUnauthorized, locale.T("system.message.unauthorized")}
	ErrBadRequest     = ResponseError{iris.StatusBadRequest, locale.T("system.message.invalid_request")}
	ErrInternalServer = ResponseError{iris.StatusInternalServerError, locale.T("system.message.internal_server_error")}
)
