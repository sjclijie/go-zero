package template

var Insert = `
func (m *default{{.upperStartCamelObject}}Model) Insert(data *{{.upperStartCamelObject}}) (int64,error) {
	err := m.Create( data ).Error
	return data.Id, err
}
`

var InsertMethod = `Insert(data *{{.upperStartCamelObject}}) (int64,error)`
