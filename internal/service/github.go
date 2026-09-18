package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/semih/github-pr-bot/internal/model"
)

type GithubService struct {
	appTransfort *ghinstallation.AppsTransport
	llm *LLMService
}

func NewGithubService(appTransfort *ghinstallation.AppsTransport, llm *LLMService) *GithubService {
	return &GithubService{
		appTransfort: appTransfort,
		llm: llm,
	}
}

const maxDiffChars = 1200

func (s *GithubService) clientForInstallation(installationID int64) *http.Client {
	itr := ghinstallation.NewFromAppsTransport(s.appTransfort, installationID)
	return &http.Client{Transport: itr, Timeout: 15 * time.Second}
}

func (s *GithubService) ProcessPullRequest(ctx context.Context, payload model.PullRequestEvent) error {
	if payload.Action != "opened" && payload.Action != "reopened"{
		log.Printf("action atlandi: %s (pr=#%d)", payload.Action, payload.Number)
		return nil
	}
	
	client := s.clientForInstallation(payload.Installation.ID)

	diff, err := s.fetchDiff(
		ctx,
		client,
		payload.Repository.Owner.Login,
		payload.Repository.Name,
		payload.Number,
	)

	if err != nil {
		return err
	}

	log.Printf("diff alindi: repo=%s pr=#%d boyut=%d byte", payload.Repository.FullName, payload.Number, len(diff))
	
	if strings.TrimSpace(diff) == ""{
		log.Printf("bos diff, atlaniyor: pr=#%d", payload.Number)
		return nil
	}

	if len(diff) > maxDiffChars{
		log.Printf("diff kirpildi: %d -> %d byte (pr=#%d)", len(diff), maxDiffChars, payload.Number)
		diff = diff[:maxDiffChars] + "\n\n[... diff cok buyuk, kirpildi ...]"
	}
	
	description, err := s.llm.GenerateDescription(ctx, diff)
	if err != nil{
		return err
	}
	log.Printf("LLM aciklama uretti (%d karakter)", len(description))

	// 3. PR'a yorum yaz.
	if err := s.updatePullRequestBody(ctx, client, payload.Repository.Owner.Login, payload.Repository.Name, payload.Number, description); err != nil {
		return err
	}
	log.Printf("yorum yazildi: pr=#%d", payload.Number)

	return nil
}

const botMarker = "\n\n<!-- generated-by-pr-bot -->"

func (s *GithubService) updatePullRequestBody(ctx context.Context, client *http.Client, owner, repo string, number int, body string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d",owner,repo,number)

	payload := map[string]string{"body": body + botMarker}
	jsonBody, err := json.Marshal(payload)
	if err != nil{
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(jsonBody))
	if err != nil{
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil{
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("PR body guncelleme basarisiz: %d %s", resp.StatusCode, string(b))
	}
	return nil
}

// fetchDiff = PR'ın ham diff metnini GitHub'dan çeker
func(s *GithubService) fetchDiff(ctx context.Context, client *http.Client, owner, repo string, number int) (string,error){
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, number)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil{
		return "", err
	}

	// Diff formatı için Header
	req.Header.Set("Accept", "application/vnd.github.v3.diff")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil{
		return "", err
	}

	// fonksiyon biterken çalışır, body kapatma garanti
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("github diff cekme basarisiz: %d %s", resp.StatusCode, string(b))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		return "", err
	}

	return string(body), nil
}