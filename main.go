package main

import (
	"context"
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
	err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("❌ Errore caricamento config: %v", err)
	}

	db, err := services.NewDatabaseConnection(config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	calendarService, err := services.NewCalendarService(context.Background(), "google-calendar-key.json")
	if err != nil {
		log.Fatalf("❌ Errore inizializzazione Google Calendar: %v", err)
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
		Calendar: calendarService,
		DB:       db,
	}

	// 2. Definisci le rotte
	http.HandleFunc("/firefly/webhook", fireflyService.HandleWebhook)
	http.HandleFunc("/vikunja/reminder_webhook", vikunjaService.HandleReminderWebhook)
	http.HandleFunc("/vikunja/create_task_webhook", vikunjaService.HandleCreateTaskWebhook)
	http.HandleFunc("/vikunja/update_task_webhook", vikunjaService.HandleUpdateTaskWebhook)
	http.HandleFunc("/vikunja/complete_task/{id}", vikunjaService.CompleteTask)

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
