package response

type StatusCode int

const (
	StatusOk            StatusCode = 200
	StatusInternalError StatusCode = 503
	StatusUnauthorized  StatusCode = 401
	StatusForbidden     StatusCode = 403
	StatusNotFound      StatusCode = 404
	StatusInternalerror StatusCode = 400
)

type Response interface {
	Write(data []byte, code StatusCode) error
}
