package services

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarService struct {
	srv *calendar.Service
}

// NewCalendarService inizializza il client usando il percorso al file JSON della Service Account
func NewCalendarService(ctx context.Context, serviceAccountKeyPath string) (*CalendarService, error) {
	srv, err := calendar.NewService(ctx, option.WithAuthCredentialsFile(option.ServiceAccount, serviceAccountKeyPath))
	if err != nil {
		return nil, fmt.Errorf("errore inizializzazione Google Calendar: %v", err)
	}
	return &CalendarService{srv: srv}, nil
}

// UpsertEvent riceve solo dati primitivi, rendendolo agnostico rispetto a chi lo chiama
func (s *CalendarService) UpsertEvent(ctx context.Context, calendarID, eventID, title, description string, start, end time.Time) (string, error) {

	// Configurazione dell'evento
	event := &calendar.Event{
		Summary:     title,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: "UTC", // O la tua timezone locale
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: "UTC",
		},
	}

	if eventID == "" {
		// CREAZIONE: Se non abbiamo un ID, inseriamo un nuovo evento
		res, err := s.srv.Events.Insert(calendarID, event).Context(ctx).Do()
		if err != nil {
			return "", fmt.Errorf("errore creazione evento: %v", err)
		}
		return res.Id, nil
	}

	// AGGIORNAMENTO: Se l'ID esiste, aggiorniamo l'evento esistente
	res, err := s.srv.Events.Update(calendarID, eventID, event).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("errore aggiornamento evento %s: %v", eventID, err)
	}
	return res.Id, nil
}
