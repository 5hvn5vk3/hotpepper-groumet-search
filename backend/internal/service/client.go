package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const hotpepperBaseURL = "https://webservice.recruit.co.jp/hotpepper"

type HotpepperService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewHotpepperService(apiKey string) *HotpepperService {

	return &HotpepperService{
		apiKey:  apiKey,
		baseURL: hotpepperBaseURL,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
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

// convertAPIError はホットペッパーAPIのエラーをHotpepperAPIErrorに変換する。
// 生のエラーメッセージは情報漏洩防止のためログのみに出力し、外部には返さない。
func convertAPIError(apiErr *hotPepperResponseError) error {
	log.Printf("hotpepper API error: code=%d, message=%s", apiErr.Code, apiErr.Message)
	return &HotpepperAPIError{Code: apiErr.Code}
}

// hotPepperResponseError は HotPepper API レスポンスの "error" フィールドをデコードするための構造体。
// Go の error インターフェースを実装した HotpepperAPIError とは役割が異なる。
type hotPepperResponseError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
