package commons

import "errors"

var (
	ErrInsufficientCredits = errors.New("insufficient credits, please purchase more and try again later")
	ErrInvalidParam        = errors.New("invalid parameter, please refresh the page and try again later")
	ErrInternalError       = errors.New("internal application error, please refresh the page and try again later")
)
