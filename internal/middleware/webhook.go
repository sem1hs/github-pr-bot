package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/semih/github-pr-bot/internal/service"
)

func VerifyAndFilter(webhookService *service.WebhookService) gin.HandlerFunc{
	return func(c *gin.Context){
		// 1. Gövde
		body, err := io.ReadAll(c.Request.Body)
		if err != nil{
			c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"error":"govde okunamadi"})
			return
		}

		// 2. İmza
		signature := c.GetHeader("X-Hub-Signature-256")
		if !webhookService.VerifySignature(body,signature){
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"gecersiz imza"})
			return
		}

		// 3. Event filtre
		if c.GetHeader("X-GitHub-Event") != "pull_request"{
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status":"yok sayildi"})
			return
		}

		// 4. Gövdeyi geri ver
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Next()
	}
}