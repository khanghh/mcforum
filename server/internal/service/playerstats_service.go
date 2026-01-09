package service

import (
	"bbs-go/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var PlayerStatsService = newPlayerStatsService()

type playerStatsService struct {
	httpClient http.Client
}

func newPlayerStatsService() *playerStatsService {
	return &playerStatsService{
		httpClient: http.Client{Timeout: 5 * time.Second},
	}
}

type PlayTime struct {
	ServerName string `json:"serverName"`
	TotalTime  int64  `json:"totalTime"`
}

type PlayerStats struct {
	PlayTimes []PlayTime `json:"playTimes"`
}

func (s *playerStatsService) GetStats(username string) (*PlayerStats, error) {
	reqURL := fmt.Sprintf("%s/players/%s/stats", config.Instance().PlayerStatsAPIURL, username)

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp struct {
		Data  []PlayTime `json:"data"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	data, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &apiResp); err != nil {
		log.Printf("Error decoding playerstats response: err=%v, status=%d, body=%s", err, resp.StatusCode, string(data))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return &PlayerStats{PlayTimes: apiResp.Data}, nil
	}

	if apiResp.Error != nil {
		return nil, fmt.Errorf("playerstats api error: %s", apiResp.Error.Message)
	}
	return nil, fmt.Errorf("playerstats api return unexpected error: %d", resp.StatusCode)
}

func (s *playerStatsService) GetTotalPlayTimeSec(username string) (int64, error) {
	stats, err := s.GetStats(username)
	if err != nil {
		return 0, err
	}
	var total int64
	for _, pt := range stats.PlayTimes {
		total += pt.TotalTime
	}
	return total, nil
}
