package commons

type contextKey int

// http context
const (
	ContextTypeUnknown contextKey = iota
	ContextTypeAuth
	ContextTypeRequestTimestamp
	ContextTypeLogger
)
