package services

import "net/http"

type Service interface {
	HandleWebhook(w http.ResponseWriter, r *http.Request)
	ValidateWebhookSignature(signatureHeader string, rawBody []byte, secret string) (bool, error)
}
