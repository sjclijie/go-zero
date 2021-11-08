package httpx

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sjclijie/go-zero/core/mapping"
	"github.com/sjclijie/go-zero/rest/internal/context"
)

const (
	formKey           = "form"
	pathKey           = "path"
	emptyJson         = "{}"
	maxMemory         = 32 << 20 // 32MB
	maxBodyLen        = 8 << 20  // 8MB
	separator         = ";"
	tokensInAttribute = 2
)

var (
	formUnmarshaler = mapping.NewUnmarshaler(formKey, mapping.WithStringValues())
	pathUnmarshaler = mapping.NewUnmarshaler(pathKey, mapping.WithStringValues())
)

func Parse(r *http.Request, v interface{}) error {
	if err := ParsePath(r, v); err != nil {
		return err
	}

	if err := ParseForm(r, v); err != nil {
		return err
	}

	return ParseJsonBody(r, v)
}

// Parses the form request.
func ParseForm(r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}

	params := make(map[string]interface{}, len(r.Form))
	for name := range r.Form {
		formValue := r.Form.Get(name)
		if len(formValue) > 0 {
			params[name] = formValue
		}
	}
	if r.MultipartForm != nil {
		for key, fileHeaders := range r.MultipartForm.File {
			file := []byte{}
			for i := 0; i < len(fileHeaders); i++ {
				if fileHeaders[i] == nil {
					continue
				}
				f, _ := fileHeaders[i].Open()
				buf, _ := ioutil.ReadAll(f)
				file = append(file, buf...)
			}
			params[key] = file
		}
	}
	return formUnmarshaler.Unmarshal(params, v)
}

func ParseHeader(headerValue string) map[string]string {
	ret := make(map[string]string)
	fields := strings.Split(headerValue, separator)

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) == 0 {
			continue
		}

		kv := strings.SplitN(field, "=", tokensInAttribute)
		if len(kv) != tokensInAttribute {
			continue
		}

		ret[kv[0]] = kv[1]
	}

	return ret
}

// Parses the post request which contains json in body.
func ParseJsonBody(r *http.Request, v interface{}) error {
	var reader io.Reader
	if withJsonBody(r) {
		reader = io.LimitReader(r.Body, maxBodyLen)
	} else {
		reader = strings.NewReader(emptyJson)
	}

	return mapping.UnmarshalJsonReader(reader, v)
}

// Parses the symbols reside in url path.
// Like http://localhost/bag/:name
func ParsePath(r *http.Request, v interface{}) error {
	vars := context.Vars(r)
	m := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		m[k] = v
	}

	return pathUnmarshaler.Unmarshal(m, v)
}

func withJsonBody(r *http.Request) bool {
	return r.ContentLength > 0 && strings.Contains(r.Header.Get(ContentType), ApplicationJson)
}
