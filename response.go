package web

import (
	"encoding/xml"
	json "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
)

type MediaType string

const (
	DefaultMediaType            = MediaTypeTextHtml
	MediaTypeTextHtml MediaType = "text/html"
	MediaTypeXml      MediaType = "application/xml"
	MediaTypeJson     MediaType = "application/json"
)

type ResponseHeaderBuilder interface {
	Location(location string) ResponseHeaderBuilder
	Header(name string, value string) ResponseHeaderBuilder
}

type ResponseBodyBuilder interface {
	ResponseHeaderBuilder
	Body(body interface{}) ResponseBodyBuilder
	ContentType(contentType MediaType) ResponseBodyBuilder
}

type Response interface {
	Accepted() ResponseBodyBuilder
	BadRequest() ResponseBodyBuilder
	Created(location string) ResponseBodyBuilder
	Ok() ResponseBodyBuilder
	NoContent() ResponseHeaderBuilder
	NotFound() ResponseHeaderBuilder
	Status(status int) ResponseBodyBuilder
}

type ResponseEntity struct {
	body        interface{}
	location    string
	status      int
	contentType MediaType
	ctx         *fasthttp.Response
}

func (responseEntity *ResponseEntity) reset() {
	responseEntity.body = nil
	responseEntity.location = ""
	responseEntity.status = http.StatusOK
	responseEntity.contentType = MediaTypeJson
}

func (responseEntity *ResponseEntity) GetStatus() int {
	return responseEntity.status
}

func (responseEntity *ResponseEntity) GetBody() interface{} {
	return responseEntity.body
}

func (responseEntity *ResponseEntity) HasBody() bool {
	return responseEntity.body != nil
}

func (responseEntity *ResponseEntity) Accepted() ResponseBodyBuilder {
	responseEntity.reset()
	responseEntity.status = http.StatusAccepted
	return responseEntity
}

func (responseEntity *ResponseEntity) BadRequest() ResponseBodyBuilder {
	responseEntity.reset()
	responseEntity.status = http.StatusBadRequest
	return responseEntity
}

func (responseEntity *ResponseEntity) Created(location string) ResponseBodyBuilder {
	responseEntity.reset()
	responseEntity.status = http.StatusCreated
	responseEntity.location = location
	return responseEntity
}

func (responseEntity *ResponseEntity) Ok() ResponseBodyBuilder {
	responseEntity.reset()
	responseEntity.status = http.StatusOK
	return responseEntity
}

func (responseEntity *ResponseEntity) NoContent() ResponseHeaderBuilder {
	responseEntity.reset()
	responseEntity.status = http.StatusNoContent
	return responseEntity
}

func (responseEntity *ResponseEntity) NotFound() ResponseHeaderBuilder {
	responseEntity.reset()
	responseEntity.status = http.StatusNotFound
	return responseEntity
}

func (responseEntity *ResponseEntity) Status(status int) ResponseBodyBuilder {
	responseEntity.reset()
	responseEntity.status = status
	return responseEntity
}

func (responseEntity *ResponseEntity) Location(location string) ResponseHeaderBuilder {
	responseEntity.location = location
	return responseEntity
}

func (responseEntity *ResponseEntity) Header(name string, value string) ResponseHeaderBuilder {

	return responseEntity
}

func (responseEntity *ResponseEntity) Body(body interface{}) ResponseBodyBuilder {
	responseEntity.body = body
	return responseEntity
}

func (responseEntity *ResponseEntity) ContentType(contentType MediaType) ResponseBodyBuilder {
	responseEntity.contentType = contentType
	return responseEntity
}

type ResponseWriter struct {
}

func (responseWriter ResponseWriter) WriteResponse(ctx *WebRequestContext, body []byte) {
	if ctx == nil {
		return
	}
	ctx.fastHttpRequestContext.SetBody(body)
}

type ResponseBodyWriter interface {
	WriteResponseBody(ctx *WebRequestContext, responseWriter ResponseWriter) error
}

type defaultResponseBodyWriter struct {
}

func newDefaultResponseBodyWriter() defaultResponseBodyWriter {
	return defaultResponseBodyWriter{}
}

func (bodyWriter defaultResponseBodyWriter) WriteResponseBody(ctx *WebRequestContext, responseWriter ResponseWriter) error {
	if ctx.responseEntity.contentType == MediaTypeApplicationJson {
		if ctx.responseEntity.model == nil {
			return nil
		}

		result, err := json.Marshal(ctx.responseEntity.model)
		if err != nil {
			return err
		}
		responseWriter.WriteResponse(ctx, result)
	} else if ctx.responseEntity.contentType == MediaTypeApplicationTextHtml {
		if ctx.responseEntity.model == nil {
			return nil
		}

		switch ctx.responseEntity.model.(type) {
		case string:
			result := []byte(ctx.responseEntity.model.(string))
			responseWriter.WriteResponse(ctx, result)
		}
	} else {
		if ctx.responseEntity.model == nil {
			return nil
		}

		result, err := xml.Marshal(ctx.responseEntity.model)
		if err != nil {
			return err
		}
		responseWriter.WriteResponse(ctx, result)
	}
	return nil
}
