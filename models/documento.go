package models

type Documento struct {
	TipoDocumento string `json:"tipo_documento"`
	Categoria     string `json:"categoria"`
	Operacao      string `json:"operacao"`
	Base64        string `json:"base64"`
}
