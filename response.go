package web

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
	AddHeader(key string, value string) ResponseHeaderBuilder
}

type ResponseBodyBuilder interface {
	ResponseHeaderBuilder
	SetResponseStatus(status int) ResponseBodyBuilder
	SetResponseBody(body interface{}) ResponseBodyBuilder
	SetResponseContentType(mediaType MediaType) ResponseBodyBuilder
}

type Response interface {
	GetResponseStatus() int
	GetResponseBody() interface{}
	GetResponseContentType() MediaType
}

type ResponseEntity struct {
	body        interface{}
	status      int
	contentType MediaType
}
