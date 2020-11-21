package web

type MediaType string

const (
	DefaultMediaType                   = MediaTypeApplicationJson
	MediaTypeApplicationXml  MediaType = "application/xml"
	MediaTypeApplicationJson MediaType = "application/json"
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
