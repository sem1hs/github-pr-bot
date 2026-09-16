package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config = uygulamanın tüm ayarları tek yerde.
type Config struct {
	Port                string
	GitHubWebhookSecret string
	GitHubToken 		string
}

// Load = env değişkenlerini okur, Config döner.
func Load() Config{

	if err:= godotenv.Load(); err !=nil{
		log.Println(".env dosyasi bulunamadi, sistem env kullanilacak")
	}
	return Config{
		Port:                getEnv("PORT", "8080"),
		GitHubWebhookSecret: mustEnv("GITHUB_WEBHOOK_SECRET"),
		GitHubToken: mustEnv("GITHUB_TOKEN"),
	}
}

// getEnv = env varsa onu, yoksa varsayılanı döner.
func getEnv(key, fallback string) string{
	
	if v:= os.Getenv(key); v!= "" {
		return v
	}
	return fallback
}

// mustEnv = env yoksa uygulamayı çökertir. Kritik ayarlar için.
func mustEnv(key string) string{
	v:= os.Getenv(key)
	if v == ""{
		log.Fatal("zorunlu env eksik %s",key)
	}
	return v
}