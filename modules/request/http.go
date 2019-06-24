package request

import (
	"io"
	"net/http"
)

type HttpRequest struct {
	request *http.Request
}

func FromHttp(r *http.Request) Request {
	h := &HttpRequest{
		request: r,
	}
	return h
}

func (h *HttpRequest) Headers() map[string][]string {
	return h.request.Header
}

func (h *HttpRequest) Body() io.ReadCloser {
	return h.request.Body
}

func (h *HttpRequest) Ip() string {
	return h.request.RemoteAddr
}

func (h *HttpRequest) RealIp() string {
	return h.request.RemoteAddr
}

func (h *HttpRequest) UserAgent() string {
	return h.request.UserAgent()
}

func NewHttp(r *http.Request) Request {
	return &HttpRequest{
		request: r,
	}

}
