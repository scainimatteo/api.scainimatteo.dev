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

func (s *CalendarService) UpsertEvent(ctx context.Context, calendarID, eventID, title, description string, start time.Time) (string, error) {
	timezone, _ := time.LoadLocation("Europe/Rome")
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, timezone)

	event := &calendar.Event{
		Summary:     title,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startDate.Format(time.RFC3339),
			TimeZone: timezone.String(),
		},
		End: &calendar.EventDateTime{
			DateTime: startDate.Format(time.RFC3339),
			TimeZone: timezone.String(),
		},
	}

	if eventID == "" {
		res, err := s.srv.Events.Insert(calendarID, event).Context(ctx).Do()
		if err != nil {
			return "", fmt.Errorf("errore creazione evento: %v", err)
		}
		return res.Id, nil
	}

	res, err := s.srv.Events.Update(calendarID, eventID, event).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("errore aggiornamento evento %s: %v", eventID, err)
	}
	return res.Id, nil
}
