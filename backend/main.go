// このファイルが属するパッケージ名を宣言する。
package main

// ここから使用する外部パッケージの列挙を始める。
import (
	// この行は処理の一部として必要な操作を実行する。
	"log"
	// この行は処理の一部として必要な操作を実行する。
	"net/http"
	// この行は処理の一部として必要な操作を実行する。
	"os"

	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/handler"
	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/middleware"
	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/service"
	// このブロックや引数リストを閉じる。
)

// この行で関数定義を始める。
func main() {

	// この行で新しい変数を宣言しつつ値を代入する。
	port := os.Getenv("PORT")

	// 条件を評価し、真のときだけ次の処理を実行する。
	if port == "" {
		// この行で変数やフィールドへ値を代入する。
		port = "8080"
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	apiKey := os.Getenv("HOTPEPPER_API_KEY")

	// 条件を評価し、真のときだけ次の処理を実行する。
	if apiKey == "" {
		// デバッグや障害調査のためにログを出力する。
		log.Fatal("HOTPEPPER_API_KEY environment variable is required")
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	// 条件を評価し、真のときだけ次の処理を実行する。
	if allowedOrigin == "" {
		// この行で変数やフィールドへ値を代入する。
		allowedOrigin = "http://localhost:3000"
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	hotpepperService := service.NewHotpepperService(apiKey)

	// この行で新しい変数を宣言しつつ値を代入する。
	gourmetHandler := handler.NewGourmetHandler(hotpepperService)
	// この行で新しい変数を宣言しつつ値を代入する。
	gourmetDetailHandler := handler.NewGourmetDetailHandler(hotpepperService)

	// この行で新しい変数を宣言しつつ値を代入する。
	genreHandler := handler.NewGenreHandler(hotpepperService)

	// この行で新しい変数を宣言しつつ値を代入する。
	mux := http.NewServeMux()

	// この行は処理の一部として必要な操作を実行する。
	mux.HandleFunc("/api/gourmet", middleware.CORS(allowedOrigin, gourmetHandler.Handle))
	// この行は処理の一部として必要な操作を実行する。
	mux.HandleFunc("/api/gourmet/detail", middleware.CORS(allowedOrigin, gourmetDetailHandler.Handle))

	// この行は処理の一部として必要な操作を実行する。
	mux.HandleFunc("/api/genre", middleware.CORS(allowedOrigin, genreHandler.Handle))

	// デバッグや障害調査のためにログを出力する。
	log.Printf("Server starting on port %s", port)

	// 条件を評価し、真のときだけ次の処理を実行する。
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		// デバッグや障害調査のためにログを出力する。
		log.Fatal(err)
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}
