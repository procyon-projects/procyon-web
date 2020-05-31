package web

import "net/http"

type Response interface {
	GetStatus() int
	GetBody() interface{}
}

type ResponseBody struct {
	status int
	body   interface{}
}

type ResponseBodyOption func(body *ResponseBody)

func NewResponseBody(options ...ResponseBodyOption) *ResponseBody {
	body := &ResponseBody{
		status: http.StatusOK,
	}
	for _, option := range options {
		option(body)
	}
	return body
}

func WithStatus(status int) ResponseBodyOption {
	return func(body *ResponseBody) {
		body.status = status
	}
}

func WithBody(body interface{}) ResponseBodyOption {
	return func(body *ResponseBody) {
		body.body = body
	}
}
