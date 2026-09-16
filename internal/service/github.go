package service

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/semih/github-pr-bot/internal/model"
)

type GithubService struct {
	token string
	httpClient *http.Client
}

func NewGithubService(token string) *GithubService {
	return &GithubService{
		token: token,
		httpClient: &http.Client{Timeout: 15* time.Second},
	}
}

func (s *GithubService) ProcessPullRequest(payload model.PullRequestEvent) error {
	return nil
}

// fetchDiff = PR'ın ham diff metnini GitHub'dan çeker
func(s *GithubService) fetchDiff(owner, repo string, number int) (string,error){
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, number)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil{
		return "", err
	}

	// Diff formatı için Header
	req.Header.Set("Authorization","Bearer "+s.token)
	req.Header.Set("Accept","application/vnd.github.v3.diff")
	req.Header.Set("X-GitHub-Api-Version","2022-11-28")

	resp, err := s.httpClient.Do(req)
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
