package web

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWebRequestContext_prepare(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)

	ctx.prepare(false)
	assert.Equal(t, 0, len(ctx.contextIdStr))

	ctx.prepare(true)
	assert.NotNil(t, 0, len(ctx.contextIdStr))
}

func TestWebRequestContext_reset(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)

	ctx.handlerIndex = 1
	ctx.pathVariableCount = 1
	ctx.valueMap = make(map[string]interface{})
	ctx.responseEntity.status = http.StatusCreated
	ctx.responseEntity.body = "test-body"
	ctx.responseEntity.contentType = MediaTypeApplicationJson

	ctx.reset()

	assert.Equal(t, 0, ctx.handlerIndex)
	assert.Equal(t, 0, ctx.pathVariableCount)
	assert.Nil(t, ctx.valueMap)
	assert.Equal(t, http.StatusOK, ctx.responseEntity.status)
	assert.Nil(t, ctx.responseEntity.body)
	assert.Equal(t, DefaultMediaType, ctx.responseEntity.contentType)
}

func TestWebRequestContext_ValueMap(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.Put("test-key", "test-value")
	assert.Equal(t, "test-value", ctx.Get("test-key"))
}

func TestWebRequestContext_Status(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.SetStatus(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, ctx.GetStatus())
}

func TestWebRequestContext_Body(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.SetBody("test-body")
	assert.Equal(t, "test-body", ctx.GetBody())
}

func TestWebRequestContext_ContextType(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.SetContentType(MediaTypeApplicationJson)
	assert.Equal(t, MediaTypeApplicationJson, ctx.GetContentType())
}

func TestWebRequestContext_Ok(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.Ok()
	assert.Equal(t, http.StatusOK, ctx.responseEntity.status)
}

func TestWebRequestContext_NotFound(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.NotFound()
	assert.Equal(t, http.StatusNotFound, ctx.responseEntity.status)
}

func TestWebRequestContext_NoContent(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.NoContent()
	assert.Equal(t, http.StatusNoContent, ctx.responseEntity.status)
}

func TestWebRequestContext_BadRequest(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.BadRequest()
	assert.Equal(t, http.StatusBadRequest, ctx.responseEntity.status)
}

func TestWebRequestContext_Accepted(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.Accepted()
	assert.Equal(t, http.StatusAccepted, ctx.responseEntity.status)
}

func TestWebRequestContext_Created(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.Created("")
	assert.Equal(t, http.StatusCreated, ctx.responseEntity.status)
}

func TestWebRequestContext_Error(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.SetError(errors.New("test-error"))
	assert.Equal(t, "test-error", ctx.GetError().Error())
}

func TestWebRequestContext_ThrowError(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	assert.Panics(t, func() {
		ctx.ThrowError(errors.New("test-error"))
	})
}