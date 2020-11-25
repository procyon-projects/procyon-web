package web

type MediaType byte

const (
	DefaultMediaType                       = MediaTypeApplicationTextHtml
	MediaTypeApplicationTextHtml MediaType = 0
	MediaTypeApplicationJson     MediaType = 0
	MediaTypeApplicationXml      MediaType = 1
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
	SetStatus(status int) ResponseBodyBuilder
	SetBody(body interface{}) ResponseBodyBuilder
	SetContentType(mediaType MediaType) ResponseBodyBuilder
}

type Response interface {
	GetStatus() int
	GetBody() interface{}
	GetContentType() MediaType
}

type ResponseEntity struct {
	body        interface{}
	status      int
	contentType MediaType
}
