package template

var (
	Imports = `import (
	"context"
	"database/sql"
	"github.com/sjclijie/go-zero/core/stores/redis"
	"github.com/tmsong/gorm"

	"hll-iam-server/rpc/permission/internal/helper/mysql"
		"hll-iam-server/rpc/permission/internal/model"


	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"github.com/sjclijie/go-zero/core/stores/cache"
	"github.com/sjclijie/go-zero/core/stores/sqlc"
	"github.com/sjclijie/go-zero/core/stores/sqlx"
	"github.com/sjclijie/go-zero/core/stringx"
	"github.com/sjclijie/go-zero/tools/goctl/model/sql/builderx"
)
`
	ImportsNoCache = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"github.com/sjclijie/go-zero/core/stores/sqlc"
	"github.com/sjclijie/go-zero/core/stores/sqlx"
	"github.com/sjclijie/go-zero/core/stringx"
	"github.com/sjclijie/go-zero/tools/goctl/model/sql/builderx"
)
`
)
