package httpx

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/sjclijie/go-zero/core/logx"
	"github.com/stretchr/testify/assert"
)

type message struct {
	Name string `json:"name"`
}

func init() {
	logx.Disable()
}

func TestError(t *testing.T) {

	w1 := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	Error(&w1, errors.New("hello world"))
	assert.Equal(t, "{\"code\":400,\"msg\":\"hello world\",\"data\":{}}", w1.builder.String())

	w2 := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	e := NewResponseError(CodeOption(50000), MsgOption("this is msg"), DataOption(struct {
		Name string
		Age  int64
	}{Name: "lijie", Age: 40}))
	Error(&w2, e)
	assert.Equal(t, "{\"code\":50000,\"msg\":\"this is msg\",\"data\":{\"Name\":\"lijie\",\"Age\":40}}", w2.builder.String())

	w3 := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	e2 := NewResponseError(CodeOption(50000), MsgOption("this is msg"))
	Error(&w3, e2)
	assert.Equal(t, "{\"code\":50000,\"msg\":\"this is msg\",\"data\":{}}", w3.builder.String())

	w4 := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	e3 := NewResponseError(CodeOption(50000))
	Error(&w4, e3)
	assert.Equal(t, "{\"code\":50000,\"msg\":\"unknown error.\",\"data\":{}}", w4.builder.String())

	w5 := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	e4 := NewResponseError()
	Error(&w5, e4)
	assert.Equal(t, "{\"code\":400,\"msg\":\"unknown error.\",\"data\":{}}", w5.builder.String())

	w6 := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	Error(&w6, errors.New(""))
	assert.Equal(t, "{\"code\":400,\"msg\":\"request failed.\",\"data\":{}}", w6.builder.String())
}

func TestOk(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	Ok(&w)
	assert.Equal(t, http.StatusOK, w.code)
	assert.Equal(t, "{\"code\":0,\"msg\":\"request successes.\",\"data\":{}}", w.builder.String())
}

func TestOkJson(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	msg := message{Name: "anyone"}
	OkJson(&w, msg)
	assert.Equal(t, http.StatusOK, w.code)
	assert.Equal(t, "{\"code\":0,\"msg\":\"request successes.\",\"data\":{\"name\":\"anyone\"}}", w.builder.String())
}

func TestWriteJsonTimeout(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		timeout: true,
	}
	msg := message{Name: "anyone"}
	WriteJson(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteJsonLessWritten(t *testing.T) {
	w := tracedResponseWriter{
		headers:     make(map[string][]string),
		lessWritten: true,
	}
	msg := message{Name: "anyone"}
	WriteJson(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

type tracedResponseWriter struct {
	headers     map[string][]string
	builder     strings.Builder
	code        int
	lessWritten bool
	timeout     bool
}

func (w *tracedResponseWriter) Header() http.Header {
	return w.headers
}

func (w *tracedResponseWriter) Write(bytes []byte) (n int, err error) {
	if w.timeout {
		return 0, http.ErrHandlerTimeout
	}

	n, err = w.builder.Write(bytes)
	if w.lessWritten {
		n -= 1
	}
	return
}

func (w *tracedResponseWriter) WriteHeader(code int) {
	w.code = code
}
