package outline

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"api.scainimatteo.dev/services"
)

type OutlineService struct {
	Config services.Config
}

//go:embed templates/calculator.html
var calculatorTemplate string

//go:embed templates/sum-list.html
var sumListTemplate string

//go:embed templates/copy-text.html
var copyTextTemplate string

func (s OutlineService) GetTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	templateName := r.PathValue("templateName")

	var template string

	switch templateName {
	case "transfer_calculator":
		template = calculatorTemplate
	case "sum_list":
		template = sumListTemplate
	case "copy_month_table":
		var jsonBytes []byte
		currentMonth := time.Now().Format("01/06")
		placeholderPlain := strings.ReplaceAll(monthTablePlaceholderPlain, "\"{MM/YY}\"", currentMonth)
		placeholderHTML := strings.ReplaceAll(monthTablePlaceholderHTML, "\"{MM/YY}\"", currentMonth)
		template = copyTextTemplate
		template = strings.ReplaceAll(template, "\"{placeholderTitle}\"", monthTablePlaceholderTitle)
		jsonBytes, _ = json.Marshal(placeholderPlain)
		template = strings.ReplaceAll(template, "\"{placeholderPlain}\"", string(jsonBytes))
		jsonBytes, _ = json.Marshal(placeholderHTML)
		template = strings.ReplaceAll(template, "\"{placeholderHTML}\"", string(jsonBytes))
	case "copy_1_1_odg":
		var jsonBytes []byte
		oneOnOneDate := nextOneOnOneWednesday(time.Now())
		oneOnOneDateText := oneOnOneDate.Format("01/02")
		placeholderPlain := strings.ReplaceAll(oneOnOneOdgPlaceholderPlain, "\"{MM/DD}\"", oneOnOneDateText)
		placeholderHTML := strings.ReplaceAll(oneOnOneOdgPlaceholderHTML, "\"{MM/DD}\"", oneOnOneDateText)
		template = copyTextTemplate
		template = strings.ReplaceAll(template, "\"{placeholderTitle}\"", oneOnOneOdgPlaceholderTitle)
		jsonBytes, _ = json.Marshal(placeholderPlain)
		template = strings.ReplaceAll(template, "\"{placeholderPlain}\"", string(jsonBytes))
		jsonBytes, _ = json.Marshal(placeholderHTML)
		template = strings.ReplaceAll(template, "\"{placeholderHTML}\"", string(jsonBytes))
	case "copy_1_1_mail":
		var jsonBytes []byte
		oneOnOneDate := nextOneOnOneWednesday(time.Now())
		oneOnOneDateText := oneOnOneDate.Format("02/01")
		placeholderPlain := strings.ReplaceAll(oneOnOneMailPlaceholderPlain, "\"{DD/MM}\"", oneOnOneDateText)
		placeholderHTML := strings.ReplaceAll(oneOnOneMailPlaceholderHTML, "\"{DD/MM}\"", oneOnOneDateText)
		template = copyTextTemplate
		template = strings.ReplaceAll(template, "\"{placeholderTitle}\"", oneOnOneMailPlaceholderTitle)
		jsonBytes, _ = json.Marshal(placeholderPlain)
		template = strings.ReplaceAll(template, "\"{placeholderPlain}\"", string(jsonBytes))
		jsonBytes, _ = json.Marshal(placeholderHTML)
		template = strings.ReplaceAll(template, "\"{placeholderHTML}\"", string(jsonBytes))
	case "attivita_settimanali":
		var jsonBytes []byte
		now := time.Now()
		daysUntilMonday := (8 - int(now.Weekday())) % 7
		monday := now.AddDate(0, 0, daysUntilMonday)
		friday := monday.AddDate(0, 0, 4)
		mondayText := monday.Format("02/01")
		fridayText := friday.Format("02/01")
		placeholderPlain := strings.ReplaceAll(weeklyActivitiesPlaceholderPlain, "\"{MONDAY}\"", mondayText)
		placeholderPlain = strings.ReplaceAll(placeholderPlain, "\"{FRIDAY}\"", fridayText)
		placeholderHTML := strings.ReplaceAll(weeklyActivitiesPlaceholderHTML, "\"{MONDAY}\"", mondayText)
		placeholderHTML = strings.ReplaceAll(placeholderHTML, "\"{FRIDAY}\"", fridayText)
		template = copyTextTemplate
		template = strings.ReplaceAll(template, "\"{placeholderTitle}\"", weeklyActivitiesPlaceholderTitle)
		jsonBytes, _ = json.Marshal(placeholderPlain)
		template = strings.ReplaceAll(template, "\"{placeholderPlain}\"", string(jsonBytes))
		jsonBytes, _ = json.Marshal(placeholderHTML)
		template = strings.ReplaceAll(template, "\"{placeholderHTML}\"", string(jsonBytes))
	default:
		http.Error(w, "Template non trovato", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(template))
}

// nextOneOnOneWednesday restituisce il primo mercoledì a partire da now (incluso)
// la cui settimana ISO ha la parità indicata da oneOnOneOddWeek.
func nextOneOnOneWednesday(now time.Time) time.Time {
	daysUntilWednesday := (int(time.Wednesday) - int(now.Weekday()) + 7) % 7
	wednesday := now.AddDate(0, 0, daysUntilWednesday)
	_, week := wednesday.ISOWeek()
	isOddWeek := week%2 == 1
	if isOddWeek != oneOnOneOddWeek {
		wednesday = wednesday.AddDate(0, 0, 7)
	}
	return wednesday
}
