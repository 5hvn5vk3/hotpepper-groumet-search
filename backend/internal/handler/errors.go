package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"backend/internal/types"
)

// regexp.MustCompile はプログラム起動時に正規表現をコンパイルしてパターンオブジェクトを作る。
// var（パッケージ変数）として定義することで、リクエストごとではなく1回だけコンパイルされ、効率的。
// MustCompile は正規表現が不正な場合に panic を起こすため、起動時に問題が検出できる。
//
// このパターンは URLクエリ文字列中の "key=<値>" の部分（APIキー）にマッチする。
//   例: "?key=abc123&keyword=寿司" → "?key=[REDACTED]&keyword=寿司"
//
// キャプチャグループ `([?&]key=)` で "key=" の手前の区切り文字も一緒に取り、
// 置換後に `${1}` で復元することで、クエリ文字列の構造を壊さない。
var apiKeyQueryPattern = regexp.MustCompile(`([?&]key=)[^&\s]+`)

// handleServiceError はservice層から返ってきたエラーを分析し、
// クライアント（フロントエンド）に返すべき適切なHTTPステータスコードと
// ユーザー向けメッセージをJSON形式で書き込む関数。
//
// 【セキュリティ上の重要な設計方針】
// ホットペッパーAPIが返す生のエラーメッセージには、内部URL・パラメータ・
// スタックトレースなどの機密情報が含まれる可能性がある。
// そのため、外部（フロント）には「サービスが一時的に利用できません」という
// 汎用メッセージだけを返し、詳細はサーバーのログにのみ記録する。
//
// handleServiceError はservice層のエラーを判別し、適切なHTTPステータスと
// ユーザー向けメッセージ（自前定義）をJSONで返す。
// ホットペッパーAPIの生のエラーメッセージはログのみに出力し、フロントには返さない。
func handleServiceError(w http.ResponseWriter, err error) {
	// errors.As は、エラーチェーンの中から指定した型（ここでは *types.HotpepperAPIError）の
	// エラーを探し出す標準ライブラリ関数。
	// 単純な型アサーション（err.(*types.HotpepperAPIError)）と違い、
	// errors.Wrap などでラップされた深いエラーも透過的に見つけられる。
	// 第2引数にはポインタのポインタ（**型）を渡す点に注意。
	var apiErr *types.HotpepperAPIError
	if errors.As(err, &apiErr) {
		log.Printf("hotpepper API error handled: code=%d", apiErr.Code)
		switch apiErr.Code {
		case 3000:
			// コード3000: ホットペッパーAPIへ送ったリクエストパラメータが不正。
			// これはクライアント（ユーザー）側が修正できる問題なので、
			// 400 Bad Request を返す。「あなたのリクエストに問題があります」の意味。
			writeErrorJSON(w, http.StatusBadRequest, "検索条件が正しくありません")
			return
		case 1000:
			// コード1000: ホットペッパー側のサーバー障害。
			// 502 Bad Gateway は「このサーバーは正常だが、上流（外部API）が壊れている」を示す。
			// 503（Service Unavailable）と似ているが、502 は上流依存の障害に使うのが慣習。
			writeErrorJSON(w, http.StatusBadGateway, "サービスが一時的に利用できません")
			return
		case 2000:
			// コード2000: APIキーまたはIP認証エラー。
			// これはバックエンドの設定ミス（無効なAPIキーなど）が原因であり、
			// ユーザーには解決できないため 500 Internal Server Error を返す。
			// フロントに「認証エラー」と伝えると攻撃のヒントになるので汎用メッセージにする。
			writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
			return
		}
	}
	// 上記いずれのエラーコードにも該当しない場合（ネットワーク障害・未知のエラーなど）は
	// 500 Internal Server Error を返す。ログにはエラー詳細を記録するが、
	// その際に maskAPIKeyInLog でAPIキーをマスクしてログ漏洩を防ぐ。
	log.Printf("internal server error: %s", maskAPIKeyInLog(err))
	writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
}

// maskAPIKeyInLog はエラーメッセージ文字列の中に含まれるAPIキーをマスクする。
// サーバーログはオペレータが確認するため、APIキーが平文でログに残ると
// 漏洩リスクがある。"key=実際のキー値" を "key=[REDACTED]" に置き換えることで
// セキュリティを保ちながら、ログの他の情報（エンドポイント・パラメータ等）は保持できる。
func maskAPIKeyInLog(err error) string {
	if err == nil {
		return ""
	}

	// ReplaceAllString は正規表現にマッチする全箇所を第2引数の文字列で置換する。
	// `${1}` はキャプチャグループ1番（"?key=" or "&key=" の部分）を指す後方参照。
	return apiKeyQueryPattern.ReplaceAllString(err.Error(), "${1}[REDACTED]")
}

// writeErrorJSON はHTTPレスポンスにエラー情報をJSON形式で書き込むユーティリティ関数。
// 全ハンドラーから共通利用することで、エラーレスポンスの形式を統一できる。
//
// 【レスポンスボディのフォーマット例】
//   {"error": {"message": "検索条件が正しくありません"}}
//
// w.Header().Set は必ず w.WriteHeader より前に呼ぶこと。
// WriteHeader を呼んだ後にヘッダーを設定しても反映されない（Goの仕様）。
func writeErrorJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// json.NewEncoder(w).Encode は構造体・マップをJSONに変換しつつ直接 w（ResponseWriter）へ書き込む。
	// json.Marshal + w.Write の2ステップを1行で行う省略記法。
	json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"message": message},
	})
}
