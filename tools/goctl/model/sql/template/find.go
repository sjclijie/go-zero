package template

// 通过id查询
var FindOne = `

func (m *default{{.upperStartCamelObject}}Model) Query(condition map[string]interface{}) (*{{.upperStartCamelObject}}, error) {
	var ret {{.upperStartCamelObject}}
	err := m.Where(condition).Where(&{{.upperStartCamelObject}}{
		Status: model.StatusNormal,
	}).Find(&ret).Error
	return &ret, err
}
`

// 通过id查询
var FindList = `

func (m *default{{.upperStartCamelObject}}Model) QueryList(condition map[string]interface{}) ([]*{{.upperStartCamelObject}}, error) {
	var ret []*{{.upperStartCamelObject}}
	err := m.Where(condition).Where(&{{.upperStartCamelObject}}{
		Status: model.StatusNormal,
	}).Find(&ret).Error
	return ret, err
}
`

// 通过指定字段查询
var FindOneByField = `
func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}({{.in}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRowIndex(&resp, {{.cacheKeyVariable}}, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where {{.originalField}} = ? limit 1", {{.lowerStartCamelObject}}Rows, m.table)
		if err := conn.QueryRow(&resp, query, {{.lowerStartCamelField}}); err != nil {
			return nil, err
		}
		return resp.{{.upperStartCamelPrimaryKey}}, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{else}}var resp {{.upperStartCamelObject}}
	query := fmt.Sprintf("select %s from %s where {{.originalField}} = ? limit 1", {{.lowerStartCamelObject}}Rows, m.table )
	err := m.conn.QueryRow(&resp, query, {{.lowerStartCamelField}})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{end}}
`
var FindOneByFieldExtraMethod = `
func (m *default{{.upperStartCamelObject}}Model) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("#%s%v", {{.primaryKeyLeft}}, primary)
}
`

var FindOneMethod = `Query( condition map[string]interface{} ) ( *{{.upperStartCamelObject}}, error)`
var FindListMethod = `QueryList( condition map[string]interface{} ) ( []*{{.upperStartCamelObject}}, error)`
var FindOneByFieldMethod = `FindOneBy{{.upperField}}({{.in}}) (*{{.upperStartCamelObject}}, error) `
