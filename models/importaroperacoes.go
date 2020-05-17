package models

import "time"

type ImportarOperacoes struct {
	NumeroCCB              int           `json:"numero_ccb"`
	DataAceite             time.Time     `json:"data_aceite"`
	DataLiberacao          time.Time     `json:"data_liberacao"`
	ValorTotal             string        `json:"valor_total"`
	ValorIof               string        `json:"valor_iof"`
	Taxt                   string        `json:"taxa"`
	Prazo                  string        `json:"prazo"`
	Modalidade             string        `json:"modalidade"`
	Periodicidade          string        `json:"periodicidade"`
	PercIof                string        `json:"parc_iof"`
	DataPrimeiroVencimento string        `json:"data_primeiro_vencimento"`
	ValorParcela           string        `json:"valor_parcela"`
	ValorTarifas           string        `json:"valor_tarifa"`
	ValorTitulo            string        `json:"valor_titulo"`
	Tac                    string        `json:"tac"`
	PercIofAdicional       string        `json:"perc_iof_adicional"`
	IDCliente              string        `json:"id_cliente"`
	CpfCnpjCliente         string        `json:"cpfcnpj_cliente"`
	RazaoSocial            string        `json:"razao_social"`
	EmailCliente           string        `json:"email_cliente"`
	EnderecoCliente        string        `json:"endereco_cliente"`
	BairroCliente          string        `json:"bairro_cliente"`
	CidadeCliente          string        `json:"cidade_cliente"`
	EstadoCliente          string        `json:"estado_cliente"`
	CepCliente             string        `json:"cep_cliente"`
	TelefoneCliente        string        `json:"telefone_cliente"`
	Receita01              string        `json:"receita_mes_01"`
	Receita02              string        `json:"receita_mes_02"`
	Receita03              string        `json:"receita_mes_03"`
	Receita04              string        `json:"receita_mes_04"`
	Receita05              string        `json:"receita_mes_05"`
	Receita06              string        `json:"receita_mes_06"`
	Receita07              string        `json:"receita_mes_07"`
	Receita08              string        `json:"receita_mes_08"`
	Receita09              string        `json:"receita_mes_09"`
	Receita10              string        `json:"receita_mes_10"`
	Receita11              string        `json:"receita_mes_11"`
	Receita12              string        `json:"receita_mes_12"`
	TipoDocumento          string        `json:"tipo_documento"`
	Escolaridade           string        `json:"escolaridade"`
	EstadoCivil            string        `json:"estado_civil"`
	TipoResidencia         string        `json:"tiporesidencia"`
	Inscestadual           string        `json:"inscestadual"`
	RG                     string        `json:"rg"`
	OrgaoRG                string        `json:"orgaorg"`
	UfOrgaoRG              string        `json:"uf_orgaorg"`
	Nacionalidade          string        `json:"nacionalidade"`
	Complemento            string        `json:"complemento"`
	Nascimento             string        `json:"nascimento"`
	Sexo                   string        `json:"sexo"`
	DtexpRG                string        `json:"dtexprg"`
	Naturalidade           string        `json:"naturalidade"`
	CpfConj                string        `json:"cpfconj"`
	NomeConjuge            string        `json:"nomeconjuge"`
	Dtnascconjuge          string        `json:"dtnascconjuge"`
	Mae                    string        `json:"mae"`
	Pai                    string        `json:"pai"`
	NumDependentes         string        `json:"numdependentes"`
	Fonecel                string        `json:"fonecel"`
	NumeroEndereco         string        `json:"numeroendereco"`
	TempoResidencia        string        `json:"temporesidencia"`
	ValorAluguel           string        `json:"vlaluguel"`
	TipoEndereco           string        `json:"tipo_endereco"`
	ContaBancaria          *Conta        `json:"conta_bancaria"`
	Avalistas              []*Avalista   `json:"avalistas"`
	Socios                 []*Socio      `json:"socios"`
	Vencimentos            []*Vencimento `json:"vencimentos"`
}
