package template

var Error = `package {{.pkg}}

import "github.com/sjclijie/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound
`
