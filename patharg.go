package chidoc

type PathArgType struct {
	Kind   string `json:"type"`
	Format string `json:"format"`
}

type PathArg struct {
	Name     string      `json:"name"`
	In       string      `json:"in"`
	Required bool        `json:"required"`
	Schema   PathArgType `json:"schema"`
}
