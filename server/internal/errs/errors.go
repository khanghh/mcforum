package errs

import (
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
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

func NewResponseError(code int, message string) error {
	return &ResponseError{
		code:    code,
		message: message,
	}
}

func NewBadRequestError(msg string) error {
	return &ResponseError{
		code:    iris.StatusBadRequest,
		message: msg,
	}
}

var (
	ErrForbidden          = &ResponseError{iris.StatusForbidden, "Forbidden"}
	ErrUnauthorized       = &ResponseError{iris.StatusUnauthorized, "Unauthorized"}
	ErrBadRequest         = &ResponseError{iris.StatusBadRequest, "Bad Request"}
	ErrInternalServer     = &ResponseError{iris.StatusInternalServerError, "Internal Server Error"}
	ErrBadGateway         = &ResponseError{iris.StatusBadGateway, "Bad Gateway"}
	ErrServiceUnavailable = &ResponseError{iris.StatusServiceUnavailable, "Service Unavailable"}
)

var (
	ErrForumNotFound       = &ResponseError{iris.StatusNotFound, locale.T("errors.forum_not_found")}
	ErrTopicNotFound       = &ResponseError{iris.StatusNotFound, locale.T("errors.topic_not_found")}
	ErrCommentNotFound     = &ResponseError{iris.StatusNotFound, locale.T("errors.comment_not_found")}
	ErrCommentDeleted      = &ResponseError{iris.StatusNotFound, locale.T("errors.comment_deleted")}
	ErrUserNotFound        = &ResponseError{iris.StatusNotFound, locale.T("errors.user_not_found")}
	ErrTopicNotUnderReview = &ResponseError{iris.StatusNotFound, locale.T("errors.topic_not_under_review")}
)

var (
	ErrUploadIncomplete    = &ResponseError{iris.StatusBadRequest, locale.T("upload.upload_incomplete")}
	ErrUnsupportedFileType = &ResponseError{iris.StatusUnsupportedMediaType, locale.T("upload.unsupported_file_type")}
	ErrFileTooLarge        = &ResponseError{iris.StatusRequestEntityTooLarge, locale.T("upload.file_too_large")}
	ErrAvatarTooLarge      = &ResponseError{iris.StatusRequestEntityTooLarge, locale.T("upload.avatar_too_large")}
	ErrAvatarSizeTooSmall  = &ResponseError{iris.StatusRequestEntityTooLarge, locale.T("upload.avatar_size_too_small")}
)

var (
	ErrAlreadyFollowing = &ResponseError{iris.StatusBadRequest, locale.T("follow.already_following")}
)

var (
	ErrTopicExceedMaxImages        = &ResponseError{iris.StatusBadRequest, locale.T("topic.exceed_max_images", constants.TopicMaxImageCount)}
	ErrTopicInvalidImageURL        = &ResponseError{iris.StatusBadRequest, locale.T("topic.invalid_image_url")}
	ErrTopicTitleRequired          = &ResponseError{iris.StatusBadRequest, locale.T("topic.title_required")}
	ErrTopicTitleMaxLengthExceeded = &ResponseError{iris.StatusBadRequest, locale.T("topic.title_max_length_exceeded", constants.TopicTitleMaxLength)}
	ErrTopicContentRequired        = &ResponseError{iris.StatusBadRequest, locale.T("topic.content_required")}
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
