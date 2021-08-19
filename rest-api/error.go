package restapi

import "errors"

var (
	ErrCreateRequest  = errors.New("failed to create request")
	ErrSendRequest    = errors.New("failed to send request")
	ErrRequestTimeout = errors.New("request timed out")
	ErrInvalidStatus  = errors.New("invalid status")
	ErrReadResponse   = errors.New("failed to read response")
)

type Error struct {
	Code    int		`json:"code"`
	Message string 	`json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}
