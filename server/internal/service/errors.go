package service

import (
	"bbs-go/internal/locale"
	"errors"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

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

func NewBadRequestError(msg string) error {
	return ResponseError{
		code:    iris.StatusBadRequest,
		message: msg,
	}
}

var (
	ErrTopicNotFound   = ResponseError{iris.StatusNotFound, locale.T("topic.not_found")}
	ErrCommentNotFound = ResponseError{iris.StatusNotFound, locale.T("comment.not_found")}
	ErrCommentDeleted  = ResponseError{iris.StatusNotFound, locale.T("comment.deleted")}
	ErrUserNotFound    = ResponseError{iris.StatusNotFound, locale.T("user.not_found")}
	ErrForbidden       = ResponseError{iris.StatusForbidden, locale.T("errors.permission_denied")}
	ErrUnauthorized    = ResponseError{iris.StatusUnauthorized, locale.T("errors.unauthorized")}
	ErrBadRequest      = ResponseError{iris.StatusBadRequest, locale.T("errors.invalid_request")}
	ErrInternalServer  = ResponseError{iris.StatusInternalServerError, locale.T("errors.internal_server_error")}
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
