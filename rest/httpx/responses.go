package httpx

import (
	"encoding/json"
	"errors"
	"google.golang.org/grpc/status"
	"net/http"
	"sync"

	"github.com/sjclijie/go-zero/core/logx"
)

var (
	errorHandler func(error) (int, interface{})
	lock         sync.RWMutex
)

var GrpcCodeMap = map[string]int64{
	"OK":                 10000,
	"Canceled":           11000,
	"Unknown":            12000,
	"InvalidArgument":    13000,
	"DeadlineExceeded":   14000,
	"NotFound":           15000,
	"AlreadyExists":      16000,
	"PermissionDenied":   17000,
	"ResourceExhausted":  18000,
	"FailedPrecondition": 19000,
	"Aborted":            20000,
	"OutOfRange":         21000,
	"Unimplemented":      22000,
	"Internal":           23000,
	"Unavailable":        24000,
	"DataLoss":           25000,
	"Unauthenticated":    26000,
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
		if ret, ok := GrpcCodeMap[st.Code().String()]; ok {
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
