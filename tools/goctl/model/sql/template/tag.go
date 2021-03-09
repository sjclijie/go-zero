package template

var Tag = "`gorm:\"column:{{.field}}{{if .nullAble}};default:{{.defaultValue}}{{end}}\"`"
