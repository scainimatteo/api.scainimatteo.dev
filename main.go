package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"api.scainimatteo.dev/firefly"
	grandigiochinigiorno "api.scainimatteo.dev/grandi-giochini-giorno"
	"api.scainimatteo.dev/outline"
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
	outlineService := outline.OutlineService{
		Config: config,
	}
	grandiGiochiniGiornoService := grandigiochinigiorno.GrandiGiochiniGiornoService{
		Config: config,
	}

	http.HandleFunc("/firefly/webhook", fireflyService.HandleWebhook)
	http.HandleFunc("/firefly/csv", fireflyService.HandleCSVImport)

	http.HandleFunc("/vikunja/reminder_webhook", vikunjaService.HandleReminderWebhook)
	http.HandleFunc("/vikunja/create_task_webhook", vikunjaService.HandleCreateTaskWebhook)
	http.HandleFunc("/vikunja/update_task_webhook", vikunjaService.HandleUpdateTaskWebhook)
	http.HandleFunc("/vikunja/complete_task/{id}", vikunjaService.CompleteTask)

	http.HandleFunc("/outline/{templateName}", outlineService.GetTemplate)

	http.HandleFunc("/grandi-giochini-giorno/parola-del-giorno/landing", handleDynamicRedirect)
	http.HandleFunc("/grandi-giochini-giorno/parola-del-giorno", grandiGiochiniGiornoService.GetParolaDelGiorno)
	http.HandleFunc("/grandi-giochini-giorno/bandiera-del-giorno/landing", handleDynamicRedirect)
	http.HandleFunc("/grandi-giochini-giorno/bandiera-del-giorno", grandiGiochiniGiornoService.GetBandieraDelGiorno)

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

func handleDynamicRedirect(w http.ResponseWriter, r *http.Request) {
	targetPath := path.Dir(r.URL.Path)

	if targetPath == "" {
		targetPath = "/"
	}

	if r.URL.RawQuery != "" {
		targetPath += "?" + r.URL.RawQuery
	}

	http.Redirect(w, r, targetPath, http.StatusTemporaryRedirect)
}
