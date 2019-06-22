package response

import "net/http"

type HttpResponse struct {
	response http.ResponseWriter
}

func (h *HttpResponse) Write(data []byte, code StatusCode) error {
	var err error
	if code == StatusOk {
		_, err = h.response.Write(data)
	} else {
	}
	return err
}
