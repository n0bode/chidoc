package chidoc

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/n0bode/chi-doc/api"
	"github.com/n0bode/chi-doc/models"
)

func TestAddRouteDoc(t *testing.T) {
	route := chi.NewRouter()
	api := api.NewAPI()
	api.Route(route)

	AddRouteDoc(route, "/", DocSettings{
		"title": "Lucree CCB",
		"icon":  "",
		"logo":  "",
		"info": map[string]interface{}{
			"description": "Lucree CCB",
			"version":     "1.0",
		},
		"security": map[string]interface{}{
			"token": map[string]interface{}{
				"description": "## Gerar token\n\n" +
					"Toda a autorização é feita pelo cognito então, " +
					"para acessar a api é nessário que tenha um " +
					"usuario e senha para a geração do token.\n\n" +
					"O token é gerado na url:\n" +
					"```POST https://ccb.lucree.com.br/auth```\n\n" +
					"O escopo do body:\n" +
					"```json\n{\n\t\"username\":\"USUARIO\",\n\t\"password\":\"SENHA\"\n{\n```\n" +
					"Passe o token no header dos endpoints da API:\n" +
					"```bash\n$ POST https://ccc.lucree.com.br/endpoint --HEADER '{\"token\":\"TOKEN\"}'\n```",
				"type":  "http",
				"sheme": "bearer",
			},
		},
		"definitions": []interface{}{
			models.Retorno{},
			models.Documento{},
			models.ImportarOperacoes{},
			models.ConsultaPeriodo{},
			models.ConsultaOperacao{},
		},
	})

	http.ListenAndServe(":9000", route)
}
