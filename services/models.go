package services

type Config struct {
	DB                   DB      `json:"db"`
	VikunjaPushoverToken string  `json:"vikunja_pushover_token"`
	Vikunja              Vikunja `json:"vikunja"`
	Firefly              Firefly `json:"firefly"`
	PushoverUser         string  `json:"pushover_user"`
	Port                 string  `json:"port"`
}

type DB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Vikunja struct {
	BaseURL    string `json:"base_url"`
	APIToken   string `json:"api_token"`
	CalendarID string `json:"calendar_id"`
}

type Firefly struct {
	BaseURL       string         `json:"base_url"`
	PushoverToken string         `json:"pushover_token"`
	APIKey        string         `json:"api_key"`
	Sources       FireflySources `json:"sources"`
}

type FireflySources struct {
	Bper   string `json:"bper"`
	Ticket string `json:"ticket"`
}
