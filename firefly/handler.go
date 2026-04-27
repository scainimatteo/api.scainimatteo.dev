package firefly

import (
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"api.scainimatteo.dev/services"
	"golang.org/x/crypto/sha3"
)

type FireflyService struct {
	Pushover *services.PushoverService
}

func (s FireflyService) HandleWebhook(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var data FireflyWebhookResponse
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Print(err)
		http.Error(w, "Errore JSON", http.StatusBadRequest)
		return
	}

	if len(data.Content.Transactions) > 0 && data.Content.Transactions[0].Description == "Rata macchina" {
		transaction := data.Content.Transactions[0]
		title := "Firefly III - Pagamento rata salvato"
		message := fmt.Sprintf("Descrizione: %s | Cifra: %s | Id: %v", transaction.Description, transaction.Amount, transaction.RecurrenceID)
		s.Pushover.Send(title, message)
	}

	w.WriteHeader(http.StatusOK)
}

func (s FireflyService) ValidateWebhookSignature(signatureHeader string, rawBody []byte, secret string) (bool, error) {
	var timestamp, v1Signature string

	parts := strings.SplitSeq(signatureHeader, ",")
	for part := range parts {
		if after, ok := strings.CutPrefix(part, "t="); ok {
			timestamp = after
		} else if after, ok := strings.CutPrefix(part, "v1="); ok {
			v1Signature = after
		}
	}

	if timestamp == "" || v1Signature == "" {
		return false, errors.New("invalid_signature")
	}

	providedMAC, err := hex.DecodeString(v1Signature)
	if err != nil {
		return false, fmt.Errorf("decode_error")
	}

	signedString := fmt.Sprintf("%s.%s", timestamp, string(rawBody))

	mac := hmac.New(sha3.New256, []byte(secret))
	mac.Write([]byte(signedString))
	expectedMAC := mac.Sum(nil)

	if hmac.Equal(providedMAC, expectedMAC) {
		return true, nil
	}

	return false, nil
}
