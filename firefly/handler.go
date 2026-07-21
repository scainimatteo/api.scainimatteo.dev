package firefly

import (
	"bytes"
	"crypto/hmac"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"api.scainimatteo.dev/services"
	"golang.org/x/crypto/sha3"
)

type FireflyService struct {
	Config   services.Config
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
		s.Pushover.Send(title, message, s.Config.Firefly.PushoverToken)
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

func (s FireflyService) HandleCSVImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	reader := csv.NewReader(r.Body)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Errore lettura CSV: %v", err)
		http.Error(w, "Errore nel formato CSV", http.StatusBadRequest)
		return
	}

	date := time.Now().Format("2006-01-02")

	for i, record := range records {
		if len(record) < 3 {
			log.Printf("Riga %d ignorata: campi insufficienti", i+1)
			continue
		}

		title := record[0]
		category := record[1]
		amount := record[2]

		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			log.Printf("Errore conversione importo riga %d: %v", i+1, err)
			continue
		}

		transaction := Transaction{
			Date:         date,
			Amount:       fmt.Sprintf("%.2f", math.Abs(amountFloat)),
			Description:  title,
			CategoryName: category,
		}

		if amountFloat < 0 {
			transaction.Type = "deposit"
			transaction.DestinationID = s.Config.Firefly.Sources.Bper
			transaction.SourceName = title
		} else {
			transaction.Type = "withdrawal"
			transaction.SourceID = s.Config.Firefly.Sources.Bper
		}

		err = s.saveTransaction(transaction)
		if err != nil {
			log.Printf("Errore creazione transazione '%s': %v", title, err)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Importazione completata"))
}

func (s FireflyService) saveTransaction(transaction Transaction) error {
	baseURL := strings.TrimSuffix(s.Config.Firefly.BaseURL, "/")
	apiURL := baseURL + "/api/v1/transactions"

	payload := Payload{
		Transactions: []Transaction{transaction},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.Config.Firefly.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status: %s, body: %s", resp.Status, string(b))
	}

	return nil
}
