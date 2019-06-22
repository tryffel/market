package request

import "io"

// Request interface has all methods that can be used by api
type Request interface {
	Headers() map[string][]string
	Body() io.ReadCloser
	Ip() string
	RealIp() string
	UserAgent() string
}
