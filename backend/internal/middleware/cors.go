package middleware

import "net/http"

// CORS は CORS ヘッダーを追加するミドルウェア関数
// CORS (Cross-Origin Resource Sharing) は異なるオリジン（ドメイン）からのリクエストを許可する仕組み
// ミドルウェアとは、リクエスト処理の前後に共通処理を挟み込むための仕組み
// allowedOrigin: 許可するオリジン（例: "http://localhost:5173"）。環境変数 ALLOWED_ORIGIN から渡す
func CORS(allowedOrigin string, next http.HandlerFunc) http.HandlerFunc {
	// 新しいハンドラー関数を返す（クロージャー）
	return func(w http.ResponseWriter, r *http.Request) {
		// Access-Control-Allow-Origin: 引数で受け取ったオリジンのみを許可
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)

		// Access-Control-Allow-Methods: 許可するHTTPメソッドを指定
		// GETとOPTIONSメソッドのみを許可
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		// Access-Control-Allow-Headers: クライアントが送信できるヘッダーを指定
		// Content-Typeヘッダーの送信を許可
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// OPTIONSリクエスト（プリフライトリクエスト）の処理
		// ブラウザが実際のリクエストの前に送信する確認用のリクエスト
		if r.Method == http.MethodOptions {
			// 200 OKを返してプリフライトリクエストを完了
			w.WriteHeader(http.StatusOK)
			// プリフライトリクエストの場合は、次のハンドラーを実行せずに終了
			return
		}

		// 通常のリクエストの場合は、次のハンドラー（本来の処理）を実行
		next(w, r)
	}
}
