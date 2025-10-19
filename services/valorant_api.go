package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ValorantAPIClient клиент для работы с Valorant API
type ValorantAPIClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewValorantAPIClient создает новый клиент
func NewValorantAPIClient(apiKey string) *ValorantAPIClient {
	return &ValorantAPIClient{
		BaseURL: "https://pd.riotgames.com", // или другой базовый URL
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PlayerInfo информация об игроке
type PlayerInfo struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// MatchInfo информация о матче
type MatchInfo struct {
	MatchID string `json:"matchId"`
	Map     string `json:"map"`
	Mode    string `json:"mode"`
	Result  string `json:"result"`
	Score   string `json:"score"`
	Date    string `json:"date"`
}

// GetPlayerByName получает информацию об игроке по имени
func (c *ValorantAPIClient) GetPlayerByName(gameName, tag string) (*PlayerInfo, error) {
	// TODO: Реализовать запрос к API
	// Пример URL: /riot/account/v1/accounts/by-riot-id/{gameName}/{tag}

	url := fmt.Sprintf("%s/riot/account/v1/accounts/by-riot-id/%s/%s", c.BaseURL, gameName, tag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var player PlayerInfo
	err = json.Unmarshal(body, &player)
	if err != nil {
		return nil, err
	}

	return &player, nil
}

// GetPlayerMatches получает матчи игрока
func (c *ValorantAPIClient) GetPlayerMatches(puuid string, count int) ([]string, error) {
	// TODO: Реализовать запрос к API
	// Пример URL: /match/v1/matchlists/by-puuid/{puuid}

	url := fmt.Sprintf("%s/match/v1/matchlists/by-puuid/%s?count=%d", c.BaseURL, puuid, count)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Matches []struct {
			MatchID string `json:"matchId"`
		} `json:"matches"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var matchIDs []string
	for _, match := range response.Matches {
		matchIDs = append(matchIDs, match.MatchID)
	}

	return matchIDs, nil
}

// GetMatchDetails получает детали матча
func (c *ValorantAPIClient) GetMatchDetails(matchID string) (*MatchInfo, error) {
	// TODO: Реализовать запрос к API
	// Пример URL: /match/v1/matches/{matchId}

	url := fmt.Sprintf("%s/match/v1/matches/%s", c.BaseURL, matchID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var match MatchInfo
	err = json.Unmarshal(body, &match)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

// GetPlayerRank получает ранг игрока
func (c *ValorantAPIClient) GetPlayerRank(puuid string) (string, int, error) {
	// TODO: Реализовать запрос к API
	// Пример URL: /mmr/v1/players/{puuid}

	url := fmt.Sprintf("%s/mmr/v1/players/%s", c.BaseURL, puuid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("X-Riot-Token", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	var response struct {
		CurrentTier string `json:"currentTier"`
		RankRating  int    `json:"rankRating"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", 0, err
	}

	return response.CurrentTier, response.RankRating, nil
}
