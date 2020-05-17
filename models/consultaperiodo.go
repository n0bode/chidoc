package models

import "time"

type ConsultaPeriodo struct {
	DataInicio time.Time `json:"data_inicio" docs:"required"`
	DataFinal  time.Time `json:"data_final" docs:"required"`
}
