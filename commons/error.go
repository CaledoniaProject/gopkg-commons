package commons

import "errors"

var (
	ErrAuthFailed          = errors.New("authentication failed")
	ErrNoKey               = errors.New("key field missing")
	ErrNoCoupon            = errors.New("no coupon available")
	ErrInsufficientCredits = errors.New("insufficient credits, please purchase more and try again later")
	ErrInvalidParam        = errors.New("invalid parameter, please refresh the page and try again later")
	ErrInternalError       = errors.New("internal application error, please refresh the page and try again later")
)
