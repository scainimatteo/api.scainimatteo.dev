package grandigiochinigiorno

type restCountry struct {
	Names restCountryNames `json:"names"`
	Codes restCountryCodes `json:"codes"`
	Flag  restCountryFlag  `json:"flag"`
}

type restCountryNames struct {
	Common string `json:"common"`
}

type restCountryCodes struct {
	Alpha2 string `json:"alpha_2"`
}

type restCountryFlag struct {
	URLSvg string `json:"url_svg"`
	URLPng string `json:"url_png"`
	Emoji  string `json:"emoji"`
}

type restCountriesResponse struct {
	Data restCountriesData `json:"data"`
}

type restCountriesData struct {
	Objects []restCountry     `json:"objects"`
	Meta    restCountriesMeta `json:"meta"`
}

type restCountriesMeta struct {
	Total  int  `json:"total"`
	Count  int  `json:"count"`
	Limit  int  `json:"limit"`
	Offset int  `json:"offset"`
	More   bool `json:"more"`
}
