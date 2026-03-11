package service

import (
	"encoding/json"
	"fmt"
	"net/url"

	"backend/internal/types"
)

func (s *HotpepperService) GetGenreMaster() (*types.GenreMasterResult, error) {

	queryParams := url.Values{}
	queryParams.Set("key", s.apiKey)
	queryParams.Set("format", "json")

	apiURL := fmt.Sprintf("%s/genre/v1/?%s", s.baseURL, queryParams.Encode())

	body, err := s.fetchAPI(apiURL)
	if err != nil {
		return nil, err
	}

	var raw genreMasterResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(raw.Results.Error) > 0 {
		return nil, convertAPIError(&raw.Results.Error[0])
	}

	return &types.GenreMasterResult{
		Results: types.GenreMasterResultData{
			Genre: raw.Results.Genre,
		},
	}, nil
}

// --- HotPepper ジャンルマスタ API レスポンスの JSON デコード専用プライベート型 ---
// これらは service 層の実装詳細であり、外部パッケージに公開しない。

type genreMasterResponse struct {
	Results genreMasterResults `json:"results"`
}

type genreMasterResults struct {
	Genre []types.HotpepperGenre   `json:"genre"`
	Error []hotPepperResponseError `json:"error"`
}
