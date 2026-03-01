// package service は、外部APIとの通信やビジネスロジックを担うレイヤーです。
// HTTP ハンドラー（presentation層）とデータ型（types層）の橋渡し役を担います。
package service

import (
	// encoding/json: JSON のエンコード・デコードを行う標準パッケージ
	"encoding/json"
	// fmt: 文字列フォーマットやエラー生成に使用する標準パッケージ
	"fmt"
	// io: ストリーム I/O の基本インターフェースを提供する標準パッケージ
	"io"
	// log: 標準出力へのログ出力に使用するパッケージ
	"log"
	// net/http: HTTP クライアント・サーバーを実装する標準パッケージ
	"net/http"
	// net/url: URL のパース・クエリパラメータ組み立てに使用するパッケージ
	"net/url"
	// strconv: 数値と文字列の相互変換に使用するパッケージ
	"strconv"
	// strings: 文字列操作ユーティリティを提供するパッケージ
	"strings"
	// time: 時間・タイムアウト制御に使用するパッケージ
	"time"

	"backend/internal/types"
)

// hotpepperBaseURL は Hot Pepper グルメ Web サービスのベース URL です。
// const（定数）にすることで、誤って書き換えられるのを防ぎます。
const hotpepperBaseURL = "https://webservice.recruit.co.jp/hotpepper"

// HotpepperService はホットペッパー API との通信に必要な設定を保持する構造体です。
// フィールドを小文字（unexported）にすることで、パッケージ外からの直接アクセスを禁止し、
// NewHotpepperService 経由での初期化を強制します。
type HotpepperService struct {
	apiKey  string       // APIキー（認証情報）
	baseURL string       // APIのベースURL（テスト時にモックサーバーに差し替えられるよう分離）
	client  *http.Client // HTTP クライアント（タイムアウト設定付き）
}

