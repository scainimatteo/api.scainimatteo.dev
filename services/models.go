package services

type Config struct {
	FireflyPushoverToken string `json:"firefly_pushover_token"`
	VikunjaPushoverToken string `json:"vikunja_pushover_token"`
	PushoverUser         string `json:"pushover_user"`
	Port                 string `json:"port"`
}
