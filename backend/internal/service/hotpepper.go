// このファイルが属するパッケージ名を宣言する。
package service

// ここから使用する外部パッケージの列挙を始める。
import (
	// この行は処理の一部として必要な操作を実行する。
	"encoding/json"
	// この行は処理の一部として必要な操作を実行する。
	"fmt"
	// この行は処理の一部として必要な操作を実行する。
	"io"
	// この行は処理の一部として必要な操作を実行する。
	"log"
	// この行は処理の一部として必要な操作を実行する。
	"net/http"
	// この行は処理の一部として必要な操作を実行する。
	"net/url"
	// この行は処理の一部として必要な操作を実行する。
	"strconv"
	// この行は処理の一部として必要な操作を実行する。
	"strings"
	// この行は処理の一部として必要な操作を実行する。
	"time"

	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/types"
	// このブロックや引数リストを閉じる。
)

// この行で定数を定義する。
const hotpepperBaseURL = "https://webservice.recruit.co.jp/hotpepper"

// この行で構造体の型定義を始める。
type HotpepperService struct {
	// この行は処理の一部として必要な操作を実行する。
	apiKey string
	// この行は処理の一部として必要な操作を実行する。
	baseURL string
	// この行は処理の一部として必要な操作を実行する。
	client *http.Client
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func NewHotpepperService(apiKey string) *HotpepperService {

	// この時点で関数の結果を返して処理を終了する。
	return &HotpepperService{
		// この行は処理の一部として必要な操作を実行する。
		apiKey: apiKey,
		// この行は処理の一部として必要な操作を実行する。
		baseURL: hotpepperBaseURL,
		// 構造体やマップのフィールドへ値を設定する。
		client: &http.Client{Timeout: 60 * time.Second},
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (s *HotpepperService) SearchGourmet(params types.GourmetSearchParams) (*types.GourmetSearchResponse, error) {

	// この行で新しい変数を宣言しつつ値を代入する。
	queryParams := url.Values{}

	// この行は処理の一部として必要な操作を実行する。
	queryParams.Set("key", s.apiKey)
	// この行は処理の一部として必要な操作を実行する。
	queryParams.Set("type", "lite+credit_card")
	// この行は処理の一部として必要な操作を実行する。
	queryParams.Set("format", "json")

	// 条件を評価し、真のときだけ次の処理を実行する。
	if params.Lat != 0 && params.Lng != 0 {
		// この行は処理の一部として必要な操作を実行する。
		queryParams.Set("lat", strconv.FormatFloat(params.Lat, 'f', -1, 64))
		// この行は処理の一部として必要な操作を実行する。
		queryParams.Set("lng", strconv.FormatFloat(params.Lng, 'f', -1, 64))
		// 条件を評価し、真のときだけ次の処理を実行する。
		if params.Range > 0 {
			// この行は処理の一部として必要な操作を実行する。
			queryParams.Set("range", strconv.Itoa(params.Range))
			// ここで処理ブロックを終了する。
		}
		// ここで処理ブロックを終了する。
	}

	// 条件を評価し、真のときだけ次の処理を実行する。
	if params.Address != "" {
		// この行は処理の一部として必要な操作を実行する。
		queryParams.Set("address", params.Address)
		// ここで処理ブロックを終了する。
	}
	// 条件を評価し、真のときだけ次の処理を実行する。
	if params.Genre != "" {
		// この行は処理の一部として必要な操作を実行する。
		queryParams.Set("genre", params.Genre)
		// ここで処理ブロックを終了する。
	}
	// 条件を評価し、真のときだけ次の処理を実行する。
	if params.Keyword != "" {
		// この行は処理の一部として必要な操作を実行する。
		queryParams.Set("keyword", params.Keyword)
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	start := clampInt(params.Start, 1, 1000)
	// この行で新しい変数を宣言しつつ値を代入する。
	count := clampInt(params.Count, 1, 100)

	// この行は処理の一部として必要な操作を実行する。
	queryParams.Set("start", strconv.Itoa(start))
	// この行は処理の一部として必要な操作を実行する。
	queryParams.Set("count", strconv.Itoa(count))

	// この行で新しい変数を宣言しつつ値を代入する。
	apiURL := fmt.Sprintf("%s/gourmet/v1/?%s", s.baseURL, queryParams.Encode())

	// この行で新しい変数を宣言しつつ値を代入する。
	body, err := s.fetchAPI(apiURL)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この時点で関数の結果を返して処理を終了する。
		return nil, err
		// ここで処理ブロックを終了する。
	}

	// この行で変数を定義する。
	var response types.GourmetSearchResponse
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err := json.Unmarshal(body, &response); err != nil {
		// この時点で関数の結果を返して処理を終了する。
		return nil, fmt.Errorf("failed to decode response: %w", err)
		// ここで処理ブロックを終了する。
	}

	// 条件を評価し、真のときだけ次の処理を実行する。
	if len(response.Results.Error) > 0 {
		// この時点で関数の結果を返して処理を終了する。
		return nil, convertAPIError(&response.Results.Error[0])
		// ここで処理ブロックを終了する。
	}

	// この時点で関数の結果を返して処理を終了する。
	return &response, nil
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (s *HotpepperService) GetGourmetDetail(id string) (*types.GourmetSearchResponse, error) {
	// この行で新しい変数を宣言しつつ値を代入する。
	shopID := strings.TrimSpace(id)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if shopID == "" {
		// この時点で関数の結果を返して処理を終了する。
		return nil, fmt.Errorf("id is required")
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	queryParams := url.Values{}
	// この行は処理の一部として必要な操作を実行する。
	queryParams.Set("key", s.apiKey)
	// この行は処理の一部として必要な操作を実行する。
	queryParams.Set("format", "json")
	// この行は処理の一部として必要な操作を実行する。
	queryParams.Set("id", shopID)

	// この行で新しい変数を宣言しつつ値を代入する。
	apiURL := fmt.Sprintf("%s/gourmet/v1/?%s", s.baseURL, queryParams.Encode())

	// この行で新しい変数を宣言しつつ値を代入する。
	body, err := s.fetchAPI(apiURL)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この時点で関数の結果を返して処理を終了する。
		return nil, err
		// ここで処理ブロックを終了する。
	}

	// この行で変数を定義する。
	var response types.GourmetSearchResponse
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err := json.Unmarshal(body, &response); err != nil {
		// この時点で関数の結果を返して処理を終了する。
		return nil, fmt.Errorf("failed to decode response: %w", err)
		// ここで処理ブロックを終了する。
	}

	// 条件を評価し、真のときだけ次の処理を実行する。
	if len(response.Results.Error) > 0 {
		// この時点で関数の結果を返して処理を終了する。
		return nil, convertAPIError(&response.Results.Error[0])
		// ここで処理ブロックを終了する。
	}

	// この時点で関数の結果を返して処理を終了する。
	return &response, nil
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (s *HotpepperService) GetGenreMaster() (*types.GenreMasterResponse, error) {

	// この行で新しい変数を宣言しつつ値を代入する。
	queryParams := url.Values{}
	// この行は処理の一部として必要な操作を実行する。
	queryParams.Set("key", s.apiKey)
	// この行は処理の一部として必要な操作を実行する。
	queryParams.Set("format", "json")

	// この行で新しい変数を宣言しつつ値を代入する。
	apiURL := fmt.Sprintf("%s/genre/v1/?%s", s.baseURL, queryParams.Encode())

	// この行で新しい変数を宣言しつつ値を代入する。
	body, err := s.fetchAPI(apiURL)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この時点で関数の結果を返して処理を終了する。
		return nil, err
		// ここで処理ブロックを終了する。
	}

	// この行で変数を定義する。
	var response types.GenreMasterResponse
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err := json.Unmarshal(body, &response); err != nil {
		// この時点で関数の結果を返して処理を終了する。
		return nil, fmt.Errorf("failed to decode response: %w", err)
		// ここで処理ブロックを終了する。
	}

	// 条件を評価し、真のときだけ次の処理を実行する。
	if len(response.Results.Error) > 0 {
		// この時点で関数の結果を返して処理を終了する。
		return nil, convertAPIError(&response.Results.Error[0])
		// ここで処理ブロックを終了する。
	}

	// この時点で関数の結果を返して処理を終了する。
	return &response, nil
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (s *HotpepperService) fetchAPI(apiURL string) ([]byte, error) {

	// この行で新しい変数を宣言しつつ値を代入する。
	resp, err := s.client.Get(apiURL)

	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {

		// この時点で関数の結果を返して処理を終了する。
		return nil, fmt.Errorf("failed to fetch API: %w", err)
		// ここで処理ブロックを終了する。
	}

	// 関数終了時に実行したい後処理を予約する。
	defer resp.Body.Close()

	// この行で新しい変数を宣言しつつ値を代入する。
	body, err := io.ReadAll(resp.Body)

	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この時点で関数の結果を返して処理を終了する。
		return nil, fmt.Errorf("failed to read response: %w", err)
		// ここで処理ブロックを終了する。
	}

	// 条件を評価し、真のときだけ次の処理を実行する。
	if resp.StatusCode != http.StatusOK {

		// この時点で関数の結果を返して処理を終了する。
		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(body))
		// ここで処理ブロックを終了する。
	}

	// この時点で関数の結果を返して処理を終了する。
	return body, nil
	// ここで処理ブロックを終了する。
}

// この行は補足のためのコメントを記述する。
// convertAPIError はホットペッパーAPIのエラーをHotpepperAPIErrorに変換する。
// この行は補足のためのコメントを記述する。
// 生のエラーメッセージは情報漏洩防止のためログのみに出力し、外部には返さない。
// この行で関数定義を始める。
func convertAPIError(apiErr *types.APIError) error {
	// デバッグや障害調査のためにログを出力する。
	log.Printf("hotpepper API error: code=%d, message=%s", apiErr.Code, apiErr.Message)
	// この時点で関数の結果を返して処理を終了する。
	return &types.HotpepperAPIError{Code: apiErr.Code}
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func clampInt(value, min, max int) int {

	// 条件を評価し、真のときだけ次の処理を実行する。
	if value < min {
		// この時点で関数の結果を返して処理を終了する。
		return min
		// ここで処理ブロックを終了する。
	}

	// 条件を評価し、真のときだけ次の処理を実行する。
	if value > max {
		// この時点で関数の結果を返して処理を終了する。
		return max
		// ここで処理ブロックを終了する。
	}

	// この時点で関数の結果を返して処理を終了する。
	return value
	// ここで処理ブロックを終了する。
}
