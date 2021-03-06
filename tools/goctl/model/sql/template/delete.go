package template

var Delete = `
func (m *default{{.upperStartCamelObject}}Model) Delete( condition map[string]interface{} ) error {
	now := time.Now()
	_, err := m.Update(condition, &{{.upperStartCamelObject}}{
		Status:    model.StatusDeleted,
		DeletedAt: &now,
	})
	return err
}
`

var DeleteMethod = `Delete( condition map[string]interface{} ) error`
