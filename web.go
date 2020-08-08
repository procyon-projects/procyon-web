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

func (req HttpRequest) GetMethod() string {
	return req.request.Method
}

func (req HttpRequest) AddAttribute(key string, value interface{}) {
	req.attributes[key] = value
}

func (req HttpRequest) clearAttributes() {
	for key := range req.attributes {
		delete(req.attributes, key)
	}
}

type HttpResponse struct {
	responseWriter http.ResponseWriter
}

func newHttpResponse() interface{} {
	return HttpResponse{}
}
