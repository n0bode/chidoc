package chidoc

type StructEnum struct {
	Name        string
	Type        string
	Enum        []interface{}
	Description string
}

// Enum generate models for enum fields
func Enum(name, description string, values ...interface{}) StructEnum {
	return StructEnum{
		Name:        name,
		Enum:        values,
		Type:        typeName(values[0]),
		Description: description,
	}
}

func (s StructEnum) Parse() map[string]interface{} {
	m := make(map[string]interface{})
	m["type"] = s.Type
	m["enum"] = s.Enum
	m["description"] = s.Description
	return m
}
