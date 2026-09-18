package router

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/gin-gonic/gin"
	"github.com/semih/github-pr-bot/internal/config"
	"github.com/semih/github-pr-bot/internal/controller"
	"github.com/semih/github-pr-bot/internal/middleware"
	"github.com/semih/github-pr-bot/internal/service"
)

func Setup(cfg config.Config) (*gin.Engine, error){
	router := gin.Default()

	// Sağlık kontrolü endpoint
	router.GET("/health",func (c *gin.Context)  {
		c.JSON(http.StatusOK, gin.H{"status":"ok"})
	})

	// LLM client
	llmService, err := service.NewLLMService(
		context.Background(),
		cfg.GCPProjectID,
		cfg.GCPLocation,
		cfg.VertexModel,
	)

	if err != nil{
		return nil, err
	}

	appTransport, err := ghinstallation.NewAppsTransportKeyFromFile(
		http.DefaultTransport,
		cfg.GitHubAppID,
		cfg.GitHubAppPrivateKeyPath,
	)

	if err != nil{
		return nil, err
	}

	// Bağımlılıklar
	webhookService := service.NewWebhookService(cfg.GitHubWebhookSecret)
	githubService := service.NewGithubService(appTransport, llmService)
	webhook := controller.NewWebhookController(githubService)

	// Route zinciri
	router.POST("/webhook", middleware.VerifyAndFilter(webhookService), webhook.Handle)

	return router, nil
}