package handler

import (
	"encoding/json" // JSONのエンコード・デコードを提供
	"net/http"      // HTTPクライアント・サーバー機能を提供

	"backend/internal/service" // サービス層をインポート
)

// GenreHandler はジャンルマスタAPIのリクエストを処理する構造体
// 構造体（struct）は関連するデータをグループ化するための型
type GenreHandler struct {
	// serviceフィールド：HotpepperServiceへのポインタを保持
	// *（アスタリスク）はポインタ型を表す（メモリアドレスを保持）
	service *service.HotpepperService
}

// NewGenreHandler は新しい GenreHandler を作成するコンストラクタ関数
// Goにはコンストラクタの概念がないため、Newで始まる関数を慣例的に使用
func NewGenreHandler(service *service.HotpepperService) *GenreHandler {
	// GenreHandler構造体を作成し、そのポインタを返す
	// &はアドレス演算子で、値のメモリアドレスを取得する
	return &GenreHandler{service: service}
}

// Handle はHTTPリクエストを処理するメソッド
// (h *GenreHandler)はレシーバー：このメソッドがGenreHandler型に属することを示す
// w http.ResponseWriter：クライアントへのレスポンスを書き込むためのインターフェース
// r *http.Request：クライアントからのリクエスト情報を保持する構造体
func (h *GenreHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// リクエストメソッドがGETでない場合、エラーを返す
	// http.MethodGetは定数で"GET"を表す
	if r.Method != http.MethodGet {
		// http.Error()はエラーメッセージとステータスコードをクライアントに返す
		// 405 Method Not Allowedは許可されていないHTTPメソッドが使用されたことを示す
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		// return文で関数を終了し、これ以降の処理を実行しない
		return
	}

	// サービス層を呼び出してジャンルマスタデータを取得
	// h.serviceはこのハンドラーが保持するHotpepperServiceのインスタンス
	response, err := h.service.GetGenreMaster()
	// エラーチェック：errがnilでない場合、エラーが発生している
	if err != nil {
		// 500 Internal Server Errorはサーバー側でエラーが発生したことを示す
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	// レスポンスをJSONにエンコード
	body, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// レスポンスヘッダーを設定
	// Content-Type: application/jsonはレスポンスがJSON形式であることを示す
	w.Header().Set("Content-Type", "application/json")
	// HTTPステータスコードを200 OK（成功）に設定
	w.WriteHeader(http.StatusOK)
	// レスポンスボディ（JSON データ）をクライアントに書き込む
	w.Write(body)
}
