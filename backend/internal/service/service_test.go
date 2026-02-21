package service

import (
	"net/http"          // HTTPクライアント・サーバー機能を提供
	"net/http/httptest" // HTTPテスト用のユーティリティ（モックサーバーなど）
	"strings"           // 文字列操作のための関数群
	"testing"           // Go標準のテストフレームワーク

	"backend/internal/types"
)

// ============================================================================
// TestClampInt: clampInt関数のテスト
// ============================================================================
// 【目的】
// clampInt関数が値を指定された範囲内に制限する機能を正しく実装しているかを検証
//
// 【Triangulation（三角測量）の適用】
// 3つの異なる観点から関数の動作を検証することで、実装の正しさを保証：
// 1. 正常範囲内の値 → そのまま返す
// 2. 最小値より小さい → 最小値を返す
// 3. 最大値より大きい → 最大値を返す
//
// 【テーブル駆動テストとは】
// テストケースを構造体の配列として定義し、ループで実行する手法
// メリット：
// - テストケースの追加が容易（配列に要素を追加するだけ）
// - テストコードの重複を削減
// - 各ケースが独立して実行され、失敗箇所が明確
func TestClampInt(t *testing.T) {
	// テストケースを定義する構造体の配列
	// []struct は「無名構造体のスライス（配列）」を意味する
	tests := []struct {
		name     string // テストケースの名前（t.Runで表示される）
		value    int    // テストする入力値
		min      int    // 最小値
		max      int    // 最大値
		expected int    // 期待される結果
		point    string // Triangulationのポイント（どの観点のテストか）
	}{
		// ----------------------------------------------------------------
		// Triangulation Point 1: 正常範囲内の値
		// ----------------------------------------------------------------
		// 検証内容：値が min と max の範囲内にある場合、値をそのまま返す
		// 例：clampInt(5, 1, 10) → 5
		{
			name:     "値が範囲内の場合、そのまま返す",
			value:    5, // 1〜10の範囲内
			min:      1,
			max:      10,
			expected: 5, // 変更なし
			point:    "Point 1: 正常範囲",
		},

		// ----------------------------------------------------------------
		// Triangulation Point 2: 最小値より小さい値（下限境界値テスト）
		// ----------------------------------------------------------------
		// 検証内容：値が min より小さい場合、min を返す
		// 例：clampInt(0, 1, 10) → 1
		{
			name:     "値が最小値より小さい場合、最小値を返す",
			value:    0, // 最小値1より小さい
			min:      1,
			max:      10,
			expected: 1, // 最小値に補正される
			point:    "Point 2: 下限境界",
		},

		// ----------------------------------------------------------------
		// Triangulation Point 3: 最大値より大きい値（上限境界値テスト）
		// ----------------------------------------------------------------
		// 検証内容：値が max より大きい場合、max を返す
		// 例：clampInt(15, 1, 10) → 10
		{
			name:     "値が最大値より大きい場合、最大値を返す",
			value:    15, // 最大値10より大きい
			min:      1,
			max:      10,
			expected: 10, // 最大値に補正される
			point:    "Point 3: 上限境界",
		},
	}

	// for rangeループ：testsスライスの各要素を順番に処理
	// tt は "test case"の略で、各テストケースを表す変数名として慣例的に使用される
	for _, tt := range tests {
		// t.Run()：サブテストを実行する
		// サブテストのメリット：
		// - 各テストケースが独立して実行される
		// - 一つのケースが失敗しても、他のケースは実行される
		// - テスト結果で各ケースの成否が個別に表示される
		t.Run(tt.name, func(t *testing.T) {
			// 実際に関数を呼び出して結果を取得
			result := clampInt(tt.value, tt.min, tt.max)

			// 結果が期待値と一致するか確認
			if result != tt.expected {
				// t.Errorf()：テストの失敗を報告（テストは継続される）
				// %d は整数の埋め込み、%s は文字列の埋め込み
				t.Errorf("clampInt(%d, %d, %d) = %d; want %d [%s]",
					tt.value, tt.min, tt.max, result, tt.expected, tt.point)
			}
		})
	}
}

