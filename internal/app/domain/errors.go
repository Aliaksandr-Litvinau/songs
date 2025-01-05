package domain

import "songs/internal/app/common/slugerrors"

var (
	ErrSongNotFound = slugerrors.NewError(
		"song-not-found",
		slugerrors.ErrorTypeNotFound,
		"song not found",
	)

	ErrInvalidSongID = slugerrors.NewError(
		"invalid-song-id",
		slugerrors.ErrorTypeBadRequest,
		"invalid song ID",
	)

	ErrInvalidSongData = slugerrors.NewError(
		"invalid-song-data",
		slugerrors.ErrorTypeBadRequest,
		"invalid song data",
	)
)
