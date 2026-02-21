package handler

import (
	"encoding/json" // JSONのエンコード・デコードを提供
	"net/http"      // HTTPサーバーとクライアント機能を提供
	"strconv"       // 文字列と数値の変換を行うパッケージ

	"backend/internal/service" // サービス層をインポート
	"backend/internal/types"   // リクエスト・レスポンスの型定義をインポート
)

// ValidationError はバリデーション（入力検証）エラーを表すカスタムエラー型
type ValidationError struct {
	Message string // エラーメッセージを保持
}

// Error はerrorインターフェースを実装するメソッド
// このメソッドを実装することで、ValidationErrorをerror型として扱える
func (e *ValidationError) Error() string {
	// エラーメッセージを返す
	return e.Message
}

// GourmetHandler はグルメサーチAPIのリクエストを処理する構造体
type GourmetHandler struct {
	// serviceフィールド：HotpepperServiceのインスタンスへのポインタを保持
	service *service.HotpepperService
}

// NewGourmetHandler は新しい GourmetHandler を作成するコンストラクタ関数
func NewGourmetHandler(service *service.HotpepperService) *GourmetHandler {
	// GourmetHandler構造体を初期化し、そのポインタを返す
	return &GourmetHandler{service: service}
}

// Handle はHTTPリクエストを処理するメインメソッド
// w: レスポンスを書き込むためのResponseWriter
// r: クライアントからのリクエスト情報
func (h *GourmetHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// HTTPメソッドがGETであることを確認
	// このAPIはGETメソッドのみを受け付ける
	if r.Method != http.MethodGet {
		// 405エラー：許可されていないメソッドが使用された
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// リクエストからクエリパラメータを取得し、バリデーション（検証）を実行
	// parseParamsメソッドでパラメータを解析し、GourmetSearchParams構造体に変換
	params, err := h.parseParams(r)
	// パラメータの解析でエラーが発生した場合
	if err != nil {
		// 400エラー：不正なリクエスト（パラメータが不足または無効）
		// err.Error()でエラーメッセージを文字列として取得
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// サービス層のメソッドを呼び出して、ホットペッパーAPIから店舗データを検索
	response, err := h.service.SearchGourmet(params)
	// API呼び出しでエラーが発生した場合
	if err != nil {
		// 500エラー：サーバー内部エラー
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	// レスポンスをJSONにエンコード
	body, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// 正常にデータを取得できた場合、レスポンスを返す
	// Content-Typeヘッダーを設定：レスポンスがJSON形式であることを示す
	w.Header().Set("Content-Type", "application/json")
	// HTTPステータスコードを200（成功）に設定
	w.WriteHeader(http.StatusOK)
	// JSONデータをレスポンスボディに書き込む
	w.Write(body)
}

// parseParams はHTTPリクエストからクエリパラメータを解析するメソッド
// 戻り値：GourmetSearchParams構造体とエラー（エラーがない場合はnil）
func (h *GourmetHandler) parseParams(r *http.Request) (types.GourmetSearchParams, error) {
	// r.URL.Query()でURLからクエリパラメータを取得
	// 例：/api/gourmet?service_area=SA11&keyword=居酒屋
	query := r.URL.Query()

	// 必須パラメータ（service_area）のチェック
	// Get()メソッドで指定されたキーの値を取得（存在しない場合は空文字列）
	serviceArea := query.Get("service_area")
	// service_areaが空の場合、エラーを返す
	if serviceArea == "" {
		// GourmetSearchParams{}は空の構造体を返す
		// &ValidationError{...}はValidationErrorのポインタを返す
		return types.GourmetSearchParams{}, &ValidationError{Message: "service_area is required"}
	}

	// オプショナルなパラメータの解析
	// parseIntOrDefaultで文字列を整数に変換、失敗時はデフォルト値を使用
	start := parseIntOrDefault(query.Get("start"), 1)  // 開始位置（デフォルト：1）
	count := parseIntOrDefault(query.Get("count"), 20) // 取得件数（デフォルト：20）

	// GourmetSearchParams構造体を作成して返す
	// 構造体リテラルを使用してフィールドに値を設定
	return types.GourmetSearchParams{
		ServiceArea: serviceArea,          // 必須：都道府県コード
		Address:     query.Get("address"), // オプション：住所
		Genre:       query.Get("genre"),   // オプション：ジャンルコード
		Keyword:     query.Get("keyword"), // オプション：キーワード
		Start:       start,                // ページング：開始位置
		Count:       count,                // ページング：取得件数
	}, nil // nilはエラーなしを意味する
}

// parseIntOrDefault は文字列を整数に変換し、失敗時はデフォルト値を返すヘルパー関数
// value: 変換する文字列
// fallback: 変換失敗時に返すデフォルト値
func parseIntOrDefault(value string, fallback int) int {
	// 値が空文字列の場合、デフォルト値を返す
	if value == "" {
		return fallback
	}
	// strconv.Atoi()で文字列を整数に変換
	// v: 変換された整数値、err: エラー情報
	if v, err := strconv.Atoi(value); err == nil {
		// 変換成功（err == nil）の場合、変換された値を返す
		return v
	}
	// 変換失敗の場合、デフォルト値を返す
	return fallback
}
