package controller

import (
	"log"
	"net/http"

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

	// İş mantığı servicede
	if err:= wc.githubService.ProcessPullRequest(payload); err != nil{
		log.Printf("PR isleme hatasi: %v",err)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"isleme basarisiz"})
		return
	}

	// Başarılı cevap
	c.JSON(http.StatusOK, gin.H{"status":"alindi","pr":payload.Number})
}