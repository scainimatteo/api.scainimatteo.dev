package services

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

type PushoverService struct {
	User string
}

func (s *PushoverService) Send(title, message, token string) {
	apiURL := "https://api.pushover.net/1/messages.json"
	formData := url.Values{
		"token":   {token},
		"user":    {s.User},
		"message": {message},
		"title":   {title},
	}
	resp, err := http.PostForm(apiURL, formData)
	if err != nil {
		log.Printf("❌ Errore Pushover: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("✅ Notifica inviata!")
	} else {
		log.Printf("⚠️ Errore Pushover Status: %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		log.Printf("⚠️ Errore Pushover Body: %s", string(body))
	}
}