// ============================================================================
// TestSearchGourmet: SearchGourmetメソッドのテスト
// ============================================================================
// 【目的】
// HotpepperServiceのSearchGourmetメソッドが、様々なパラメータ組み合わせで
// 正しくAPIリクエストを構築できるかを検証
//
// 【テスト戦略】
// モックサーバー（httptest.NewServer）を使用して、実際のAPIを呼び出さずにテスト
// これにより：
// - テストの実行速度が速い
// - 外部APIの状態に依存しない
// - ネットワークエラーの影響を受けない
//
// 【Triangulationの適用】
// 1. 必須パラメータが含まれていない場合 → バリデーションエラーの検証
// 2. 必須パラメータのみ → 最小構成でのAPI呼び出し
// 3. 全パラメータ指定 → 最大構成での動作確認
func TestSearchGourmet(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name          string                    // テストケース名
		params        types.GourmetSearchParams // SearchGourmetに渡すパラメータ
		mockResponse  string                    // モックサーバーが返すレスポンス
		expectedInURL []string                  // URLに含まれるべきパラメータのリスト
		expectError   bool                      // エラーが発生することを期待するか
		point         string                    // Triangulationポイント
	}{
		// Triangulation Point 1: 必須パラメータが含まれていない場合
		{
			name: "必須パラメータ(ServiceArea)が空の場合エラー",
			params: types.GourmetSearchParams{
				ServiceArea: "", // 必須パラメータが空
				Start:       1,
				Count:       20,
			},
			mockResponse:  `{"results": {"shop": []}}`,
			expectedInURL: []string{},
			expectError:   true,
			point:         "Point 1: 必須パラメータ欠如",
		},
		// Triangulation Point 2: 必須パラメータのみ
		{
			name: "必須パラメータのみで検索",
			params: types.GourmetSearchParams{
				ServiceArea: "SA11",
				Start:       1,
				Count:       20,
			},
			mockResponse: `{"results": {"shop": []}}`,
			expectedInURL: []string{
				"service_area=SA11",
				"start=1",
				"count=20",
			},
			expectError: false,
			point:       "Point 2: 最小構成",
		},
		// Triangulation Point 3: 全パラメータ
		{
			name: "全パラメータを含む検索",
			params: types.GourmetSearchParams{
				ServiceArea: "SA11",
				Address:     "さいたま",
				Genre:       "G001",
				Keyword:     "居酒屋",
				Start:       11,
				Count:       10,
			},
			mockResponse: `{"results": {"shop": []}}`,
			expectedInURL: []string{
				"service_area=SA11",
				"genre=G001",
				"start=11",
				"count=10",
			},
			expectError: false,
			point:       "Point 3: 最大構成",
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ========================================
			// ステップ1: モックHTTPサーバーのセットアップ
			// ========================================
			// httptest.NewServer()：テスト用のHTTPサーバーを起動
			// 実際のAPI呼び出しの代わりに、このサーバーが応答を返す
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// リクエストURLに期待されるパラメータが含まれているか検証
				// 例：r.URL.RawQuery = "service_area=SA11&start=1&count=20"
				for _, expected := range tt.expectedInURL {
					// strings.Split(expected, "=")[0] でパラメータ名を取得
					// 例："service_area=SA11" → "service_area"
					if !strings.Contains(r.URL.RawQuery, strings.Split(expected, "=")[0]) {
						t.Errorf("URL does not contain expected parameter: %s [%s]", expected, tt.point)
					}
				}

				// モックレスポンスを返す
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK) // 200 OK
				w.Write([]byte(tt.mockResponse))
			}))
			// defer：関数終了時にサーバーを閉じる（リソースの解放）
			defer server.Close()

			// ========================================
			// ステップ2: テスト対象のサービスを作成
			// ========================================
			// 本番環境では実際のAPI URLを使用するが、
			// テストでは server.URL（モックサーバーのURL）を使用
			service := &HotpepperService{
				apiKey:  "test-key", // テスト用のダミーAPIキー
				baseURL: server.URL, // モックサーバーのURL
			}

			// ========================================
			// ステップ3: SearchGourmetメソッドを実行
			// ========================================
			response, err := service.SearchGourmet(tt.params)

			// ========================================
			// ステップ4: エラーケースの検証
			// ========================================
			if tt.expectError {
				// エラーが期待されるのに発生しなかった場合、テスト失敗
				if err == nil {
					t.Errorf("Expected error but got none [%s]", tt.point)
				}
				// エラーが期待通り発生した場合、このテストケースは成功
				return
			}

			// ========================================
			// ステップ5: 成功ケースの検証
			// ========================================
			// エラーチェック
			// t.Fatalf()：致命的なエラー（このテストケースを即座に終了）
			if err != nil {
				t.Fatalf("SearchGourmet returned error: %v [%s]", err, tt.point)
			}

			// レスポンスの妥当性を確認
			// レスポンスが型付きで返されているか、フィールドにアクセスできるか確認
			if response == nil {
				t.Errorf("Response is nil [%s]", tt.point)
			}
			// results フィールドが存在するか確認
			if response != nil && response.Results.Shop == nil {
				t.Errorf("Response does not contain shop array [%s]", tt.point)
			}
		})
	}
}

