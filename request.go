package web

import (
	"github.com/codnect/goo"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

var cacheRequestObject = make(map[reflect.Type]*RequestObjectCache, 0)
var cacheRequestObjectMu sync.RWMutex

type RequestObjectCache struct {
	hasOnlyBody      bool
	bodyFieldIndex   int
	paramFieldIndex  int
	pathFieldIndex   int
	headerFieldIndex int
	fields           []goo.Field
}

type RequestHandlerObject interface{}
type RequestHandlerFunction = func(context *WebRequestContext)
type RequestHandlerOption func(handler *RequestHandler)

type RequestMethod string

const (
	unknownMethod        RequestMethod = "[unknown-method]"
	RequestMethodGet     RequestMethod = http.MethodGet
	RequestMethodPost    RequestMethod = http.MethodPost
	RequestMethodPut     RequestMethod = http.MethodPut
	RequestMethodDelete  RequestMethod = http.MethodDelete
	RequestMethodPatch   RequestMethod = http.MethodPatch
	RequestMethodOptions RequestMethod = http.MethodOptions
	RequestMethodHead    RequestMethod = http.MethodHead
)

type RequestHandler struct {
	Path          string
	Method        RequestMethod
	HandlerFunc   RequestHandlerFunction
	RequestObject RequestHandlerObject
}

func newHandler(handler RequestHandlerFunction, method RequestMethod, options ...RequestHandlerOption) RequestHandler {
	if handler == nil {
		panic("Handler must not be null")
	}

	handlerType := goo.GetType(handler)
	if !handlerType.IsFunction() {
		panic("Handler must be function")
	}

	requestHandler := &RequestHandler{
		HandlerFunc: handler,
		Method:      method,
	}

	for _, option := range options {
		option(requestHandler)
	}

	if requestHandler.RequestObject != nil {
		requestObjType := goo.GetType(requestHandler.RequestObject)
		if !requestObjType.IsStruct() {
			panic("Request object must be struct")
		}
		scanRequestObject(requestObjType)
	}

	return *requestHandler
}

func Get(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodGet, options...)
}

func Post(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodPost, options...)
}

func Put(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodPut, options...)
}

func Delete(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodDelete, options...)
}

func Patch(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodPatch, options...)
}

func Options(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodOptions, options...)
}

func Head(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodHead, options...)
}

func scanRequestObject(requestObjType goo.Type) {
	structType := requestObjType.ToStructType()
	if structType.GetFieldCount() == 0 {
		return
	}
	fields := structType.GetFields()

	requestObjcCache := &RequestObjectCache{
		hasOnlyBody:      false,
		bodyFieldIndex:   -1,
		paramFieldIndex:  -1,
		pathFieldIndex:   -1,
		headerFieldIndex: -1,
		fields:           structType.GetFields(),
	}

	hasField := false
	requestStruct := false
	for index, field := range fields {
		fieldType := field.GetType()

		if fieldType.IsStruct() && strings.HasPrefix(fieldType.GetName(), "struct") {
			requestStruct = true
		} else {
			hasField = true
		}

		if requestStruct && hasField {
			panic("Request Object must only consist of untyped request structs or fields completely")
		}

		if hasField {
			continue
		}

		requestTag, err := field.GetTagByName("request")

		if err != nil {
			panic("Untyped struct must have request tag in Request Object")
		}

		switch requestTag.Value {
		case "param":
			validateRequestStruct(requestTag.Value, field.GetType().ToStructType())
			requestObjcCache.paramFieldIndex = index
		case "body":
			requestObjcCache.bodyFieldIndex = index
		case "path":
			validateRequestStruct(requestTag.Value, field.GetType().ToStructType())
			requestObjcCache.pathFieldIndex = index
		case "header":
			validateRequestStruct(requestTag.Value, field.GetType().ToStructType())
			requestObjcCache.headerFieldIndex = index
		default:
			panic("Invalid request tag value")
		}

	}

	if hasField {
		requestObjcCache.hasOnlyBody = true
	}
	cacheRequestObject[structType.GetGoType()] = requestObjcCache
}

func validateRequestStruct(requestStructType string, requestStruct goo.Struct) {
	if requestStruct == nil {
		return
	}
	if "param" == requestStructType || "path" == requestStructType || "header" == requestStructType {
		fields := requestStruct.GetFields()
		for _, field := range fields {
			fieldType := field.GetType()
			if !fieldType.IsString() && !fieldType.IsBoolean() && !fieldType.IsNumber() {
				panic("Fields could be string, boolean and number types")
			}
		}
	}
}

func RequestObject(requestObject RequestHandlerObject) RequestHandlerOption {
	return func(handler *RequestHandler) {
		handler.RequestObject = requestObject
	}
}

func Path(path string) RequestHandlerOption {
	return func(handler *RequestHandler) {
		handler.Path = path
	}
}

type HandlerRegistry interface {
	Register(info ...RequestHandler)
	RegisterGroup(prefix string, info ...RequestHandler)
}

type SimpleHandlerRegistry struct {
	registryMap map[string][]RequestHandler
}

func NewSimpleHandlerRegistry() SimpleHandlerRegistry {
	return SimpleHandlerRegistry{
		registryMap: make(map[string][]RequestHandler),
	}
}

func (registry SimpleHandlerRegistry) Register(info ...RequestHandler) {
	registry.RegisterGroup("", info...)
}

func (registry SimpleHandlerRegistry) RegisterGroup(prefix string, info ...RequestHandler) {
	if len(info) == 0 {
		return
	}

	for _, handler := range info {
		if prefix+handler.Path == "" {
			panic("Specify a path or a prefix while registering a request handler")
		}
	}

	if registry.registryMap[prefix] == nil {
		registry.registryMap[prefix] = make([]RequestHandler, 0)
	}
	registry.registryMap[prefix] = append(registry.registryMap[prefix], info...)
}

func (registry SimpleHandlerRegistry) clear() {
	for key := range registry.registryMap {
		delete(registry.registryMap, key)
	}
}

func (registry SimpleHandlerRegistry) getRegistryMap() map[string][]RequestHandler {
	return registry.registryMap
}