// NewHotpepperService は HotpepperService を初期化して返すコンストラクタ関数です。
// Go にはクラスがないため、このような「New〇〇」という名前の関数でインスタンスを生成するのが慣習です。
//
// http.Client の Timeout を設定することで、API が応答しない場合でも
// 指定秒数（ここでは60秒）後に自動的にエラーを返し、ゴルーチンが永遠にブロックされるのを防ぎます。
// Timeout を設定しないと、ネットワーク障害時にリクエストが永遠に待ち続けてしまいます。
func NewHotpepperService(apiKey string) *HotpepperService {

	return &HotpepperService{
		apiKey:  apiKey,
		baseURL: hotpepperBaseURL,
		// &http.Client{...} でポインタを生成し、Timeout フィールドを設定します。
		// time.Second は time.Duration 型の定数で、60 * time.Second = 60秒のタイムアウトです。
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

// SearchGourmet はグルメ検索を行い、結果を返すメソッドです。
// レシーバー (s *HotpepperService) を使うことで、構造体のフィールド（apiKey, client など）に
// アクセスできます。Go のメソッドはこのように関数の先頭にレシーバーを付けて定義します。
//
// 戻り値に error を含めるのは Go のイディオムです。
// 呼び出し側は必ず err != nil をチェックしてエラーハンドリングを行う必要があります。
func (s *HotpepperService) SearchGourmet(params types.GourmetSearchParams) (*types.GourmetSearchResponse, error) {

	// url.Values は map[string][]string の型エイリアスで、
	// クエリパラメータ（?key=value&key2=value2 の部分）を安全に組み立てるための型です。
	// 直接文字列を連結すると特殊文字のエスケープ漏れが起きますが、
	// url.Values を使うと Encode() 時に自動でURLエンコードしてくれます。
	queryParams := url.Values{}

	// Set メソッドで key-value ペアを追加します（同じキーが既にあれば上書きされます）。
	queryParams.Set("key", s.apiKey)
	queryParams.Set("type", "lite+credit_card")
	queryParams.Set("format", "json")

	// 緯度・経度が両方とも指定されている場合のみ位置情報検索パラメータを追加します。
	// Go のゼロ値（float64 の場合は 0）を利用して「未指定」を判定しています。
	if params.Lat != 0 && params.Lng != 0 {
		// strconv.FormatFloat は float64 を文字列に変換する関数です。
		// 第2引数 'f': 指数表記なしの固定小数点形式（例: "35.681236"）
		// 第3引数 -1: 精度を自動で決定（余分なゼロを省略した最短表現）
		// 第4引数 64: 元の値が float64 であることを示す（精度丸めの制御）
		queryParams.Set("lat", strconv.FormatFloat(params.Lat, 'f', -1, 64))
		queryParams.Set("lng", strconv.FormatFloat(params.Lng, 'f', -1, 64))
		if params.Range > 0 {
			// strconv.Itoa は int を文字列に変換します（"Int to ASCII" の略）。
			queryParams.Set("range", strconv.Itoa(params.Range))
		}
	}

	// 各パラメータはオプショナルなので、値が空の場合はクエリに含めません。
	// 空文字をそのまま送ると API がエラーを返すことがあるため、このチェックが重要です。
	if params.Address != "" {
		queryParams.Set("address", params.Address)
	}
	if params.Genre != "" {
		queryParams.Set("genre", params.Genre)
	}
	if params.Keyword != "" {
		queryParams.Set("keyword", params.Keyword)
	}

	// clampInt でページング用パラメータを API の受け入れ範囲内に収めます。
	// ユーザー入力をそのまま使うと範囲外の値が送られる可能性があるため、
	// バリデーション（値の制限）をここで行います。
	start := clampInt(params.Start, 1, 1000)
	count := clampInt(params.Count, 1, 100)

	queryParams.Set("start", strconv.Itoa(start))
	queryParams.Set("count", strconv.Itoa(count))

	// queryParams.Encode() はパラメータを "key=value&key2=value2" 形式に変換します。
	// 特殊文字（日本語など）は自動的に URL エンコードされます（例: "東京" → "%E6%9D%B1%E4%BA%AC"）。
	apiURL := fmt.Sprintf("%s/gourmet/v1/?%s", s.baseURL, queryParams.Encode())

	// fetchAPI は HTTP GET リクエストを行い、レスポンスボディをバイト列で返す内部メソッドです。
	body, err := s.fetchAPI(apiURL)
	if err != nil {
		return nil, err
	}

	// json.Unmarshal は JSON バイト列を Go の構造体に変換（デシリアライズ）します。
	// &response はポインタを渡しています。ポインタを渡さないと値のコピーに書き込まれ、
	// 元の変数に反映されないため、必ずポインタを使います。
	//
	// fmt.Errorf("...: %w", err) の %w 動詞はエラーをラップします。
	// これにより errors.Is や errors.As で元のエラーを検査できるようになります。
	var response types.GourmetSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// API が 200 OK を返してもレスポンス内にエラー情報が含まれる場合があります。
	// ホットペッパー API はエラーをレスポンスボディの JSON に埋め込んで返すため、
	// HTTP ステータスコードだけでなくレスポンス内容も確認する必要があります。
	if len(response.Results.Error) > 0 {
		return nil, convertAPIError(&response.Results.Error[0])
	}

	// &response でポインタを返すことで、大きな構造体のコピーを避けてメモリ効率を高めます。
	return &response, nil
}

// GetGourmetDetail は店舗 ID を指定して1件の詳細情報を取得するメソッドです。
// ID の前後に空白が含まれていてもバリデーションが通ってしまわないよう、
// strings.TrimSpace でトリミング（先頭・末尾の空白除去）してから検証します。
func (s *HotpepperService) GetGourmetDetail(id string) (*types.GourmetSearchResponse, error) {
	// strings.TrimSpace は文字列の先頭・末尾にある空白文字（スペース・タブ・改行）を除去します。
	// ユーザーや外部システムからの入力には予期しない空白が含まれることがあるため、
	// 入力値の正規化（サニタイズ）として必ず行うべき処理です。
	shopID := strings.TrimSpace(id)
	// ID が空の場合はリクエストを送る前にエラーを返します（Fail Fast の原則）。
	// 不正な入力で無駄なAPIリクエストが発生するのを防ぎ、エラーの原因を早期に特定できます。
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

	if len(response.Results.Error) > 0 {
		return nil, convertAPIError(&response.Results.Error[0])
	}

	return &response, nil
}

// GetGenreMaster はジャンルマスターデータ（全ジャンル一覧）を取得するメソッドです。
// マスターデータとは、ジャンル名・コードなど変更頻度が低い基礎データのことです。
// 検索画面のジャンル選択プルダウンなど UI 構築に使用されます。
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

	if len(response.Results.Error) > 0 {
		return nil, convertAPIError(&response.Results.Error[0])
	}

	return &response, nil
}

// fetchAPI は指定された URL に HTTP GET リクエストを送り、
// レスポンスボディをバイト列（[]byte）で返す内部ヘルパーメソッドです。
// 小文字始まり（fetchAPI）にすることでパッケージ外からは呼び出せない unexported メソッドになります。
func (s *HotpepperService) fetchAPI(apiURL string) ([]byte, error) {

	// s.client.Get はHTTP GETリクエストを送信します。
	// s.client は NewHotpepperService でタイムアウト設定済みのため、
	// 応答が遅い場合は自動的にエラーが返ります。
	// net/http のデフォルトクライアント（http.Get）はタイムアウトが無いため、
	// プロダクションコードでは必ず自前の http.Client を使うべきです。
	resp, err := s.client.Get(apiURL)

	if err != nil {
		// fmt.Errorf の %w 動詞でエラーをラップすることで、
		// 呼び出し元が errors.Is(err, someSpecificError) で元のエラーを判別できます。
		// %v と違い、%w はエラーチェーンを保持するのが特徴です。
		return nil, fmt.Errorf("failed to fetch API: %w", err)
	}

	// defer はその関数が return する直前に実行される遅延実行の仕組みです。
	// HTTP レスポンスの Body は必ず Close() しなければなりません。
	// Close を忘れると TCP コネクションが解放されず、コネクションリークが発生します。
	// defer を使うことで、複数の return パスがある場合でも確実に Close されます。
	defer resp.Body.Close()

	// io.ReadAll はレスポンスボディ（io.Reader インターフェース）を
	// 全て読み込んでバイト列（[]byte）に変換します。
	// ストリーム（逐次読み込み）を一括でメモリに展開するため、
	// 非常に大きなレスポンスには不向きですが、通常の API 応答であれば問題ありません。
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// HTTPステータスコードが 200 OK 以外の場合はエラーとして扱います。
	// http.StatusOK は定数 200 のエイリアスで、マジックナンバーを避けるために使います。
	// 401 Unauthorized（APIキー誤り）や 429 Too Many Requests（レート制限）なども
	// ここでキャッチされ、エラーメッセージにレスポンスボディが含まれます。
	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// convertAPIError はホットペッパーAPIのエラーをHotpepperAPIErrorに変換する。
// 生のエラーメッセージは情報漏洩防止のためログのみに出力し、外部には返さない。
//
// セキュリティ上の理由から、API の生のエラーメッセージ（例: "Invalid API key"）を
// HTTP レスポンスでそのまま返すと、攻撃者にシステム情報を与えてしまいます。
// そのため、詳細はサーバーサイドのログにのみ記録し、クライアントにはエラーコードのみを返します。
// これをエラーの「情報漏洩防止（Information Disclosure 対策）」と言います。
//
// 呼び出し元で errors.As(err, &target) を使うことで、
// *types.HotpepperAPIError 型に変換してエラーコードを取り出せます。
func convertAPIError(apiErr *types.APIError) error {
	log.Printf("hotpepper API error: code=%d, message=%s", apiErr.Code, apiErr.Message)
	return &types.HotpepperAPIError{Code: apiErr.Code}
}

// clampInt は value を [min, max] の範囲内に収めて返すユーティリティ関数です。
// 「clamp（クランプ）」とは値を指定範囲に挟み込む操作のことです。
// ユーザー入力が API の許容範囲を超えないようにするバリデーション処理に使います。
//
// 例: clampInt(0, 1, 100) → 1（最小値に切り上げ）
//     clampInt(200, 1, 100) → 100（最大値に切り下げ）
//     clampInt(50, 1, 100) → 50（範囲内なのでそのまま）
func clampInt(value, min, max int) int {

	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}
