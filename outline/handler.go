package outline

import (
	_ "embed"
	"net/http"

	"api.scainimatteo.dev/services"
)

type OutlineService struct {
	Config services.Config
}

//go:embed calculator.html
var calculatorTemplate string

//go:embed sum-list.html
var sumListTemplate string

func (s OutlineService) GetTransferCalculatorTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(calculatorTemplate))
}

func (s OutlineService) GetSumListTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sumListTemplate))
}
