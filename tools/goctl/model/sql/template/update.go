package template

var Update = `
func (m *default{{.upperStartCamelObject}}Model) Update( condition map[string]interface{}, data *{{.upperStartCamelObject}}) ( int64, error ) {
	db := m.Table(m.table).Where(condition).Update(data)
	return db.RowsAffected, db.Error
}
`

var UpdateMethod = `Update(condition map[string]interface{}, data *{{.upperStartCamelObject}})(int64, error)`
