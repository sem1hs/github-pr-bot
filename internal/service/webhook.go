package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type WebhookService struct {
	secret string
}

func NewWebhookService(secret string) *WebhookService {
	return &WebhookService{secret: secret}
}

func (s *WebhookService) VerifySignature(body []byte, signature string) bool {
	if signature == "" {
		return false
	}

	mac := hmac.New(sha256.New,[]byte(s.secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}