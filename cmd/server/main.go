package main

import (
	"log"

	"github.com/semih/github-pr-bot/internal/config"
	"github.com/semih/github-pr-bot/internal/router"
)

func main() {

	// Ayarları yükle
	cfg := config.Load()

	// Gin router oluştur
	router := router.Setup(cfg)

	// Sunucuyu başlat
	addr := ":" + cfg.Port
	log.Printf("sunucu %s portunda calisiyor", addr)
	if err := router.Run(addr); err!=nil{
		log.Fatalf("sunucu baslastilamadi: %v",err)
	}
}