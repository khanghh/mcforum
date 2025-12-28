package errs

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
	ErrForbidden      = ResponseError{iris.StatusForbidden, "Forbidden"}
	ErrUnauthorized   = ResponseError{iris.StatusUnauthorized, "Unauthorized"}
	ErrBadRequest     = ResponseError{iris.StatusBadRequest, "Bad Request"}
	ErrInternalServer = ResponseError{iris.StatusInternalServerError, "Internal Server Error"}
)

var (
	ErrForumNotFound   = ResponseError{iris.StatusNotFound, locale.T("errors.forum_not_found")}
	ErrTopicNotFound   = ResponseError{iris.StatusNotFound, locale.T("errors.topic_not_found")}
	ErrCommentNotFound = ResponseError{iris.StatusNotFound, locale.T("errors.comment_not_found")}
	ErrCommentDeleted  = ResponseError{iris.StatusNotFound, locale.T("errors.comment_deleted")}
	ErrUserNotFound    = ResponseError{iris.StatusNotFound, locale.T("errors.user_not_found")}
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
