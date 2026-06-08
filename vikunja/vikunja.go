package vikunja

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

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

	s.HandleReminder(task)

	for _, label := range task.Labels {
		switch label.Title {
		case "Autocomplete":
			err := s.completeTask(task)
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

func (s VikunjaService) HandleUpdateTaskWebhook(w http.ResponseWriter, r *http.Request) {
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
		case "Fix hours":
			err := s.fixHours(task)
			if err != nil {
				log.Printf("errore fix hours: %v", err)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
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
			newEventID, err := s.Calendar.UpsertEvent(context.Background(), s.Config.Vikunja.CalendarID, "", task.Title, task.Description, task.DueDate)
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
			_, err := s.Calendar.UpsertEvent(context.Background(), s.Config.Vikunja.CalendarID, googleEventID, task.Title, task.Description, task.DueDate)
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

func (s VikunjaService) CompleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Errore conversione ID", http.StatusBadRequest)
		return
	}
	task, err := s.GetTaskByID(id)
	if err != nil {
		http.Error(w, "Errore ottenimento task", http.StatusInternalServerError)
		return
	}

	err = s.completeTask(task)
	if err != nil {
		http.Error(w, "Errore completamento task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task " + task.Title + " completato"))
}

func (s VikunjaService) completeTask(task *Task) error {
	payload := map[string]any{
		"done": true,
	}

	resp, err := s.callVikunjaAPI(task, "POST", payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vikunja API ha risposto con status %d", resp.StatusCode)
	}

	return nil
}

func (s VikunjaService) fixHours(task *Task) error {
	baseDate := task.DueDate
	if baseDate.IsZero() {
		baseDate = time.Now()
	}

	re := regexp.MustCompile(`^Hour:\s*([0-9]{2})[:.]([0-9]{2})`)

	matches := re.FindStringSubmatch(task.Description)
	if len(matches) < 2 {
		return fmt.Errorf("formato orario non trovato all'inizio della stringa")
	}

	correctHours, _ := strconv.Atoi(matches[1])
	correctMinutes, _ := strconv.Atoi(matches[2])
	newDueDate := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), correctHours, correctMinutes, 0, 0, baseDate.Location())

	payload := map[string]any{
		"due_date": newDueDate.Format(time.RFC3339),
	}

	resp, err := s.callVikunjaAPI(task, "POST", payload)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vikunja API ha risposto con status %d", resp.StatusCode)
	}

	return nil
}

func (s VikunjaService) callVikunjaAPI(task *Task, method string, payload map[string]any) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", s.Config.Vikunja.BaseURL, task.ID)

	if task.RepeatAfter > 0 {
		payload["repeat_after"] = task.RepeatAfter
		payload["repeat_mode"] = task.RepeatMode
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.Config.Vikunja.APIToken)
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}
