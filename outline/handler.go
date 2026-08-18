package outline

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"

	"api.scainimatteo.dev/services"
)

type OutlineService struct {
	Config services.Config
}

//go:embed templates/calculator.html
var calculatorTemplate string

//go:embed templates/sum-list.html
var sumListTemplate string

//go:embed templates/copy-text.html
var copyTextTemplate string
var monthTablePlaceholderTitle = "Copia template del mese negli appunti"
var monthTablePlaceholderPlain = `Descrizione
Categoria
Importo
Eseguito
Bolletta internet MM/YY
Bollette
50
`
var monthTablePlaceholderHTML = `<table>
  <colgroup>
    <col style="width: 421px;">
    <col style="width: 123px;">
    <col style="width: 96px;">
    <col style="width: 96px;">
  </colgroup>
  <tbody>
    <tr>
      <th data-colwidth="421"><p dir="auto">Descrizione</p></th>
      <th data-colwidth="123"><p dir="auto">Categoria</p></th>
      <th data-colwidth="96"><p dir="auto">Importo</p></th>
      <th data-colwidth="96"><p dir="auto">Eseguito</p></th>
    </tr>
    <tr>
      <td data-colwidth="421"><p dir="auto">Bolletta internet 08/26</p></td>
      <td data-colwidth="123"><p dir="auto">Bollette</p></td>
      <td data-colwidth="96"><p dir="auto">50</p></td>
      <td data-colwidth="96"><p dir="auto"></p></td>
    </tr>
  </tbody>
</table>`

func (s OutlineService) GetTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	templateName := r.PathValue("templateName")

	var template string

	switch templateName {
	case "transfer_calculator":
		template = calculatorTemplate
	case "sum_list":
		template = sumListTemplate
	case "copy_month_table":
		template = copyTextTemplate
		jsonBytes, _ := json.Marshal(monthTablePlaceholderTitle)
		template = strings.ReplaceAll(template, "\"{placeholderTitle}\"", string(jsonBytes))
		jsonBytes, _ = json.Marshal(monthTablePlaceholderPlain)
		template = strings.ReplaceAll(template, "\"{placeholderPlain}\"", string(jsonBytes))
		jsonBytes, _ = json.Marshal(monthTablePlaceholderHTML)
		template = strings.ReplaceAll(template, "\"{placeholderHTML}\"", string(jsonBytes))
	default:
		http.Error(w, "Template non trovato", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(template))
}
