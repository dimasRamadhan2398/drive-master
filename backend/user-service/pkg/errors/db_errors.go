package errors

import (
	"errors"

	"gorm.io/gorm"
)

func TranslateDBError(err error) error {
	if err == nil {
		return nil
	}

	// GORM-level errors
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound
	case errors.Is(err, gorm.ErrInvalidValue),
		errors.Is(err, gorm.ErrInvalidValueOfLength):
		return ErrBadRequest
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrConflict
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return ErrBadRequest
	}

	return ErrInternalServer
}
