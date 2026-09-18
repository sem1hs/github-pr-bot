package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config = uygulamanın tüm ayarları tek yerde.
type Config struct {
	Port                string
	GitHubWebhookSecret string
	GitHubAppID             int64
	GitHubAppPrivateKeyPath string
	GCPProjectID        string
	GCPLocation         string
	VertexModel         string
}

// Load = env değişkenlerini okur, Config döner.
func Load() Config{

	if err:= godotenv.Load(); err !=nil{
		log.Println(".env dosyasi bulunamadi, sistem env kullanilacak")
	}
	return Config{
		Port:                    getEnv("PORT", "8080"),
		GitHubWebhookSecret:     mustEnv("GITHUB_WEBHOOK_SECRET"),
		GitHubAppID:             mustEnvInt64("GITHUB_APP_ID"),
		GitHubAppPrivateKeyPath: mustEnv("GITHUB_APP_PRIVATE_KEY_PATH"),
		GCPProjectID:            mustEnv("GCP_PROJECT_ID"),
		GCPLocation:             getEnv("GCP_LOCATION", "us-central1"),
		VertexModel:             getEnv("VERTEX_MODEL", "gemini-2.5-flash"),
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

// mustEnvInt64 = zorunlu env'i int64'e çevirir
func mustEnvInt64(key string) int64 {
	v := mustEnv(key)
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil{
		log.Fatalf("env %s sayi olmali: %v", key, err)
	}
	return n
}