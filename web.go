package web

import (
	"net/http"
)

type HttpRequest struct {
	*http.Request
	attributes map[string]interface{}
}

func newHttpRequest(req *http.Request) HttpRequest {
	return HttpRequest{
		Request:    req,
		attributes: make(map[string]interface{}),
	}
}

func (req HttpRequest) AddAttribute(key string, value interface{}) {
	req.attributes[key] = value
}

type HttpResponse struct {
	http.ResponseWriter
}

func newHttpResponse(res http.ResponseWriter) HttpResponse {
	return HttpResponse{
		res,
	}
}
