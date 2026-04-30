package vikunja

import (
	"encoding/json"
	"net/http"

	"api.scainimatteo.dev/services"
)

type VikunjaService struct {
	Config   services.Config
	Pushover *services.PushoverService
}

func (s VikunjaService) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var data VikunjaWebhookResponse
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Errore JSON", http.StatusBadRequest)
		return
	}

	title := "Reminder - " + data.Data.Task.Title
	message := "https://vikunja.scainimatteo.dev/projects/1/1"
	s.Pushover.Send(title, message, s.Config.VikunjaPushoverToken)

	w.WriteHeader(http.StatusOK)
}

func (s VikunjaService) ValidateWebhookSignature(signatureHeader string, rawBody []byte, secret string) (bool, error) {
	return true, nil
}
