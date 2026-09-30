package domain

import "errors"

var (
	ErrPostNotFound    = errors.New("post not found")
	ErrInvalidLocation = errors.New("invalid location")
	ErrInvalidKind     = errors.New("post kind must be help or meetup")
	ErrInvalidTitle    = errors.New("invalid title")
	ErrInvalidText     = errors.New("invalid description")
	ErrInvalidComment  = errors.New("invalid comment")
)
