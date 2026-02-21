package service

import (
	"encoding/json" // JSON のエンコード・デコードを提供
	"fmt"           // フォーマット済み入出力を提供（文字列の整形など）
	"io"            // 入出力の基本インターフェースを提供
	"net/http"      // HTTPクライアント・サーバー機能を提供
	"net/url"       // URLの解析と構築を提供
	"strconv"       // 文字列と数値の変換を提供

	"backend/internal/types" // リクエスト・レスポンスの型定義をインポート
)

// hotpepperBaseURL はホットペッパーAPIのベースURL
// constは定数を宣言するキーワード（実行中に値を変更できない）
const hotpepperBaseURL = "http://webservice.recruit.co.jp/hotpepper"

// HotpepperService は ホットペッパーAPI とのやり取りを管理するサービス層の構造体
// サービス層は、ビジネスロジックとAPI通信を担当する層
type HotpepperService struct {
	apiKey  string // ホットペッパーAPIの認証キー
	baseURL string // APIのベースURL
}

// NewHotpepperService は新しい HotpepperService を作成するコンストラクタ関数
// apiKey: ホットペッパーAPIの認証キー
// 戻り値: HotpepperServiceのポインタ
func NewHotpepperService(apiKey string) *HotpepperService {
	// 構造体リテラルで初期化し、そのポインタを返す
	return &HotpepperService{
		apiKey:  apiKey,           // 引数で受け取ったAPIキーを設定
		baseURL: hotpepperBaseURL, // 定数で定義されたベースURLを設定
	}
}

// SearchGourmet はグルメサーチAPIを呼び出すメソッド
// params: 検索パラメータ
// 戻り値: 型付きのレスポンスとエラー
func (s *HotpepperService) SearchGourmet(params types.GourmetSearchParams) (*types.GourmetSearchResponse, error) {
	// 必須パラメータのバリデーション
	if params.ServiceArea == "" {
		return nil, fmt.Errorf("service_area is required")
	}

	// url.Values{}でクエリパラメータを管理するマップを作成
	// マップはキーと値のペアを格納するデータ構造
	queryParams := url.Values{}

	// 必須パラメータの設定
	queryParams.Set("key", s.apiKey)                    // APIキー
	queryParams.Set("type", "lite")                     // レスポンスタイプ（軽量版）
	queryParams.Set("format", "json")                   // レスポンス形式（JSON）
	queryParams.Set("service_area", params.ServiceArea) // 都道府県コード

	// オプションパラメータの設定（値が空でない場合のみ追加）
	if params.Address != "" {
		queryParams.Set("address", params.Address) // 住所キーワード
	}
	if params.Genre != "" {
		queryParams.Set("genre", params.Genre) // ジャンルコード
	}
	if params.Keyword != "" {
		queryParams.Set("keyword", params.Keyword) // フリーワード
	}

	// ページングパラメータの設定
	// clampIntで値を許容範囲内に制限（countに関してはAPIの制約に準拠）
	start := clampInt(params.Start, 1, 1000) // 開始位置: 1〜1000
	count := clampInt(params.Count, 1, 100)  // 取得件数: 1〜100
	// strconv.Itoa()で整数を文字列に変換
	queryParams.Set("start", strconv.Itoa(start))
	queryParams.Set("count", strconv.Itoa(count))

	// fmt.Sprintf()で文字列をフォーマット（変数を埋め込んで文字列を作成）
	// queryParams.Encode()でクエリパラメータをURLエンコード
	// 例: "key=xxx&format=json&service_area=SA11"
	apiURL := fmt.Sprintf("%s/gourmet/v1/?%s", s.baseURL, queryParams.Encode())

	// 共通のAPI呼び出しメソッドを使用してリクエストを実行
	body, err := s.fetchAPI(apiURL)
	if err != nil {
		return nil, err
	}

	// JSONをデコードして型付きレスポンスに変換
	var response types.GourmetSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetGenreMaster はジャンルマスタAPIを呼び出すメソッド
// ジャンル一覧（和食、イタリアン、中華など）を取得する
// 戻り値: 型付きのレスポンスとエラー
func (s *HotpepperService) GetGenreMaster() (*types.GenreMasterResponse, error) {
	// クエリパラメータを準備
	queryParams := url.Values{}
	queryParams.Set("key", s.apiKey)  // APIキー
	queryParams.Set("format", "json") // レスポンス形式（JSON）

	// APIのURLを構築
	// /genre/v1/ はジャンルマスタAPIのエンドポイント
	apiURL := fmt.Sprintf("%s/genre/v1/?%s", s.baseURL, queryParams.Encode())

	// 共通のAPI呼び出しメソッドを使用
	body, err := s.fetchAPI(apiURL)
	if err != nil {
		return nil, err
	}

	// JSONをデコードして型付きレスポンスに変換
	var response types.GenreMasterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// fetchAPI はAPIを呼び出してレスポンスを取得する共通メソッド
// 複数のAPIで共通する処理をまとめることでコードの重複を避ける（DRY原則）
// apiURL: 呼び出すAPIの完全なURL
// 戻り値: レスポンスボディのバイト配列とエラー
func (s *HotpepperService) fetchAPI(apiURL string) ([]byte, error) {
	// http.Get()でGETリクエストを送信
	// resp: レスポンス情報、err: エラー情報
	resp, err := http.Get(apiURL)
	// ネットワークエラーなどが発生した場合
	if err != nil {
		// nilは値がないことを表す（エラー時はデータを返さない）
		// %wはエラーをラップして、元のエラー情報を保持する
		return nil, fmt.Errorf("failed to fetch API: %w", err)
	}
	// defer文は関数終了時に実行される処理を予約
	// resp.Body.Close()でレスポンスボディを閉じ、リソースを解放
	defer resp.Body.Close()

	// io.ReadAll()でレスポンスボディを全て読み込む
	// bodyには[]byte型でデータが格納される
	body, err := io.ReadAll(resp.Body)
	// レスポンスの読み込みでエラーが発生した場合
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// HTTPステータスコードが200 OK以外の場合、エラーとして扱う
	// http.StatusOKは定数で値は200
	if resp.StatusCode != http.StatusOK {
		// %dは整数を埋め込むフォーマット指定子
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// 正常な場合、レスポンスボディとnilエラーを返す
	return body, nil
}

// clampInt は値を最小値と最大値の範囲内に制限するヘルパー関数
// value: チェックする値
// min: 最小値
// max: 最大値
// 戻り値: 範囲内に収められた値
func clampInt(value, min, max int) int {
	// 値が最小値より小さい場合、最小値を返す
	if value < min {
		return min
	}
	// 値が最大値より大きい場合、最大値を返す
	if value > max {
		return max
	}
	// 値が範囲内の場合、そのまま返す
	return value
}
