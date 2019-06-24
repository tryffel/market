package response

import "net/http"

type HttpResponse struct {
	response http.ResponseWriter
}

func (h HttpResponse) Write(data []byte, code StatusCode) error {
	var err error
	if code == StatusOk {
		_, err = h.response.Write(data)
	} else {
		http.Error(h.response, string(data), int(code))
	}
	return err
}

func NewHttp(resp http.ResponseWriter) Response {
	return HttpResponse{response: resp}
}
