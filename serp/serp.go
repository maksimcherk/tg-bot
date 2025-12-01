package serp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SerpClient struct {
	apiKey string
}

func New(apiKey string) *SerpClient {
	return &SerpClient{
		apiKey: apiKey,
	}
}

type SerpResponse struct {
	OrganicResults []struct {
		Title       string `json:"title"`
		Snippet     string `json:"snippet"`
		DisplayLink string `json:"displayed_link"`
		Link        string `json:"link"`
	} `json:"organic_results"`
}

func (s SerpClient) CheckURL(target string) (bool, string, error) {
	escaped := url.QueryEscape(target)

	apiURL := fmt.Sprintf(
		"https://serpapi.com/search?engine=google&q=%s&api_key=%s",
		escaped,
		s.apiKey,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return false, "", fmt.Errorf("ошибка подключения к SerpAPI")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, "", fmt.Errorf("SerpAPI вернул код: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("ошибка чтения ответа")
	}

	var data SerpResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, "", fmt.Errorf("ошибка JSON: %w", err)
	}

	redFlags := []string{
		"phishing",
		"scam",
		"fraud",
		"malware",
		"unsafe",
		"virus",
		"fake",
		"warning",
	}

	for _, item := range data.OrganicResults {
		txt := strings.ToLower(item.Title + " " + item.Snippet)

		for _, flag := range redFlags {
			if strings.Contains(txt, flag) {
				return true, flag, nil
			}
		}
	}

	return false, "", nil
}
