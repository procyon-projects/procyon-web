package web

import "net/http"

type Response interface {
	GetStatus() int
	GetBody() interface{}
}

type ResponseEntity struct {
	status int
	body   interface{}
}

type ResponseEntityOption func(body *ResponseEntity)

func NewResponseEntity(options ...ResponseEntityOption) *ResponseEntity {
	body := &ResponseEntity{
		status: http.StatusOK,
	}
	for _, option := range options {
		option(body)
	}
	return body
}

func WithStatus(status int) ResponseEntityOption {
	return func(body *ResponseEntity) {
		body.status = status
	}
}

func WithBody(body interface{}) ResponseEntityOption {
	return func(body *ResponseEntity) {
		body.body = body
	}
}

func (body *ResponseEntity) GetStatus() int {
	return body.status
}

func (body *ResponseEntity) GetBody() interface{} {
	return body.body
}
