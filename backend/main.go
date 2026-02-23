// main パッケージ：アプリケーションのエントリーポイント
package main

import (
	"log"      // ログ出力用のパッケージ
	"net/http" // HTTPサーバーとクライアント機能を提供するパッケージ
	"os"       // 環境変数などのOS機能を提供するパッケージ

	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/service"
)

// main関数：プログラムの実行開始地点
func main() {
	// 環境変数からポート番号を取得
	// os.Getenv()は指定された環境変数の値を文字列で返す
	port := os.Getenv("PORT")
	// ポート番号が設定されていない場合、デフォルトの8080を使用
	if port == "" {
		port = "8080"
	}

	// ホットペッパーAPIのキーを環境変数から取得
	apiKey := os.Getenv("HOTPEPPER_API_KEY")
	// APIキーが設定されていない場合、エラーメッセージを出力してプログラムを終了
	// log.Fatal()はメッセージを出力後、os.Exit(1)を呼び出してプログラムを終了する
	if apiKey == "" {
		log.Fatal("HOTPEPPER_API_KEY environment variable is required")
	}

	// 許可するオリジンを環境変数から取得
	// 未設定の場合はローカル開発用のデフォルト値を使用
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	// サービス層の初期化
	// ホットペッパーAPIとの通信を担当するサービスオブジェクトを作成
	hotpepperService := service.NewHotpepperService(apiKey)

	// ハンドラーの初期化
	// グルメサーチAPIのリクエストを処理するハンドラーを作成
	gourmetHandler := handler.NewGourmetHandler(hotpepperService)
	// ジャンルマスタAPIのリクエストを処理するハンドラーを作成
	genreHandler := handler.NewGenreHandler(hotpepperService)

	// ルーターのセットアップ
	// http.NewServeMux()は新しいHTTPリクエストマルチプレクサ（ルーター）を作成
	// マルチプレクサはURLパスに応じて適切なハンドラーにリクエストを振り分ける
	mux := http.NewServeMux()
	// "/api/gourmet"へのリクエストをグルメハンドラーに紐付け
	// CORSMiddlewareでラップすることで、クロスオリジンリクエストを許可
	mux.HandleFunc("/api/gourmet", middleware.CORS(allowedOrigin, gourmetHandler.Handle))
	// "/api/genre"へのリクエストをジャンルハンドラーに紐付け
	mux.HandleFunc("/api/genre", middleware.CORS(allowedOrigin, genreHandler.Handle))

	// サーバー起動のログを出力
	log.Printf("Server starting on port %s", port)
	// HTTPサーバーを起動
	// ":"とポート番号を連結してアドレスを作成（例：":8080"）
	// muxをハンドラーとして指定し、すべてのリクエストをmuxで処理
	// エラーが発生した場合はログに出力してプログラムを終了
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