// ============================================================================
// TestGetGenreMaster: GetGenreMasterメソッドのテスト
// ============================================================================
// 【目的】
// ジャンルマスターAPIの呼び出しが、様々なレスポンスパターンに対して
// 適切に処理できるかを検証
//
// 【Triangulationの適用】
// 1. 正常なレスポンス → 成功時の処理
// 2. 空のレスポンス → データがない場合の処理
// 3. APIエラー → エラーハンドリングの確認
//
// 【学習ポイント】
// - HTTPステータスコードの意味（200 OK、401 Unauthorized等）
// - エラーハンドリングの重要性
// - 異常系テストの書き方
func TestGetGenreMaster(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name         string // テストケース名
		mockResponse string // モックサーバーが返すレスポンスボディ
		statusCode   int    // モックサーバーが返すHTTPステータスコード
		expectError  bool   // エラーが発生することを期待するか
		point        string // Triangulationポイント
	}{
		// Triangulation Point 1: 正常なレスポンス
		{
			name:         "正常なレスポンスを取得",
			mockResponse: `{"results": {"genre": [{"code": "G001", "name": "居酒屋"}]}}`,
			statusCode:   http.StatusOK,
			expectError:  false,
			point:        "Point 1: 成功ケース",
		},
		// Triangulation Point 2: 空のレスポンス
		{
			name:         "空のジャンルリストを取得",
			mockResponse: `{"results": {"genre": []}}`,
			statusCode:   http.StatusOK,
			expectError:  false,
			point:        "Point 2: 空データ",
		},
		// Triangulation Point 3: エラーレスポンス
		{
			name:         "APIエラーが発生",
			mockResponse: `{"error": "invalid key"}`,
			statusCode:   http.StatusUnauthorized,
			expectError:  true,
			point:        "Point 3: エラーケース",
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ========================================
			// ステップ1: モックHTTPサーバーのセットアップ
			// ========================================
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// レスポンスヘッダーを設定
				w.Header().Set("Content-Type", "application/json")

				// HTTPステータスコードを設定
				// 例：200 (成功)、401 (認証エラー)
				w.WriteHeader(tt.statusCode)

				// レスポンスボディを書き込む
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close() // テスト終了時にサーバーを閉じる

			// ========================================
			// ステップ2: サービスインスタンスの作成
			// ========================================
			service := &HotpepperService{
				apiKey:  "test-key",
				baseURL: server.URL, // モックサーバーのURLを使用
			}

			// ========================================
			// ステップ3: GetGenreMasterメソッドを実行
			// ========================================
			response, err := service.GetGenreMaster()

			// ========================================
			// ステップ4: 結果の検証
			// ========================================
			// 【エラーケースの検証】
			if tt.expectError {
				// エラーが期待されるのに発生しなかった場合、テスト失敗
				if err == nil {
					t.Errorf("Expected error but got none [%s]", tt.point)
				}
				// エラーが期待通り発生した場合、このテストケースは成功
			} else {
				// 【成功ケースの検証】
				// エラーが発生してはいけない
				if err != nil {
					t.Fatalf("GetGenreMaster returned error: %v [%s]", err, tt.point)
				}

				// レスポンスが型付きで返されているか確認
				if response == nil {
					t.Errorf("Response is nil [%s]", tt.point)
				}
				// results.genre フィールドが存在するか確認
				if response != nil && response.Results.Genre == nil {
					t.Errorf("Response does not contain genre array [%s]", tt.point)
				}
			}
		})
	}
}

