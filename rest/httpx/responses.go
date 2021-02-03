package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/sjclijie/go-zero/core/logx"
)

var (
	errorHandler func(error) (int, interface{})
	lock         sync.RWMutex
)

const (
	RequestSuccessCode = 0
	RequestSuccessMsg  = "request successes."

	RequestBadCode = 400
	RequestBadMsg  = "request failed."
)

type ret struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newRet() *ret {
	return &ret{
		Code: RequestSuccessCode,
		Msg:  RequestSuccessMsg,
		Data: &struct{}{},
	}
}

func (r *ret) wrapRet(v interface{}) *ret {

	r.Code = RequestSuccessCode
	r.Msg = RequestSuccessMsg
	r.Data = v

	return r
}

func (r *ret) wrapErrRet(err error) *ret {
	if ok := errors.Is(err, &ResponseError{}); ok {
		e, _ := err.(*ResponseError)
		r.Code = e.Code
		r.Msg = e.Msg
		r.Data = e.Data

		if r.Code == 0 {
			r.Code = RequestBadCode
		}

	} else {
		r.Code = RequestBadCode
		r.Msg = err.Error()
		r.Data = &struct{}{}

		if r.Msg == "" {
			r.Msg = RequestBadMsg
		}
	}

	return r
}

func Error(w http.ResponseWriter, err error) {
	WriteJson(w, http.StatusOK, newRet().wrapErrRet(err))
}

func Ok(w http.ResponseWriter) {
	WriteJson(w, http.StatusOK, newRet())
}

func OkJson(w http.ResponseWriter, v interface{}) {
	WriteJson(w, http.StatusOK, newRet().wrapRet(v))
}

func SetErrorHandler(handler func(error) (int, interface{})) {
	lock.Lock()
	defer lock.Unlock()
	errorHandler = handler
}

func WriteJson(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(code)

	if bs, err := json.Marshal(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if n, err := w.Write(bs); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if err != http.ErrHandlerTimeout {
			logx.Errorf("write response failed, error: %s", err)
		}
	} else if n < len(bs) {
		logx.Errorf("actual bytes: %d, written bytes: %d", len(bs), n)
	}
}
