package models

type Vencimento struct {
	DataVencimento     string `json:"data_vencimento"`
	SaldoDevedor       string `json:"saldo_devedor"`
	PrazoDias          string `json:"prazo_dias"`
	ValorParcela       string `json:"valor_pacela"`
	ValorPrincipal     string `json:"valor_principal"`
	ValorJuros         string `json:"valor_juros"`
	NumeroDocumento    string `json:"numero_documento"`
	ValorTarifaParcela string `json:"valor_tarifa_parcela"`
	CapitalAmortizado  string `json:"capital_amortizado"`
}