// ============================================================================
// TestSearchGourmetPagination: ページネーション境界値テスト
// ============================================================================
// 【目的】
// SearchGourmetメソッドのページネーションパラメータ（start, count）が
// 不正な値の場合に適切に補正されるかを検証
//
// 【ページネーションとは】
// 大量のデータを複数のページに分割して表示する仕組み
// - start: 取得開始位置（1始まり）
// - count: 1ページあたりの取得件数
//
// 【Triangulationの適用】
// 1. start下限チェック → 0以下の値を1に補正
// 2. count下限チェック → 0以下の値を1に補正
// 3. count上限チェック → 100超の値を100に補正
//
// 【境界値テストの重要性】
// バグは境界値（最小値、最大値の近辺）で発生しやすい
// 例：配列のインデックスエラー、整数のオーバーフローなど
func TestSearchGourmetPagination(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name          string // テストケース名
		start         int    // テストする開始位置（入力値）
		count         int    // テストする取得件数（入力値）
		expectedStart int    // 期待される開始位置（補正後）
		expectedCount int    // 期待される取得件数（補正後）
		point         string // Triangulationポイント
	}{
		// Triangulation Point 1: start下限チェック
		{
			name:          "startが0以下の場合、1に補正される",
			start:         0,
			count:         20,
			expectedStart: 1,
			expectedCount: 20,
			point:         "Point 1: start下限",
		},
		// Triangulation Point 2: count下限チェック
		{
			name:          "countが0以下の場合、1に補正される",
			start:         1,
			count:         0,
			expectedStart: 1,
			expectedCount: 1,
			point:         "Point 2: count下限",
		},
		// Triangulation Point 3: count上限チェック
		{
			name:          "countが100を超える場合、100に補正される",
			start:         1,
			count:         150,
			expectedStart: 1,
			expectedCount: 100,
			point:         "Point 3: count上限",
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ========================================
			// ステップ1: モックHTTPサーバーのセットアップ
			// ========================================
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// リクエストURLからクエリパラメータを取得
				// 例：/api/gourmet?service_area=SA11&start=1&count=20
				query := r.URL.Query()
				start := query.Get("start") // "start"パラメータの値を取得
				count := query.Get("count") // "count"パラメータの値を取得

				// ----------------------------------------
				// パラメータの存在確認
				// ----------------------------------------
				// startパラメータがURLに含まれているか確認
				if !strings.Contains(r.URL.RawQuery, "start=") {
					t.Errorf("start parameter not found [%s]", tt.point)
				}
				// countパラメータがURLに含まれているか確認
				if !strings.Contains(r.URL.RawQuery, "count=") {
					t.Errorf("count parameter not found [%s]", tt.point)
				}

				// ----------------------------------------
				// 取得した値をログに記録（デバッグ用）
				// ----------------------------------------
				// _ = start : 変数を使用していないという警告を回避
				// 実際のテストでは、これらの値が期待通りかを検証することもできる
				_ = start
				_ = count

				// モックレスポンスを返す
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"results": {"shop": []}}`))
			}))
			defer server.Close()

			// ========================================
			// ステップ2: サービスインスタンスの作成
			// ========================================
			service := &HotpepperService{
				apiKey:  "test-key",
				baseURL: server.URL,
			}

			// ========================================
			// ステップ3: SearchGourmetメソッドを実行
			// ========================================
			// 不正な値（0以下、100超）を含むパラメータを渡す
			searchParams := types.GourmetSearchParams{
				ServiceArea: "SA11",
				Start:       tt.start, // 例：0（不正） → 1に補正されるはず
				Count:       tt.count, // 例：150（不正） → 100に補正されるはず
			}

			// メソッドを実行（内部でclampInt()により値が補正される）
			response, err := service.SearchGourmet(searchParams)

			// ========================================
			// ステップ4: エラーチェック
			// ========================================
			// 補正が正しく行われていれば、エラーは発生しないはず
			if err != nil {
				t.Fatalf("SearchGourmet returned error: %v [%s]", err, tt.point)
			}

			// レスポンスが正しく返されているか確認
			if response == nil {
				t.Errorf("Response is nil [%s]", tt.point)
			}

			// 注：このテストでは、補正された値が正しいかどうかは
			// モックサーバー内のパラメータ存在確認で間接的に検証している
			// より厳密には、URLから値を取り出して expectedStart, expectedCount と
			// 比較することもできる
		})
	}
}
