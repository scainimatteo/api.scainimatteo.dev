package vikunja

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"api.scainimatteo.dev/services"
)

type VikunjaService struct {
	Config   services.Config
	Pushover *services.PushoverService
	Calendar *services.CalendarService
	DB       *sql.DB
}

func (s VikunjaService) HandleReminderWebhook(w http.ResponseWriter, r *http.Request) {
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

	task, err := s.GetTaskByID(data.Data.Task.ID)
	if err != nil {
		http.Error(w, "Errore ottenimento task", http.StatusInternalServerError)
		return
	}

	for _, label := range task.Labels {
		switch label.Title {
		case "Reminder":
			s.HandleReminder(task)
		case "Autocomplete":
			err := s.HandleAutocomplete(task)
			if err != nil {
				log.Printf("errore completamento task: %v", err)
			}
		}

	}

	w.WriteHeader(http.StatusOK)
}

func (s VikunjaService) HandleCreateTaskWebhook(w http.ResponseWriter, r *http.Request) {
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

	task, err := s.GetTaskByID(data.Data.Task.ID)
	if err != nil {
		http.Error(w, "Errore ottenimento task", http.StatusInternalServerError)
		return
	}

	for _, label := range task.Labels {
		switch label.Title {
		case "Calendar":
			s.HandleCalendar(task)
		}
	}

}

func (s VikunjaService) GetTaskByID(taskID int) (*Task, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", s.Config.Vikunja.BaseURL, taskID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("errore creazione richiesta: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.Config.Vikunja.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("errore durante la chiamata API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vikunja API ha risposto con status: %d", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("errore decodifica risposta: %v", err)
	}

	return &task, nil
}

func (s VikunjaService) ValidateWebhookSignature(signatureHeader string, rawBody []byte, secret string) (bool, error) {
	return true, nil
}

func (s VikunjaService) HandleCalendar(task *Task) error {
	var googleEventID string

	err := s.DB.QueryRow("SELECT calendar_id FROM calendar_event WHERE vikunja_id = $1", task.ID).Scan(&googleEventID)

	switch err {
	case sql.ErrNoRows:
		{
			newEventID, err := s.Calendar.UpsertEvent(context.Background(), "primary", "", task.Title, task.Description, task.DueDate, task.DueDate)
			if err != nil {
				return err
			}

			_, err = s.DB.Exec("INSERT INTO calendar_event (vikunja_id, calendar_id) VALUES ($1, $2)", task.ID, newEventID)
			if err != nil {
				return err
			}
		}
	case nil:
		{
			_, err := s.Calendar.UpsertEvent(context.Background(), "primary", googleEventID, task.Title, task.Description, task.DueDate, task.DueDate)
			if err != nil {
				return err
			}
		}
	default:
		{
			return err
		}
	}

	return nil
}

func (s VikunjaService) HandleReminder(task *Task) {
	title := "Reminder - " + task.Title
	message := fmt.Sprintf("https://vikunja.scainimatteo.dev/tasks/%d", task.ID)
	s.Pushover.Send(title, message, s.Config.VikunjaPushoverToken)
}

func (s VikunjaService) HandleAutocomplete(task *Task) error {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", s.Config.Vikunja.BaseURL, task.ID)

	payload := map[string]interface{}{
		"done": true,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.Config.Vikunja.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vikunja API ha risposto con status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("Task completato con successo | status: ", resp.StatusCode, " | body: ", string(body))

	return nil
}
