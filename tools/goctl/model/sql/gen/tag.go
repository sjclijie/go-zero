package gen

import (
	"github.com/sjclijie/go-zero/tools/goctl/model/sql/parser"
	"github.com/sjclijie/go-zero/tools/goctl/model/sql/template"
	"github.com/sjclijie/go-zero/tools/goctl/util"
)

func genTag(in parser.Field) (string, error) {

	//不允许为null的时候，才需要默认值
	var defaultValue string
	switch in.Default.(type) {
	case []uint8:
		defaultValue = string(in.Default.([]uint8))
	default:
		defaultValue = ""
	}

	//datetime时，默认值为null
	if in.DataType == "*time.Time" || in.DataType == "time.Time" || in.DataType == "sql.NullTime" {
		defaultValue = "null"
		in.IsNullAble = false
	}

	if in.Name.Source() == "" {
		return in.Name.Source(), nil
	}

	text, err := util.LoadTemplate(category, tagTemplateFile, template.Tag)
	if err != nil {
		return "", err
	}

	output, err := util.With("tag").Parse(text).Execute(map[string]interface{}{
		"field":        in.Name.Source(),
		"nullAble":     !in.IsNullAble && defaultValue != "",
		"defaultValue": defaultValue,
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
