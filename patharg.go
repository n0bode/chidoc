package chidoc

// PathArgType struct to define type of args documentation route
type PathArgType struct {
	Kind   string `json:"type"`
	Format string `json:"format"`
}

// PathArg struct to define args documentation route
type PathArg struct {
	Name     string      `json:"name"`
	In       string      `json:"in"`
	Required bool        `json:"required"`
	Schema   PathArgType `json:"schema"`
}
