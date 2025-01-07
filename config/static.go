package config

import "time"

var (
	ErrorInvalidRequest = "INVALID_REQUEST"
	ErrorInvalidToken   = "INVALID_TOKEN"
	ErrorInvalidUser    = "INVALID_USER"
	ErrorInvalidPass    = "INVALID_PASS"
	ErrorInvalidEmail   = "INVALID_EMAIL"
	ErrorInvalidPhone   = "INVALID_PHONE"
	ErrorSessionExpired = "SESSION_EXPIRED"
	ErrorInternalServer = "INTERNAL_SERVER"
	ErrorNotFound       = "NOT_FOUND"
	ErrorUnauthorized   = "UNAUTHORIZED"
	ErrorForbidden      = "FORBIDDEN"
	ErrorConflict       = "CONFLICT"
	ErrorBadRequest     = "BAD_REQUEST"
	ErrorDuplicateKey   = "DUPLICATE_KEY"
)

var (
	TokenExpireTime = 24 * time.Hour * 7 // 7 days
)
