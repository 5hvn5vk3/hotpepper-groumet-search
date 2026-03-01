// このファイルが属するパッケージ名を宣言する。
package middleware

// この行で必要なパッケージを1つ読み込む。
import "net/http"

// この行で関数定義を始める。
func CORS(allowedOrigin string, next http.HandlerFunc) http.HandlerFunc {

	// この時点で関数の結果を返して処理を終了する。
	return func(w http.ResponseWriter, r *http.Request) {

		// HTTPレスポンスへ値を書き込む処理を行う。
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)

		// HTTPレスポンスへ値を書き込む処理を行う。
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		// HTTPレスポンスへ値を書き込む処理を行う。
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 条件を評価し、真のときだけ次の処理を実行する。
		if r.Method == http.MethodOptions {

			// HTTPレスポンスへ値を書き込む処理を行う。
			w.WriteHeader(http.StatusOK)

			// この時点で関数処理を終了する。
			return
			// ここで処理ブロックを終了する。
		}

		// この行は処理の一部として必要な操作を実行する。
		next(w, r)
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}
