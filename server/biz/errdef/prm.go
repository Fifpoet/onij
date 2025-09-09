package errdef

import "errors"

var (
	ErrInvalidFileId   = errors.New("invalid file id")
	ErrInvalidMusicId  = errors.New("invalid music id")
	ErrInvalidAlbumId  = errors.New("invalid album id")
	ErrInvalidArtistId = errors.New("invalid artist id")
	ErrInvalidTagId    = errors.New("invalid tag id")
)
