package httpx

type ResponseError struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResponseErrorOption func(responseError *ResponseError)

func CodeOption(code int64) ResponseErrorOption {
	return func(responseError *ResponseError) {
		responseError.Code = code
	}
}

func MsgOption(msg string) ResponseErrorOption {
	return func(responseError *ResponseError) {
		responseError.Msg = msg
	}
}

func DataOption(data interface{}) ResponseErrorOption {
	return func(responseError *ResponseError) {
		responseError.Data = data
	}
}

func NewResponseError(option ...ResponseErrorOption) *ResponseError {
	err := &ResponseError{Data: &struct{}{}}
	for _, opt := range option {
		opt(err)
	}

	if err.Msg == "" {
		err.Msg = "unknown error."
	}

	return err
}

func (e ResponseError) Error() string {
	return e.Msg
}

func (e *ResponseError) Is(target error) bool {
	_, ok := target.(*ResponseError)
	return ok
}
