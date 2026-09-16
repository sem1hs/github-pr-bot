package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/semih/github-pr-bot/internal/config"
	"github.com/semih/github-pr-bot/internal/controller"
	"github.com/semih/github-pr-bot/internal/middleware"
	"github.com/semih/github-pr-bot/internal/service"
)

func Setup(cfg config.Config) *gin.Engine{
	router := gin.Default()

	// Sağlık kontrolü endpoint
	router.GET("/health",func (c *gin.Context)  {
		c.JSON(http.StatusOK, gin.H{"status":"ok"})
	})

	// Bağımlılıklar
	webhookService := service.NewWebhookService(cfg.GitHubWebhookSecret)
	githubService := service.NewGithubService(cfg.GitHubToken)
	webhook := controller.NewWebhookController(githubService)

	// Route zinciri
	router.POST("/webhook", middleware.VerifyAndFilter(webhookService), webhook.Handle)

	return router
}