package grandigiochinigiorno

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"text/template"
	"time"

	"api.scainimatteo.dev/services"
)

//go:embed templates/parola-giorno.html
var parolaDelGiornoTemplate string

//go:embed templates/bandiera-giorno.html
var bandieraDelGiornoTemplate string

type GrandiGiochiniGiornoService struct {
	Config services.Config
}

type datiBandieraDelGiorno struct {
	Nome     string
	Bandiera string
}

func (s GrandiGiochiniGiornoService) GetParolaDelGiorno(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.New("parola-del-giorno").Parse(parolaDelGiornoTemplate)
	if err != nil {
		log.Fatalf("Errore nel parsing del template HTML: %v", err)
	}

	word, err := s.getRandomWord(5)
	if err != nil {
		http.Error(w, "Errore durante la ricerca nel dizionario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Sanificazione dell'output per l'HTML
	word = html.EscapeString(word)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	data := struct{ Parola string }{Parola: word}

	err = tmpl.ExecuteTemplate(w, "parola-del-giorno", data)
	if err != nil {
		log.Printf("Errore nel rendering del template: %v", err)
	}
}

func (s GrandiGiochiniGiornoService) GetBandieraDelGiorno(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.New("bandiera-del-giorno").Parse(bandieraDelGiornoTemplate)
	if err != nil {
		log.Fatalf("Errore nel parsing del template HTML: %v", err)
	}

	paese, err := s.getRandomCountry()
	if err != nil {
		http.Error(w, "Errore durante la ricerca del paese: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	data := datiBandieraDelGiorno{
		Nome:     html.EscapeString(paese.Names.Common),
		Bandiera: html.EscapeString(paese.Flag.URLSvg),
	}

	err = tmpl.ExecuteTemplate(w, "bandiera-del-giorno", data)
	if err != nil {
		log.Printf("Errore nel rendering del template: %v", err)
	}
}

// getRandomCountry recupera un singolo paese a caso all'offset casuale, senza
// interrogare prima l'API per il totale: il totale dei paesi è fisso in
// config (RestCountriesTotal), così basta una sola chiamata a REST Countries.
func (s GrandiGiochiniGiornoService) getRandomCountry() (restCountry, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	offset := r.Intn(s.Config.GrandiGiochiniGiorno.RestCountriesTotal)

	paesi, _, err := s.fetchCountriesPage(offset, 1)
	if err != nil {
		return restCountry{}, err
	}

	if len(paesi) == 0 {
		return restCountry{}, fmt.Errorf("nessun paese trovato all'offset %d", offset)
	}

	return paesi[0], nil
}

func (s GrandiGiochiniGiornoService) fetchCountriesPage(offset int, limit int) ([]restCountry, restCountriesMeta, error) {
	url := fmt.Sprintf(
		"%s?limit=%d&offset=%d&response_fields=names.common,codes.alpha_2,flag.url_svg",
		s.Config.GrandiGiochiniGiorno.RestCountriesURL, limit, offset,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, restCountriesMeta{}, fmt.Errorf("errore creazione richiesta: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.Config.GrandiGiochiniGiorno.RestCountriesAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, restCountriesMeta{}, fmt.Errorf("errore download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, restCountriesMeta{}, fmt.Errorf("richiesta a REST Countries fallita con status %d", resp.StatusCode)
	}

	var risposta restCountriesResponse
	err = json.NewDecoder(resp.Body).Decode(&risposta)
	if err != nil {
		return nil, restCountriesMeta{}, fmt.Errorf("errore decodifica risposta: %w", err)
	}

	return risposta.Data.Objects, risposta.Data.Meta, nil
}

func (s GrandiGiochiniGiornoService) getRandomWord(letters int) (string, error) {
	resp, err := http.Get(s.Config.GrandiGiochiniGiorno.Dictionary)
	if err != nil {
		return "", fmt.Errorf("errore download: %w", err)
	}

	var parole []string

	err = func() error {
		defer func() {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			parola := strings.TrimSpace(scanner.Text())
			if len([]rune(parola)) == letters {
				parole = append(parole, strings.ToUpper(parola))
			}
		}
		return scanner.Err()
	}()

	if err != nil {
		return "", err
	}

	if len(parole) == 0 {
		return "", fmt.Errorf("nessuna parola trovata")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	parolaEstratta := parole[r.Intn(len(parole))]

	parole = nil
	runtime.GC()

	return parolaEstratta, nil
}
