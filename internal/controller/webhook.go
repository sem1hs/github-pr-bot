package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/semih/github-pr-bot/internal/model"
	"github.com/semih/github-pr-bot/internal/service"
)

// WebhookController = webhook isteklerini işler.
// secret alanı imza doğrulama için tutulur.
type WebhookController struct {
	githubService *service.GithubService
}

func NewWebhookController(gh *service.GithubService) *WebhookController{
	return &WebhookController{githubService: gh}
}

// Handle = POST /webhook endpoint'i.
func (wc *WebhookController) Handle(c *gin.Context){
	// Parse, Gin JSON'u structa bağlar
	var payload model.PullRequestEvent
	if err := c.ShouldBindJSON(&payload); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "json parse edilemedi"})
		return
	}

	// Arka planda işle
	go func(){
		ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
		defer cancel()

		// İş mantığı servicede
		if err:= wc.githubService.ProcessPullRequest(ctx, payload); err != nil{
		log.Printf("PR isleme hatasi: pr=#%d %v",payload.Number, err)
		return
	}
	}()

	// Başarılı cevap
	c.JSON(http.StatusOK, gin.H{"status":"alindi","pr":payload.Number})
}