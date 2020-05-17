package models

type Retorno struct {
	StatusRetorno string      `json:"status_retorno"`
	MsgErro       string      `json:"msg_erro,ominempty"`
	Retorno       interface{} `json:"retorno"`
}
