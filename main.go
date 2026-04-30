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

var config services.Config

func main() {
	// 1. Carica la configurazione
	err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("❌ Errore caricamento config: %v", err)
	}

	pushover := services.PushoverService{
		User: config.PushoverUser,
	}
	fireflyService := firefly.FireflyService{
		Config:   config,
		Pushover: &pushover,
	}
	vikunjaService := vikunja.VikunjaService{
		Config:   config,
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
