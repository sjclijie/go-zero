package httpx

import (
	"encoding/json"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"sync"

	"github.com/sjclijie/go-zero/core/logx"
)

var (
	errorHandler func(error) (int, interface{})
	lock         sync.RWMutex
)

var GrpcCodeMap = map[codes.Code]int64{
	codes.OK:                 10000,
	codes.Canceled:           11000,
	codes.Unknown:            12000,
	codes.InvalidArgument:    13000,
	codes.DeadlineExceeded:   14000,
	codes.NotFound:           15000,
	codes.AlreadyExists:      16000,
	codes.PermissionDenied:   17000,
	codes.ResourceExhausted:  18000,
	codes.FailedPrecondition: 19000,
	codes.Aborted:            20000,
	codes.OutOfRange:         21000,
	codes.Unimplemented:      22000,
	codes.Internal:           23000,
	codes.Unavailable:        24000,
	codes.DataLoss:           25000,
	codes.Unauthenticated:    26000,
}

const (
	RequestSuccessCode = 0
	RequestSuccessMsg  = "request successes."

	RequestBadCode = 400
	RequestBadMsg  = "request failed."
)

type ret struct {
	Ret  int64       `json:"ret"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newRet() *ret {
	return &ret{
		Ret:  RequestSuccessCode,
		Msg:  RequestSuccessMsg,
		Data: &struct{}{},
	}
}

func (r *ret) wrapRet(v interface{}) *ret {
	r.Ret = RequestSuccessCode
	r.Msg = RequestSuccessMsg
	r.Data = v
	return r
}

func (r *ret) wrapErrRet(err error) *ret {
	if st, ok := status.FromError(err); ok {
		if ret, ok := GrpcCodeMap[st.Code()]; ok {
			r.Ret = ret
			r.Msg = st.Message()
		} else {
			r.Ret = RequestBadCode
			r.Msg = err.Error()
			r.Data = &struct{}{}
		}
	} else if ok := errors.Is(err, &ResponseError{}); ok {
		e, _ := err.(*ResponseError)
		r.Ret = e.Code
		r.Msg = e.Msg
		r.Data = e.Data
	} else {
		r.Ret = RequestBadCode
		r.Msg = err.Error()
		r.Data = &struct{}{}
	}

	if r.Msg == "" {
		r.Msg = RequestBadMsg
	}

	return r
}

func Error(w http.ResponseWriter, err error) {
	WriteJson(w, http.StatusOK, newRet().wrapErrRet(err))
}

func ErrorJson(w http.ResponseWriter, err error, v interface{}) {
	ret := newRet().wrapErrRet(err)
	ret.Data = v
	WriteJson(w, http.StatusOK, ret)
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
