package web

import (
	"encoding/xml"
	"errors"
	json "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"reflect"
)

type RequestBinder interface {
	BindRequest(request interface{}, ctx *WebRequestContext) error
}

type defaultRequestBinder struct {
}

func newDefaultRequestBinder() defaultRequestBinder {
	return defaultRequestBinder{}
}

func (binder defaultRequestBinder) BindRequest(request interface{}, ctx *WebRequestContext) error {
	typ := reflect.TypeOf(request)
	if typ == nil {
		return errors.New("type cannot be determined as the given object is nil")
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	metadata := ctx.handlerChain.requestObjectMetadata
	if metadata == nil {
		return errors.New("you need to specify RequestObject for handler, that's why you cannot use BindRequest function")
	}

	if metadata.typ != typ {
		return errors.New("request object and type don't match")
	}

	body := ctx.fastHttpRequestContext.Request.Body()
	if metadata.hasOnlyBody {
		contentType, ok := ctx.GetRequestHeader(fasthttp.HeaderContentType)
		if !ok {
			contentType = MediaTypeApplicationJsonValue
		}

		if contentType == MediaTypeApplicationJsonValue {
			err := json.Unmarshal(body, request)
			if err != nil {
				return err
			}
		} else {
			err := xml.Unmarshal(body, request)
			if err != nil {
				return err
			}
		}
		return nil
	}

	val := reflect.ValueOf(request)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if metadata.bodyMetadata.fieldIndex != -1 {
		bodyValue := val.Field(metadata.bodyMetadata.fieldIndex)
		contentType, ok := ctx.GetRequestHeader(fasthttp.HeaderContentType)
		if !ok {
			contentType = MediaTypeApplicationJsonValue
		}

		if contentType == MediaTypeApplicationJsonValue {
			err := json.Unmarshal(body, bodyValue.Addr().Interface())
			if err != nil {
				return err
			}
		} else if contentType == MediaTypeApplicationXmlValue {
			err := xml.Unmarshal(body, bodyValue.Addr().Interface())
			if err != nil {
				return err
			}
		}
	}

	if metadata.paramMetadata.fieldIndex != -1 {
		paramStruct := val.Field(metadata.paramMetadata.fieldIndex)
		for tagValue, fieldMetadata := range metadata.paramMetadata.paramMap {
			paramField := paramStruct.Field(fieldMetadata.index)
			paramValue, ok := ctx.GetRequestParameter(tagValue)
			if !ok {
				continue
			}

			if fieldMetadata.converter != nil {
				paramField.Set(reflect.ValueOf(fieldMetadata.converter(paramValue)))
			} else {
				paramField.SetString(paramValue)
			}
		}
	}

	if metadata.pathMetadata.fieldIndex != -1 {
		pathStruct := val.Field(metadata.pathMetadata.fieldIndex)
		for _, fieldMetadata := range metadata.pathMetadata.pathVariableMap {
			pathField := pathStruct.Field(fieldMetadata.index)
			if fieldMetadata.extra == -1 {
				continue
			}

			pathVariableValue := ctx.pathVariables[fieldMetadata.extra]
			if fieldMetadata.converter != nil {
				pathField.Set(reflect.ValueOf(fieldMetadata.converter(pathVariableValue)))
			} else {
				pathField.SetString(pathVariableValue)
			}
		}
	}

	if metadata.headerMetadata.fieldIndex != -1 {
		headerStruct := val.Field(metadata.headerMetadata.fieldIndex)
		for tagValue, fieldMetadata := range metadata.headerMetadata.headerMap {
			headerField := headerStruct.Field(fieldMetadata.index)
			headerValue, ok := ctx.GetRequestHeader(tagValue)
			if !ok {
				continue
			}

			if fieldMetadata.converter != nil {
				headerField.Set(reflect.ValueOf(fieldMetadata.converter(headerValue)))
			} else {
				headerField.SetString(headerValue)
			}
		}
	}
	return nil
}
