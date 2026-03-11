package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"backend/internal/types"
)

func (s *HotpepperService) SearchGourmet(params types.GourmetSearchParams) (*types.GourmetSearchResult, error) {

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

	var raw gourmetSearchResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(raw.Results.Error) > 0 {
		return nil, convertAPIError(&raw.Results.Error[0])
	}

	return &types.GourmetSearchResult{
		Results: types.GourmetSearchResultData{
			ResultsAvailable: raw.Results.ResultsAvailable,
			ResultsStart:     raw.Results.ResultsStart,
			Shop:             raw.Results.Shop,
		},
	}, nil
}

func (s *HotpepperService) GetGourmetDetail(id string) (*types.GourmetSearchResult, error) {
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

	var raw gourmetSearchResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(raw.Results.Error) > 0 {
		return nil, convertAPIError(&raw.Results.Error[0])
	}

	return &types.GourmetSearchResult{
		Results: types.GourmetSearchResultData{
			ResultsAvailable: raw.Results.ResultsAvailable,
			ResultsStart:     raw.Results.ResultsStart,
			Shop:             raw.Results.Shop,
		},
	}, nil
}

// --- HotPepper グルメ API レスポンスの JSON デコード専用プライベート型 ---
// これらは service 層の実装詳細であり、外部パッケージに公開しない。

type gourmetSearchResponse struct {
	Results gourmetSearchResults `json:"results"`
}

type gourmetSearchResults struct {
	ResultsAvailable int                      `json:"results_available"`
	ResultsStart     int                      `json:"results_start"`
	Shop             []types.Shop             `json:"shop"`
	Error            []hotPepperResponseError `json:"error"`
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
