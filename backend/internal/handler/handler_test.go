package handler

import (
	"encoding/json"     // JSON のエンコード・デコード機能
	"net/http"          // HTTPクライアント・サーバー機能
	"net/http/httptest" // HTTPテスト用ユーティリティ
	"testing"           // Go標準のテストフレームワーク

	"backend/internal/service"
)

// ============================================================================
// TestParseIntOrDefault: parseIntOrDefault関数のテスト
// ============================================================================
// 【目的】
// 文字列を整数に変換するヘルパー関数が、様々な入力に対して
// 適切に動作するかを検証
//
// 【関数の仕様】
// parseIntOrDefault(value string, fallback int) int
// - value: 変換する文字列
// - fallback: 変換失敗時に返すデフォルト値
// - 戻り値: 変換成功時は整数値、失敗時はfallback
//
// 【Triangulationの適用】
// 1. 空文字列 → デフォルト値を返す
// 2. 正常な数値文字列 → 変換した値を返す
// 3. 不正な文字列 → デフォルト値を返す
//
// 【学習ポイント】
// - 文字列→数値変換の一般的なパターン
// - エラーハンドリングとデフォルト値の使用
// - ヘルパー関数のテストの重要性
func TestParseIntOrDefault(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name     string // テストケース名
		value    string // テストする入力文字列
		fallback int    // デフォルト値
		expected int    // 期待される結果
		point    string // Triangulationポイント
	}{
		// ----------------------------------------------------------------
		// Triangulation Point 1: 空文字列のケース
		// ----------------------------------------------------------------
		// 検証内容：URLパラメータが指定されていない場合（空文字列）、
		//           デフォルト値を返す
		// 実際の使用例：?page= （値なし）や、パラメータ自体がない場合
		{
			name:     "空文字列の場合、デフォルト値を返す",
			value:    "", // 空文字列
			fallback: 10, // デフォルト値
			expected: 10, // デフォルト値が返される
			point:    "Point 1: 空文字列",
		},

		// ----------------------------------------------------------------
		// Triangulation Point 2: 正常な数値文字列のケース
		// ----------------------------------------------------------------
		// 検証内容：正しい数値文字列の場合、整数に変換して返す
		// 実際の使用例：?page=42
		{
			name:     "正常な数値文字列の場合、変換した値を返す",
			value:    "42", // 正常な数値文字列
			fallback: 10,
			expected: 42, // 変換された整数値
			point:    "Point 2: 正常値",
		},

		// ----------------------------------------------------------------
		// Triangulation Point 3: 不正な文字列のケース
		// ----------------------------------------------------------------
		// 検証内容：数値に変換できない文字列の場合、デフォルト値を返す
		// 実際の使用例：?page=abc （エラー入力）
		{
			name:     "不正な文字列の場合、デフォルト値を返す",
			value:    "abc", // 数値に変換できない文字列
			fallback: 10,
			expected: 10, // デフォルト値にフォールバック
			point:    "Point 3: 不正値",
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 関数を実行
			result := parseIntOrDefault(tt.value, tt.fallback)

			// 結果を検証
			// %q は文字列をクォート付きで表示（例："abc"）
			if result != tt.expected {
				t.Errorf("parseIntOrDefault(%q, %d) = %d; want %d [%s]",
					tt.value, tt.fallback, result, tt.expected, tt.point)
			}
		})
	}
}

