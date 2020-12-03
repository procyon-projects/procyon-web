
<img src="https://procyon-projects.github.io/img/logo.png" width="128">

# Procyon Web
![alt text](https://goreportcard.com/badge/github.com/procyon-projects/procyon-web)
[![Build Status](https://travis-ci.com/procyon-projects/procyon-web.svg?branch=master)](https://travis-ci.com/procyon-projects/procyon-web)


This gives you a basic understanding of Procyon Web Module. It covers
components provided by the framework, such as Controller, Handler Registry, Web Request Context
and Interceptors.


Note that you need to register your components by using the function **core.Register**.

## Controller
It's used to define a struct as a **Controller Component** and register handler methods. If it is implemented by your struct,
your struct will be considered as a **Controller**.

```go
type Controller interface {
	RegisterHandlers(registry HandlerRegistry)
}
```

## Handler Registry
It's used to register handler methods.
```go
type HandlerRegistry interface {
	Register(info ...RequestHandler)
	RegisterGroup(prefix string, info ...RequestHandler)
}
```

* **Register** is used to register a Request Handler.

* **RegisterGroup** are used to register multiple Request Handlers.

### Request Handler
They are used to get an instance of a request handler by request method.
```go
func Get(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler
func Post(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler 
func Put(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler
func Delete(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler
func Patch(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler
func Options(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler
func Head(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler
```

### Request Handler Options
Options are used to specify handler's properties like Path and Request Object.
```go
func RequestObject(requestObject RequestHandlerObject) RequestHandlerOption
func Path(path string) RequestHandlerOption
```

* **RequestObject** is used to specify the request object. If you have a request type, you have to register it.
Otherwise, **GetRequest** will throw an error.
* **Path** is used to specify the path.

## Web Request Context
**WebRequestContext** implements the interface **context.Context** in procyon-context. That's why
it has the methods the following.

```go
func (ctx *WebRequestContext) GetContextId() context.ContextId
func (ctx *WebRequestContext) Get(key string) interface{}
func (ctx *WebRequestContext) Put(key string, value interface{})
```
* **GetContextId** returns a Context Id. It is unique and consists of **UUID**. It can be used
for logging. 
* **Get** returns the value from context by the given key. If it is not found, it returns nil.
* **Put** an key-value pair into context.


```go
func (ctx *WebRequestContext) Next()
```
* **Next** can only be invoked from **HandleBefore**. If you call it from other methods or
functions, nothing will happen. When **Next** is invoked, it calls next interceptor implementing 
**HandlerInterceptorBefore**.

```go
func (ctx *WebRequestContext) GetRequest(request interface{})
func (ctx *WebRequestContext) GetPathVariable(name string) string
func (ctx *WebRequestContext) GetRequestParameter(name string) string
func (ctx *WebRequestContext) GetHeaderValue(key string) string
```

* **GetRequest** is used to bind the request data to the instance of the request object.
* **GetPathVariable** is used to get the path variable by name
* **GetRequestParameter** is used to get the request parameter by name.
* **GetHeaderValue** is used to get the header value by name.

```go
func (ctx *WebRequestContext) GetStatus() int
func (ctx *WebRequestContext) SetStatus(status int) ResponseBodyBuilder

func (ctx *WebRequestContext) SetBody(body interface{}) ResponseBodyBuilder
func (ctx *WebRequestContext) GetBody() interface{}

func (ctx *WebRequestContext) GetContentType() MediaType 
func (ctx *WebRequestContext) SetContentType(mediaType MediaType) ResponseBodyBuilder

func (ctx *WebRequestContext) AddHeader(key string, value string) ResponseHeaderBuilder
```
* **GetStatus** is used to get the status.
* **SetStatus** is used to set the status of response.
* **SetBody** is used to set the body of response.
* **GetBody** is used to get the body of response.
* **GetContentType** is used to get the content type.
* **SetContentType** is used to set the content type

```go
func (ctx *WebRequestContext) Ok() ResponseBodyBuilder
func (ctx *WebRequestContext) NotFound() ResponseHeaderBuilder
func (ctx *WebRequestContext) NoContent() ResponseHeaderBuilder
func (ctx *WebRequestContext) BadRequest() ResponseBodyBuilder
func (ctx *WebRequestContext) Accepted() ResponseBodyBuilder
func (ctx *WebRequestContext) Created(location string) ResponseBodyBuilder
```

* **Ok** sets the status to 200.
* **NotFound** sets the status to 404.
* **NoContent** sets the status to 204.
* **BadRequest** sets the status to 400.
* **Accepted** sets the status to 202.
* **Created** sets the status to 201.

```go
func (ctx *WebRequestContext) GetError() error
func (ctx *WebRequestContext) SetError(err error)
func (ctx *WebRequestContext) ThrowError(err error)
```

* **GetError** is used to get the error.
* **SetError** is used to put the error into context.
* **ThrowError** is used to throw an error. It's an alternative to **SetError**

## Interceptors
Interceptors are used to manipulate requests and responses. 

**HandlerBefore**,Handler Method and **HandlerAfter**, **HandleAfterCompletion** are invoked
respectively.

### Interceptor Before
If you want to do something before handler method is executed, implement the interface 
**HandlerInterceptorBefore**.
```go
type HandlerInterceptorBefore interface {
	HandleBefore(requestContext *WebRequestContext)
}
```
### Interceptor After
If you want to do something after handler method is executed, implement the interface
**HandlerInterceptorAfter**.
```go
type HandlerInterceptorAfter interface {
	HandleAfter(requestContext *WebRequestContext)
}
```

### Interceptor After Completion
If you want to do something after the request process is completed, implement the interface
**HandlerInterceptorAfterCompletion**. **HandlerAfterCompletion** is invoked after response is
returned successfully or when any error occurs while request is processed. In case of an error,
You can get the error from request context.
```go
type HandlerInterceptorAfterCompletion interface {
	HandleAfterCompletion(requestContext *WebRequestContext)
}
```

## License
Procyon Framework is released under version 2.0 of the Apache License
