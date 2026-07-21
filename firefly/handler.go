package firefly

import (
	"bytes"
	"crypto/hmac"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

	myAccountName := "Conto Principale"
	today := time.Now()
	startDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	for i, record := range records {
		if len(record) < 3 {
			log.Printf("Riga %d ignorata: campi insufficienti", i+1)
			continue
		}

		title := record[0]
		category := record[1]
		amountStr := record[2]

		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			log.Printf("Errore conversione importo riga %d: %v", i+1, err)
			continue
		}

		txType := "withdrawal"
		source := myAccountName
		destination := title // Se non hai un conto spesa specifico, usa il titolo o un default come "Uscite Varie"

		// Se l'importo è negativo, è un'entrata (Deposit)
		if amountFloat < 0 {
			txType = "deposit"
			source = "Entrate Varie" // Conto di origine per le entrate (Revenue account)
			destination = myAccountName
		}

		// Firefly richiede sempre importi positivi nel JSON
		finalAmount := fmt.Sprintf("%.2f", math.Abs(amountFloat))

		// 3. Prepara il payload per Firefly
		txRequest := FireflyTxRequest{
			Transactions: []Transaction{
				{
					Type:            txType,
					Date:            startDate,
					Amount:          finalAmount,
					Description:     title,
					CategoryName:    category,
					SourceName:      source,
					DestinationName: destination,
				},
			},
		}

		// 4. Invia la richiesta all'API di Firefly
		err = s.createFireflyTransaction(txRequest)
		if err != nil {
			log.Printf("Errore creazione transazione '%s': %v", title, err)
			// Puoi decidere se interrompere o continuare con le altre righe
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Importazione completata"))
}

func (s FireflyService) createFireflyTransaction(txReq FireflyTxRequest) error {
	jsonData, err := json.Marshal(txReq)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/transactions", s.Config.Firefly.BaseURL)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
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

	if resp.StatusCode >= 300 {
		return fmt.Errorf("firefly api ha risposto con status %d", resp.StatusCode)
	}

	return nil
}
