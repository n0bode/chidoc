package models

type Conta struct {
	CodBanco      string `json:"cod_banco"`
	Conta         string `json:"conta"`
	Agencia       string `json:"agencia"`
	AgenciaD      string `json:"agenciad"`
	Operacao      string `json:"operacao"`
	Conjunta      string `json:"conjunta"`
	Tipo          string `json:"tipo"`
	CpfFavorecido string `json:"cpf_favorecido"`
}
