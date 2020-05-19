package web

import "net/http"

type Handler interface {
	DoGet(res http.ResponseWriter, req *http.Request)
	DoPost(res http.ResponseWriter, req *http.Request)
	DoPatch(res http.ResponseWriter, req *http.Request)
	DoPut(res http.ResponseWriter, req *http.Request)
	DoDelete(res http.ResponseWriter, req *http.Request)
}

type DefaultHandler struct {
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (handler *DefaultHandler) DoGet(res http.ResponseWriter, req *http.Request) {

}

func (handler *DefaultHandler) DoPost(res http.ResponseWriter, req *http.Request) {

}

func (handler *DefaultHandler) DoPatch(res http.ResponseWriter, req *http.Request) {

}

func (handler *DefaultHandler) DoPut(res http.ResponseWriter, req *http.Request) {

}

func (handler *DefaultHandler) DoDelete(res http.ResponseWriter, req *http.Request) {

}
