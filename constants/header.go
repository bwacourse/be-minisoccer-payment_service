package constants

import "net/textproto"

var (
	XServiceName  = textproto.CanonicalMIMEHeaderKey("X-Service-Name")
	XAPIKey       = textproto.CanonicalMIMEHeaderKey("X-Api-Key")
	XRequestAt    = textproto.CanonicalMIMEHeaderKey("X-Request-At")
	Authorization = textproto.CanonicalMIMEHeaderKey("Authorization")
)
