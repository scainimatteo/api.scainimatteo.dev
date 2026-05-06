package services

type Config struct {
	DB                   DB      `json:"db"`
	FireflyPushoverToken string  `json:"firefly_pushover_token"`
	VikunjaPushoverToken string  `json:"vikunja_pushover_token"`
	Vikunja              Vikunja `json:"vikunja"`
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