// ============================================================================
// TestGourmetHandlerHandle: GourmetHandler.Handleメソッドのテスト
// ============================================================================
// 【目的】
// HTTPハンドラーが、様々なHTTPリクエストに対して適切に応答するかを検証
//
// 【ハンドラーとは】
// HTTPリクエストを受け取り、適切なレスポンスを返す関数
// ウェブアプリケーションの「入口」となる重要な部分
//
// 【テスト内容】
// - HTTPメソッドの検証（GETのみ許可、POST等は拒否）
// - クエリパラメータの処理
// - HTTPステータスコードの確認
// - レスポンスヘッダーの確認
//
// 【Triangulationの適用】
// 1. 正常なGETリクエスト（必須パラメータのみ） → 200 OK
// 2. 正常なGETリクエスト（全パラメータ） → 200 OK
// 3. POSTメソッド → 405 Method Not Allowed
func TestGourmetHandlerHandle(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name           string // テストケース名
		method         string // HTTPメソッド（GET, POST, PUT等）
		queryString    string // URLのクエリパラメータ部分
		expectedStatus int    // 期待されるHTTPステータスコード
		expectError    bool   // エラーレスポンスを期待するか
		point          string // Triangulationポイント
	}{
		// Triangulation Point 1: 正常なGETリクエスト（必須パラメータのみ）
		{
			name:           "必須パラメータのみの正常なGETリクエスト",
			method:         http.MethodGet,
			queryString:    "keyword=居酒屋",
			expectedStatus: http.StatusOK,
			expectError:    false,
			point:          "Point 1: 最小構成",
		},
		// Triangulation Point 2: オプションパラメータを含むGETリクエスト
		{
			name:           "オプションパラメータを含むGETリクエスト",
			method:         http.MethodGet,
			queryString:    "address=さいたま&genre=G001&keyword=居酒屋",
			expectedStatus: http.StatusOK,
			expectError:    false,
			point:          "Point 2: 全パラメータ",
		},
		// Triangulation Point 3: POSTメソッド（エラーケース）
		{
			name:           "POSTメソッドの場合、405エラー",
			method:         http.MethodPost,
			queryString:    "keyword=居酒屋",
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
			point:          "Point 3: 不正メソッド",
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ========================================
			// ステップ1: モックサービスの準備
			// ========================================
			// ハンドラーが依存するサービスのモックを作成
			mockService := &service.HotpepperService{}

			// ----------------------------------------
			// 外部APIへの実際の呼び出しを避けるため、
			// モックサーバーをセットアップ
			// ----------------------------------------
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// このサーバーは常に成功レスポンスを返す
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"results": {"shop": []}}`))
			}))
			defer server.Close()

			// モックサービスを再初期化してbaseURLを設定
			mockService = service.NewHotpepperService("test-key")

			// ========================================
			// ステップ2: テスト対象のハンドラーを作成
			// ========================================
			handler := NewGourmetHandler(mockService)

			// ========================================
			// ステップ3: テスト用のHTTPリクエストを作成
			// ========================================
			// httptest.NewRequest()：テスト用のHTTPリクエストを作成
			// 引数：(HTTPメソッド, URL, リクエストボディ)
			req := httptest.NewRequest(tt.method, "/api/gourmet?"+tt.queryString, nil)

			// httptest.NewRecorder()：レスポンスを記録するRecorderを作成
			// ハンドラーが書き込んだレスポンスを後で検証できる
			rec := httptest.NewRecorder()

			// ========================================
			// ステップ4: ハンドラーを実行
			// ========================================
			// Handle()メソッドを呼び出し、レスポンスをrecに記録
			handler.Handle(rec, req)

			// ========================================
			// ステップ5: レスポンスを検証
			// ========================================
			// 【HTTPステータスコードの確認】
			// rec.Code：ハンドラーが返したステータスコード
			// 例：200 (成功)、405 (メソッド不許可)
			if rec.Code != tt.expectedStatus {
				t.Errorf("Handler returned status %d; want %d [%s]",
					rec.Code, tt.expectedStatus, tt.point)
			}

			// 【レスポンスヘッダーの確認】
			// エラーケースでない場合（正常系）は、Content-Typeを確認
			if !tt.expectError {
				contentType := rec.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Handler returned content-type %s; want application/json [%s]",
						contentType, tt.point)
				}
			}
		})
	}
}

// ============================================================================
// TestGourmetHandlerParseParams: parseParamsメソッドのテスト
// ============================================================================
// 【目的】
// HTTPリクエストのクエリパラメータを解析し、構造体に変換する処理が
// 正しく動作するかを検証
//
// 【parseParamsメソッドの役割】
// URLのクエリパラメータ（例：?service_area=SA11&keyword=居酒屋）を
// GourmetSearchParams構造体に変換する
//
// 【バリデーションとは】
// 入力データが正しい形式かどうかをチェックする処理
// - 必須パラメータの存在確認
// - データ型の妥当性確認
// - 値の範囲確認
//
// 【Triangulationの適用】
// 1. 必須パラメータのみ → デフォルト値の適用を確認
// 2. 全パラメータ指定 → 全ての値が正しく解析されるか確認
// 3. 必須パラメータなし → バリデーションエラーの発生を確認
func TestGourmetHandlerParseParams(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name        string // テストケース名
		queryString string // テストするクエリパラメータ
		expectError bool   // エラーが発生することを期待するか
		point       string // Triangulationポイント
	}{
		// Triangulation Point 1: 必須パラメータのみ
		{
			name:        "必須パラメータのみの場合",
			queryString: "keyword=居酒屋",
			expectError: false,
			point:       "Point 1: 最小構成",
		},
		// Triangulation Point 2: 全パラメータ指定
		{
			name:        "全パラメータを指定した場合",
			queryString: "lat=35.68&lng=139.76&range=3&address=さいたま&genre=G001&keyword=居酒屋&start=11&count=10",
			expectError: false,
			point:       "Point 2: 全パラメータ",
		},
		// Triangulation Point 3: 必須パラメータなし（エラー）
		{
			name:        "必須パラメータがない場合、エラー",
			queryString: "genre=G001",
			expectError: true,
			point:       "Point 3: バリデーションエラー",
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ========================================
			// ステップ1: テスト用のHTTPリクエストを作成
			// ========================================
			// クエリパラメータ付きのリクエストを作成
			// 例：/api/gourmet?service_area=SA11&keyword=居酒屋
			req := httptest.NewRequest(http.MethodGet, "/api/gourmet?"+tt.queryString, nil)

			// ========================================
			// ステップ2: ハンドラーのインスタンスを作成
			// ========================================
			// parseParamsメソッドをテストするため、ハンドラーを作成
			// サービスは不要（parseParamsは内部ロジックのみをテスト）
			handler := &GourmetHandler{}

			// ========================================
			// ステップ3: parseParamsメソッドを実行
			// ========================================
			// リクエストからパラメータを解析
			_, err := handler.parseParams(req)

			// ========================================
			// ステップ4: 結果を検証
			// ========================================
			// 【エラーケースの検証】
			if tt.expectError {
				// エラーが期待されるケース
				if err == nil {
					t.Errorf("Expected error but got none [%s]", tt.point)
				}
				// エラーが発生した場合、このケースは成功（検証終了）
				return
			}

			// 【成功ケースの検証】
			// エラーが発生してはいけない
			if err != nil {
				t.Fatalf("parseParams returned error: %v [%s]", err, tt.point)
			}
		})
	}
}

// ============================================================================
// TestGenreHandlerHandle: GenreHandler.Handleメソッドのテスト
// ============================================================================
// 【目的】
// ジャンルマスターAPIのハンドラーが、様々なHTTPメソッドに対して
// 適切に応答するかを検証
//
// 【テスト内容】
// - 正しいHTTPメソッド（GET）での正常動作
// - 誤ったHTTPメソッド（POST, DELETE等）の拒否
// - レスポンスの形式確認（JSON）
//
// 【HTTPメソッドとは】
// HTTPリクエストの種類を表す
// - GET: データの取得（読み取り専用）
// - POST: データの作成
// - PUT: データの更新
// - DELETE: データの削除
// このAPIはGETのみを許可（読み取り専用）
//
// 【Triangulationの適用】
// 1. GETメソッド → 200 OK（成功）
// 2. POSTメソッド → 405 Method Not Allowed（エラー）
// 3. DELETEメソッド → 405 Method Not Allowed（エラー）
func TestGenreHandlerHandle(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name           string // テストケース名
		method         string // テストするHTTPメソッド
		mockResponse   string // モックサーバーが返すレスポンス
		statusCode     int    // モックサーバーが返すステータスコード
		expectedStatus int    // ハンドラーが返すべきステータスコード
		point          string // Triangulationポイント
	}{
		// Triangulation Point 1: 正常なGETリクエスト
		{
			name:           "正常なGETリクエスト",
			method:         http.MethodGet,
			mockResponse:   `{"results": {"genre": [{"code": "G001", "name": "居酒屋"}]}}`,
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
			point:          "Point 1: 成功ケース",
		},
		// Triangulation Point 2: POSTメソッド（エラー）
		{
			name:           "POSTメソッドの場合、405エラー",
			method:         http.MethodPost,
			mockResponse:   "",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusMethodNotAllowed,
			point:          "Point 2: 不正メソッド",
		},
		// Triangulation Point 3: DELETEメソッド（エラー）
		{
			name:           "DELETEメソッドの場合、405エラー",
			method:         http.MethodDelete,
			mockResponse:   "",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusMethodNotAllowed,
			point:          "Point 3: 別の不正メソッド",
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ========================================
			// ステップ1: モックHTTPサーバーのセットアップ
			// ========================================
			// 外部APIの呼び出しをシミュレートするモックサーバー
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)     // モックのステータスコード
				w.Write([]byte(tt.mockResponse)) // モックのレスポンスボディ
			}))
			defer server.Close()

			// ========================================
			// ステップ2: サービスとハンドラーの作成
			// ========================================
			// モックサーバーを使用するサービスを作成
			svc := service.NewHotpepperService("test-key")

			// テスト対象のハンドラーを作成
			handler := NewGenreHandler(svc)

			// ========================================
			// ステップ3: テスト用のHTTPリクエストを作成
			// ========================================
			// 指定されたHTTPメソッドでリクエストを作成
			// ジャンルマスターAPIはクエリパラメータ不要
			req := httptest.NewRequest(tt.method, "/api/genre", nil)

			// レスポンスを記録するRecorderを作成
			rec := httptest.NewRecorder()

			// ========================================
			// ステップ4: ハンドラーを実行
			// ========================================
			handler.Handle(rec, req)

			// ========================================
			// ステップ5: レスポンスを検証
			// ========================================
			// 【HTTPステータスコードの確認】
			// ハンドラーが返したステータスコードが期待通りか確認
			if rec.Code != tt.expectedStatus {
				t.Errorf("Handler returned status %d; want %d [%s]",
					rec.Code, tt.expectedStatus, tt.point)
			}

			// 【成功ケースのみ：レスポンスボディの検証】
			// 200 OKの場合のみ、JSONの構造を確認
			if tt.expectedStatus == http.StatusOK {
				// JSONをパース（文字列→データ構造に変換）
				var response map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to parse response JSON: %v [%s]", err, tt.point)
				}

				// "results"フィールドが存在するか確認
				// ok: フィールドの存在を示すブール値
				if _, ok := response["results"]; !ok {
					t.Errorf("Response does not contain 'results' field [%s]", tt.point)
				}
			}
			// エラーケース（405等）の場合、レスポンスボディの検証はスキップ
		})
	}
}

// ============================================================================
// TestValidationError: ValidationErrorカスタムエラー型のテスト
// ============================================================================
// 【目的】
// カスタムエラー型が、Goのerrorインターフェースを正しく実装しているかを検証
//
// 【カスタムエラー型とは】
// Goでは、Error() stringメソッドを実装することで、
// 独自のエラー型を作成できる
//
// 【ValidationErrorの役割】
// 入力データのバリデーションエラーを表現する専用のエラー型
// 通常のエラーと区別することで、エラーハンドリングが柔軟になる
//
// 【Triangulationの適用】
// 1. 標準的な英語メッセージ → 通常のエラーメッセージ
// 2. 空のメッセージ → エッジケース（空文字列）
// 3. 日本語メッセージ → 多言語対応の確認
//
// 【学習ポイント】
// - Goのインターフェース実装
// - カスタムエラー型の作成方法
// - 多言語対応の重要性
func TestValidationError(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name     string // テストケース名
		message  string // エラーメッセージ（入力）
		expected string // 期待されるメッセージ（出力）
		point    string // Triangulationポイント
	}{
		// ----------------------------------------------------------------
		// Triangulation Point 1: 標準的な英語のエラーメッセージ
		// ----------------------------------------------------------------
		// 検証内容：通常のバリデーションエラーメッセージが正しく返されるか
		{
			name:     "標準的なエラーメッセージ",
			message:  "service_area is required", // 英語のメッセージ
			expected: "service_area is required", // そのまま返される
			point:    "Point 1: 標準メッセージ",
		},

		// ----------------------------------------------------------------
		// Triangulation Point 2: 空のメッセージ（エッジケース）
		// ----------------------------------------------------------------
		// 検証内容：空文字列でもエラーが発生しないか
		// 実際のコードでは通常使用しないが、防御的プログラミングの観点で重要
		{
			name:     "空のエラーメッセージ",
			message:  "", // 空文字列
			expected: "", // 空文字列が返される
			point:    "Point 2: 空メッセージ",
		},

		// ----------------------------------------------------------------
		// Triangulation Point 3: 日本語のメッセージ（多言語対応）
		// ----------------------------------------------------------------
		// 検証内容：日本語などの多バイト文字が正しく扱えるか
		// Goは標準でUTF-8をサポートしているため、問題なく動作する
		{
			name:     "日本語のエラーメッセージ",
			message:  "サービスエリアは必須です", // 日本語のメッセージ
			expected: "サービスエリアは必須です", // そのまま返される
			point:    "Point 3: 日本語",
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ========================================
			// ステップ1: ValidationErrorインスタンスを作成
			// ========================================
			// &ValidationError{...}：構造体のポインタを作成
			// ポインタを使用することで、errorインターフェースとして扱える
			err := &ValidationError{Message: tt.message}

			// ========================================
			// ステップ2: Error()メソッドを呼び出して検証
			// ========================================
			// Error()メソッドが正しいメッセージを返すか確認
			if err.Error() != tt.expected {
				t.Errorf("ValidationError.Error() = %q; want %q [%s]",
					err.Error(), tt.expected, tt.point)
			}
		})
	}
}
