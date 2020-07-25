package web

import (
	"net/http"
)

type HttpRequest struct {
	request    *http.Request
	attributes map[string]interface{}
}

func newHttpRequest() interface{} {
	return HttpRequest{
		attributes: make(map[string]interface{}),
	}
}

func (req HttpRequest) AddAttribute(key string, value interface{}) {
	req.attributes[key] = value
}

type HttpResponse struct {
	responseWriter http.ResponseWriter
}

func newHttpResponse() interface{} {
	return HttpResponse{}
}
