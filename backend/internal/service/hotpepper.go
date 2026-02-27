package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"backend/internal/types"
)

const hotpepperBaseURL = "http://webservice.recruit.co.jp/hotpepper"

type HotpepperService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewHotpepperService(apiKey string) *HotpepperService {

	return &HotpepperService{
		apiKey:  apiKey,
		baseURL: hotpepperBaseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *HotpepperService) SearchGourmet(params types.GourmetSearchParams) (*types.GourmetSearchResponse, error) {

	queryParams := url.Values{}

	queryParams.Set("key", s.apiKey)
	queryParams.Set("type", "lite+credit_card")
	queryParams.Set("format", "json")

	if params.Lat != 0 && params.Lng != 0 {
		queryParams.Set("lat", strconv.FormatFloat(params.Lat, 'f', -1, 64))
		queryParams.Set("lng", strconv.FormatFloat(params.Lng, 'f', -1, 64))
		if params.Range > 0 {

			queryParams.Set("range", strconv.Itoa(params.Range))
		}
	}

	if params.Address != "" {
		queryParams.Set("address", params.Address)
	}
	if params.Genre != "" {
		queryParams.Set("genre", params.Genre)
	}
	if params.Keyword != "" {
		queryParams.Set("keyword", params.Keyword)
	}

	start := clampInt(params.Start, 1, 1000)
	count := clampInt(params.Count, 1, 100)

	queryParams.Set("start", strconv.Itoa(start))
	queryParams.Set("count", strconv.Itoa(count))

	apiURL := fmt.Sprintf("%s/gourmet/v1/?%s", s.baseURL, queryParams.Encode())

	body, err := s.fetchAPI(apiURL)
	if err != nil {
		return nil, err
	}

	var response types.GourmetSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (s *HotpepperService) GetGourmetDetail(id string) (*types.GourmetSearchResponse, error) {
	shopID := strings.TrimSpace(id)
	if shopID == "" {
		return nil, fmt.Errorf("id is required")
	}

	queryParams := url.Values{}
	queryParams.Set("key", s.apiKey)
	queryParams.Set("format", "json")
	queryParams.Set("id", shopID)

	apiURL := fmt.Sprintf("%s/gourmet/v1/?%s", s.baseURL, queryParams.Encode())

	body, err := s.fetchAPI(apiURL)
	if err != nil {
		return nil, err
	}

	var response types.GourmetSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (s *HotpepperService) GetGenreMaster() (*types.GenreMasterResponse, error) {

	queryParams := url.Values{}
	queryParams.Set("key", s.apiKey)
	queryParams.Set("format", "json")

	apiURL := fmt.Sprintf("%s/genre/v1/?%s", s.baseURL, queryParams.Encode())

	body, err := s.fetchAPI(apiURL)
	if err != nil {
		return nil, err
	}

	var response types.GenreMasterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (s *HotpepperService) fetchAPI(apiURL string) ([]byte, error) {

	resp, err := s.client.Get(apiURL)

	if err != nil {

		return nil, fmt.Errorf("failed to fetch API: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func clampInt(value, min, max int) int {

	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}
