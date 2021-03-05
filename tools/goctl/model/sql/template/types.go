package template

var Types = `
type (
	{{.upperStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
		model.BaseModel
		*gorm.DB
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
`
