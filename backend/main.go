// package main はGoプログラムのエントリーポイントとなるパッケージです。
// Goでは必ず1つの "main" パッケージに "main" 関数を定義することで実行可能バイナリが作られます。
package main

import (
	"log"
	"net/http"
	"os"

	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/service"
)

// main はプログラム起動時に最初に呼ばれる関数です。
// サーバーの設定・起動処理をすべてここに記述しています。
func main() {

	// os.Getenv は「環境変数」を読み取るGoの標準関数です。
	// 環境変数とは、OSレベルで設定されたキー=値のペアで、
	// コードを変更せずに動作を切り替えられる仕組みです（例: 本番/開発環境の切り替え）。
	// ここではサーバーが待ち受けるポート番号を環境変数から取得しています。
	port := os.Getenv("PORT")

	// 環境変数が設定されていない場合（空文字列が返る）はデフォルト値を使います。
	// こうすることでローカル開発中は何も設定せず 8080 ポートで動かせます。
	if port == "" {
		port = "8080"
	}

	// ホットペッパーAPIを呼び出すための認証キーを環境変数から取得します。
	// APIキーのような秘密情報はソースコードに直接書かず、
	// 環境変数で管理するのがセキュリティ上の基本的なベストプラクティスです。
	apiKey := os.Getenv("HOTPEPPER_API_KEY")

	// APIキーは必須のため、未設定ならプログラムを即座に終了させます。
	// log.Fatal は「エラーメッセージをログ出力してからos.Exit(1)で終了する」関数です。
	// os.Exit(1) は「異常終了」を表す終了コードで、CI/CDツールやDockerなどが
	// この値を見てエラー発生を検知できます。
	if apiKey == "" {
		log.Fatal("HOTPEPPER_API_KEY environment variable is required")
	}

	// CORSで許可するフロントエンドのオリジン（スキーム+ホスト+ポート）を環境変数から取得します。
	// オリジンの詳細は middleware/cors.go のコメントを参照してください。
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		// 未設定時はローカル開発用のNext.jsデフォルトオリジンをフォールバックとして使います。
		allowedOrigin = "http://localhost:3000"
	}

	// --- 依存性の注入 (Dependency Injection) ---
	// サービス層・ハンドラー層の順にオブジェクトを生成し、
	// 上位レイヤーが下位レイヤーのインスタンスを受け取る構造にしています。
	// こうすることでテスト時にモックに差し替えやすくなります。

	// ホットペッパーAPIとの通信を担うサービスを生成します。
	hotpepperService := service.NewHotpepperService(apiKey)

	// 各エンドポイントへのリクエストを処理するハンドラーを生成します。
	// ハンドラーはサービスに依存するため、サービスを引数として渡します。
	gourmetHandler := handler.NewGourmetHandler(hotpepperService)
	gourmetDetailHandler := handler.NewGourmetDetailHandler(hotpepperService)

	genreHandler := handler.NewGenreHandler(hotpepperService)

	// --- ルーティング設定 ---
	// http.NewServeMux は「URLパスとハンドラーを対応付ける」ルーターです。
	// リクエストが来たとき、登録されたパスに一致するハンドラーを自動で呼び出します。
	mux := http.NewServeMux()

	// mux.HandleFunc でURLパスとハンドラー関数を紐付けます。
	// ここでは直接ハンドラーを渡すのではなく、middleware.CORS() でラップしています。
	// これが「ミドルウェアパターン」と呼ばれる手法で、
	// 「本来の処理の前後に共通処理（ここではCORSヘッダー付与）を挟む」設計です。
	// 詳細は middleware/cors.go を参照してください。
	mux.HandleFunc("/api/gourmet", middleware.CORS(allowedOrigin, gourmetHandler.Handle))
	mux.HandleFunc("/api/gourmet/detail", middleware.CORS(allowedOrigin, gourmetDetailHandler.Handle))

	mux.HandleFunc("/api/genre", middleware.CORS(allowedOrigin, genreHandler.Handle))

	// log.Printf はフォーマット文字列を使ってログを出力します。
	// %s は文字列プレースホルダーで、port の値が埋め込まれます。
	log.Printf("Server starting on port %s", port)

	// http.ListenAndServe はTCPソケットを開いてHTTPリクエストの待ち受けを開始します。
	// この関数は正常時には返らず（ブロックし続け）、エラー発生時のみ error を返します。
	// そのため if err := ...; err != nil パターンで即座にエラーを検知しています。
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
