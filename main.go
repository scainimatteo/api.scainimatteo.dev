package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"api.scainimatteo.dev/firefly"
	"api.scainimatteo.dev/services"
	"api.scainimatteo.dev/vikunja"
)

// Config definisce la struttura del file config.json
type Config struct {
	PushoverToken string `json:"pushover_token"`
	PushoverUser  string `json:"pushover_user"`
	Port          string `json:"port"`
}

var config Config

func main() {
	// 1. Carica la configurazione
	err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("❌ Errore caricamento config: %v", err)
	}

	pushover := services.PushoverService{
		Token: config.PushoverToken,
		User:  config.PushoverUser,
	}
	fireflyService := firefly.FireflyService{
		Pushover: &pushover,
	}
	vikunjaService := vikunja.VikunjaService{
		Pushover: &pushover,
	}

	// 2. Definisci le rotte
	http.HandleFunc("/firefly/webhook", fireflyService.HandleWebhook)
	http.HandleFunc("/vikunja/webhook", vikunjaService.HandleWebhook)

	fmt.Printf("🚀 Server in ascolto sulla porta %s...\n", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

// loadConfig legge il file JSON e popola la variabile globale config
func loadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	return json.Unmarshal(byteValue, &config)
}
