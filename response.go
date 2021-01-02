package web

import (
	"encoding/xml"
	json "github.com/json-iterator/go"
)

type MediaType byte

const (
	DefaultMediaType                       = MediaTypeApplicationTextHtml
	MediaTypeApplicationTextHtml MediaType = 0
	MediaTypeApplicationJson     MediaType = 1
	MediaTypeApplicationXml      MediaType = 2
)

const (
	DefaultMediaTypeValue             = MediaTypeApplicationTextHtmlValue
	MediaTypeApplicationTextHtmlValue = "text/html"
	MediaTypeApplicationXmlValue      = "application/xml"
	MediaTypeApplicationJsonValue     = "application/json"
)

type ResponseHeaderBuilder interface {
	AddResponseHeader(key string, value string) ResponseHeaderBuilder
}

type ResponseBodyBuilder interface {
	ResponseHeaderBuilder
	SetResponseStatus(status int) ResponseBodyBuilder
	SetModel(body interface{}) ResponseBodyBuilder
	SetResponseContentType(mediaType MediaType) ResponseBodyBuilder
}

type Response interface {
	GetResponseLocation() string
	GetResponseStatus() int
	GetModel() interface{}
	GetResponseBody() []byte
	GetResponseContentType() MediaType
	GetResponseHeader(key string) (string, bool)
}

type ResponseEntity struct {
	model       interface{}
	location    string
	status      int
	contentType MediaType
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
