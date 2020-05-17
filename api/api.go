package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

// API estrutura que guarda as variaveis da API
type API struct{}

// NewAPI cria a estrutura da API
func NewAPI() *API {
	return &API{}
}

// Route seta as rotas da api
func (api *API) Route(root *chi.Mux) {
	root.Route("/bancarizacao", func(r chi.Router) {
		r.Get("/", api.GetBancarizacao)
		r.Post("/", api.PostBancarizacao)
	})
}

// GetBancarizacao retorna
// summary: Bancariazacao
// tags:
// - bancarizacao
// security:
// - token: []
// requestBody:
//   description: Insira uns destes formatos
//   required: true
//   content:
//    application/json:
//     schema:
//      anyOf:
//      - ConsultaOperacao:
//        type: object
//        properties:
//         consultaroperacao:
//          "$ref": "#/components/schemes/ConsultaOperacao"
//      - ConsultaPerido:
//        type: object
//        properties:
//         consultaperiodo:
//          "$ref": "#/components/schemes/ConsultaPeriodo"
// responses:
//  '200':
//    description: Sucesso
func (api *API) GetBancarizacao(w http.ResponseWriter, r *http.Request) {

}

// PostBancarizacao posta uma nova
// summary: Bancariazacao
// tags:
//  - bancarizacao
// security:
//  - token: []
// requestBody:
//  description: Insira uns destes formatos
//  required: true
//  content:
//   application/json:
//    schema:
//     anyOf:
//      - "$ref": "#/components/schemes/Documento"
//      - "$ref": "#/components/schemes/ImportarOperacoes"
// responses:
//  '201':
//    description: Sucesso
//    content:
//     application/json:
//      schema:
//       "$ref": "#/components/schemes/Retorno"
func (api *API) PostBancarizacao(w http.ResponseWriter, r *http.Request) {

}
