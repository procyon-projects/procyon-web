package web

import (
	"io"
	"net/http"
	"net/url"
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

func (req HttpRequest) GetBody() io.ReadCloser {
	return req.request.Body
}

func (req HttpRequest) GetHeader() http.Header {
	return req.request.Header
}

func (req HttpRequest) GetPath() string {
	return req.request.URL.Path
}

func (req HttpRequest) GetQueryParameters() url.Values {
	return req.request.URL.Query()
}

func (req HttpRequest) AddAttribute(key string, value interface{}) {
	req.attributes[key] = value
}

func (req HttpRequest) HasAttribute(key string) bool {
	if _, ok := req.attributes[key]; ok {
		return true
	}
	return false
}

func (req HttpRequest) GetAttribute(key string) interface{} {
	return req.attributes[key]
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
