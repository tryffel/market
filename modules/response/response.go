package response

import (
	"encoding/json"
	"fmt"
	"github.com/tryffel/market/modules/Error"
)

type StatusCode int

const (
	StatusOk            StatusCode = 200
	StatusInternalError StatusCode = 503
	StatusBadRequest    StatusCode = 400
	StatusUnauthorized  StatusCode = 401
	StatusForbidden     StatusCode = 403
	StatusNotFound      StatusCode = 404
)

type Response interface {
	Write(data []byte, code StatusCode) error
}

func Ok(msg string, resp Response) error {
	return resp.Write([]byte(msg), StatusOk)
}

func BadRequest(errMsg string, resp Response) error {
	return resp.Write([]byte(fmt.Sprintf(`{"Error": "%s"}`, errMsg)), StatusBadRequest)
}

func Forbidden(resp Response) error {
	msg := `{"Error": "Forbidden"}`
	return resp.Write([]byte(msg), StatusForbidden)
}

func Unauthorized(resp Response) error {
	msg := `{"Error": "Unauthorized"}`
	return resp.Write([]byte(msg), StatusUnauthorized)
}

func JsonOk(body interface{}, resp Response) error {
	data, err := json.Marshal(body)
	if err != nil {
		return Error.Wrap(&err, "failed to marshal body")
	}
	err = resp.Write([]byte(data), StatusOk)
	return err
}
