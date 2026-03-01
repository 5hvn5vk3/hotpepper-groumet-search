// package middleware はHTTPハンドラーに共通処理を付加するミドルウェアをまとめるパッケージです。
// ミドルウェアとは「本来の処理の前後に割り込む薄いレイヤー」のことで、
// 認証・ロギング・CORSヘッダー付与など、複数エンドポイントで使い回す処理を
// 一か所に集約するために使われます。
package middleware

import "net/http"

// CORS は「クロスオリジンリソース共有 (Cross-Origin Resource Sharing)」のための
// ミドルウェア関数です。
//
// --- CORSとは ---
// ブラウザには「同一オリジンポリシー (Same-Origin Policy)」というセキュリティ制約があり、
// デフォルトでは異なるオリジン（スキーム+ホスト+ポートの組み合わせ）への
// JavaScriptからのHTTPリクエストをブロックします。
// 例えばフロントエンド (http://localhost:3000) からバックエンド (http://localhost:8080) へ
// fetch() でリクエストを送る場合、オリジンが異なるためブロックされます。
// これを解除するために、サーバー側がレスポンスヘッダーで「このオリジンからのアクセスを許可する」
// と明示的に宣言する仕組みがCORSです。
//
// --- 関数シグネチャについて ---
// この関数は "http.HandlerFunc を受け取り、http.HandlerFunc を返す" 高階関数です。
// 引数 next は「CORSチェック後に実行する本来のハンドラー」です。
// 戻り値は「CORSヘッダーを付与してから next を呼ぶ新しいハンドラー」です。
// このパターンにより、ルーティング側でラップするだけでCORSを適用できます。
func CORS(allowedOrigin string, next http.HandlerFunc) http.HandlerFunc {

	// ここで「クロージャ」を返しています。
	// クロージャとは「外側のスコープの変数（ここでは allowedOrigin と next）を
	// 捕捉した関数」のことです。
	// 返された関数はどこで呼ばれても allowedOrigin と next を参照し続けます。
	// これによって CORS() を呼んだ時点で設定値が束縛され、
	// 実際にリクエストが来るたびに同じ設定で動作します。
	return func(w http.ResponseWriter, r *http.Request) {

		// --- CORSレスポンスヘッダーの設定 ---
		// w.Header().Set() で HTTP レスポンスヘッダーをセットします。
		// これらのヘッダーをブラウザが読み取り、クロスオリジンリクエストを許可するかどうかを判断します。

		// Access-Control-Allow-Origin: このオリジンからのリクエストを許可することをブラウザに伝えます。
		// 特定のオリジン文字列を指定することで、それ以外のオリジンからはアクセスを拒否できます。
		// ワイルドカード "*" を使うと全オリジンを許可しますが、セキュリティリスクが高まります。
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)

		// Access-Control-Allow-Methods: 許可するHTTPメソッドを指定します。
		// このAPIはデータ取得のみなので GET と OPTIONS だけを許可しています。
		// POST / PUT / DELETE などの変更系メソッドはここに追加しない限りブロックされます。
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		// Access-Control-Allow-Headers: クロスオリジンリクエストで送信を許可するリクエストヘッダーを指定します。
		// ここでは Content-Type ヘッダーの送信を許可しています。
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// --- プリフライトリクエスト (Preflight Request) の処理 ---
		// ブラウザは本番リクエストを送る前に、まず OPTIONS メソッドで「このリクエストを
		// 送っていいか？」とサーバーに確認するリクエストを自動的に送ります。
		// これを「プリフライトリクエスト」と呼びます。
		// サーバーは上で設定したCORSヘッダーと 200 OK を返すだけでよく、
		// 本来のビジネスロジック（next）を実行する必要はありません。
		if r.Method == http.MethodOptions {

			// 200 OK を返してプリフライトリクエストへの応答を完了します。
			// http.StatusOK は整数値 200 の定数エイリアスで、可読性のために使います。
			w.WriteHeader(http.StatusOK)

			// return で関数を即座に終了し、next（本来のハンドラー）は呼びません。
			return
		}

		// OPTIONS 以外のリクエスト（実際のGETリクエストなど）は
		// CORSヘッダーを付与した後、本来のハンドラーに処理を委譲します。
		// next(w, r) で w（レスポンス書き込み先）と r（リクエスト情報）をそのまま渡します。
		next(w, r)
	}
}
