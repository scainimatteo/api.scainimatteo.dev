package vikunja

import (
	"encoding/json"
	"log"
	"net/http"

	"api.scainimatteo.dev/services"
)

type VikunjaService struct {
	Pushover *services.PushoverService
}

func (s VikunjaService) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var data any
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Errore JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Data: %v", data)

	w.WriteHeader(http.StatusOK)
}

func (s VikunjaService) ValidateWebhookSignature(signatureHeader string, rawBody []byte, secret string) (bool, error) {
	return true, nil
}
