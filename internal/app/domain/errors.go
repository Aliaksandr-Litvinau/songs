package domain

import (
	"songs/internal/app/common/slugerrors"
)

var (
	ErrNotFound = slugerrors.NewError(
		"not-found",
		slugerrors.ErrorTypeNotFound,
		"resource not found",
	)

	ErrInvalidData = slugerrors.NewError(
		"invalid-data",
		slugerrors.ErrorTypeBadRequest,
		"invalid data provided",
	)

	ErrInvalidID = slugerrors.NewError(
		"invalid-id",
		slugerrors.ErrorTypeBadRequest,
		"invalid ID provided",
	)

	ErrRequired = slugerrors.NewError(
		"required-field",
		slugerrors.ErrorTypeBadRequest,
		"required field is missing",
	)

	ErrDuplicate = slugerrors.NewError(
		"duplicate-entry",
		slugerrors.ErrorTypeBadRequest,
		"duplicate entry",
	)

	ErrDatabase = slugerrors.NewError(
		"database-error",
		slugerrors.ErrorTypeInternal,
		"database error occurred",
	)

	ErrValidation = slugerrors.NewError(
		"validation-error",
		slugerrors.ErrorTypeBadRequest,
		"validation failed",
	)

	ErrInternal = slugerrors.NewError(
		"internal-error",
		slugerrors.ErrorTypeInternal,
		"internal server error",
	)
)
