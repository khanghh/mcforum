package service

import (
	"errors"

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
