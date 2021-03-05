package template

var New = `
func New{{.upperStartCamelObject}}Model(ctx context.Context, redis *redis.Redis, readOnly bool) ({{.upperStartCamelObject}}Model, error )  {

	model := &default{{.upperStartCamelObject}}Model{
		BaseModel: model.BaseModel{Ctx: ctx, RedisClient: redis},
		table:     {{.lowerStartCamelObject}}TableName,
	}

	var err error
	if readOnly {
		model.Gorm, err = mysql.GetIamDbReadOnly()
	} else {
		model.Gorm, err = mysql.GetIamDb()
	}

	model.DB = model.Gorm.Table(model.table).Model(&default{{.upperStartCamelObject}}Model{})

	return model, err
}

func ({{.upperStartCamelObject}}) TableName() string {
	return {{.lowerStartCamelObject}}TableName
}

func (m *default{{.upperStartCamelObject}}Model) TableName() string {
	return m.table
}
`
