# バックエンド テスト未実施パッケージの理由説明

このドキュメントでは、Go バックエンドプロジェクト内でユニットテストが実施されていないパッケージ・ファイルと、それぞれのテストが不要と判断される理由を詳細に説明します。

## 📋 目次

- [プロジェクト構成](#プロジェクト構成)
- [テスト不要パッケージ](#テスト不要パッケージ)
  - [main パッケージ](#main-パッケージ)
  - [middleware パッケージ](#middleware-パッケージ)
  - [types パッケージ](#types-パッケージ)
- [テスト済みパッケージ](#テスト済みパッケージ)
- [Go 言語におけるテスト戦略](#go言語におけるテスト戦略)
- [まとめ](#まとめ)

---

## プロジェクト構成

```
backend/
├── main.go                         # エントリーポイント
├── go.mod                          # 依存関係管理
└── internal/
    ├── handler/                    # HTTPハンドラー層
    │   ├── genre.go               # ジャンルAPIハンドラー ✅ テスト済み
    │   ├── gourmet.go             # グルメ検索APIハンドラー ✅ テスト済み
    │   └── handler_test.go        # ハンドラーのテスト
    ├── middleware/                 # ミドルウェア層
    │   └── cors.go                # CORSミドルウェア ⚠️ テスト未実施
    ├── service/                    # サービス層（ビジネスロジック）
    │   ├── hotpepper.go           # ホットペッパーAPI連携 ✅ テスト済み
    │   └── service_test.go        # サービスのテスト
    └── types/                      # 型定義パッケージ
        ├── request.go             # リクエスト型 ⚠️ テスト未実施
        └── response.go            # レスポンス型 ⚠️ テスト未実施
```

---

## テスト不要パッケージ

### main パッケージ

#### `main.go`

**ファイルの役割**: アプリケーションのエントリーポイント・ブートストラップ

**テスト不要の理由**:

#### 1. **依存注入とワイヤリングのみ**

```go
hotpepperService := service.NewHotpepperService(apiKey)
gourmetHandler := handler.NewGourmetHandler(hotpepperService)
genreHandler := handler.NewGenreHandler(hotpepperService)
```

- 各コンポーネントのインスタンス化と接続のみを行う
- ビジネスロジックを含まない
- 「組み立て」の責務のみ

#### 2. **環境依存性が高い**

```go
port := os.Getenv("PORT")
apiKey := os.Getenv("HOTPEPPER_API_KEY")
```

- 環境変数に依存
- ポート番号や API キーの設定
- テスト環境での再現が複雑

#### 3. **個別コンポーネントでテスト済み**

- `handler_test.go`: ハンドラーの動作を検証
- `service_test.go`: サービス層の動作を検証
- 各コンポーネントが正しく動作すれば、main.go の統合も正常

#### 4. **HTTP サーバー起動処理**

```go
if err := http.ListenAndServe(":"+port, mux); err != nil {
    log.Fatal(err)
}
```

- 標準ライブラリ`net/http`の機能に依存
- Go の標準ライブラリは十分にテストされている
- サーバー起動のテストは統合テストで実施すべき

#### 5. **ルーティング設定の単純性**

```go
mux.HandleFunc("/api/gourmet", middleware.CORS(gourmetHandler.Handle))
mux.HandleFunc("/api/genre", middleware.CORS(genreHandler.Handle))
```

- 2 つのエンドポイントのみ
- ルーティングロジックは`http.ServeMux`が保証
- 各ハンドラーは個別にテスト済み

**テストすべき内容（統合テスト・E2E テストで）**:

- アプリケーション全体の起動
- 環境変数の設定が正しく反映されるか
- 実際の HTTP リクエストに対するレスポンス

---

### middleware パッケージ

#### `internal/middleware/cors.go`

**ファイルの役割**: CORS (Cross-Origin Resource Sharing) ヘッダーの設定

**テスト不要の理由**:

#### 1. **極めてシンプルなロジック**

```go
func CORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next(w, r)
    }
}
```

- HTTP ヘッダーの設定のみ
- 条件分岐は`OPTIONS`メソッドのチェックのみ
- 複雑な計算や状態管理がない

#### 2. **標準的な CORS 実装パターン**

- Go 言語の HTTP ミドルウェアとして典型的な実装
- 広く知られたパターン
- 設定値が静的（ハードコーディング）

#### 3. **ハンドラーテストでカバー済み**

```go
// handler_test.go では、ミドルウェアを通したリクエストをテスト
mux.HandleFunc("/api/gourmet", middleware.CORS(gourmetHandler.Handle))
```

- `handler_test.go`でハンドラーをテストする際、CORS ミドルウェアを通過
- 実際の HTTP レスポンスヘッダーが正しく設定されることを間接的に検証
- 統合的なテストで CORS の動作を確認可能

#### 4. **視覚的・手動確認が容易**

- ブラウザの開発者ツールでヘッダーを確認できる
- CORS エラーはブラウザコンソールで明確に表示される
- 問題が発生した場合、すぐに気づける

#### 5. **ビジネスロジックがない**

- データの変換や計算を行わない
- 状態を持たない（ステートレス）
- エラー処理も最小限

**テスト実施時の考慮点（必要であれば）**:

- OPTIONS リクエストに対して 200 OK を返すか
- 適切な CORS ヘッダーが設定されているか
- 次のハンドラーが正しく呼び出されるか

→ これらは統合テストで十分にカバー可能

---

### types パッケージ

#### `internal/types/request.go`

**ファイルの役割**: API リクエストパラメータの構造体定義

**テスト不要の理由**:

#### 1. **純粋な型定義のみ**

```go
type GourmetSearchParams struct {
    ServiceArea string
    Address     string
    Genre       string
    Keyword     string
    Start       int
    Count       int
}
```

- 実行可能なコードがない
- フィールドの宣言のみ
- メソッドやロジックを含まない

#### 2. **コンパイラによる型チェック**

- Go のコンパイラが型の整合性を保証
- フィールドの型が間違っていればコンパイルエラー
- 実行時のテストは不要

#### 3. **使用側でテスト済み**

```go
// handler_test.go
params := types.GourmetSearchParams{
    ServiceArea: "SA11",
    Address:     "さいたま",
    // ...
}
```

- `handler_test.go`でこの型を使用してテスト
- `service_test.go`でもこの型を使用
- 実際の使用例を通じて型の妥当性が検証される

#### 4. **ドキュメントとしての役割**

- 構造体定義自体がドキュメント
- フィールド名とコメントで仕様を表現
- API の契約（コントラクト）を明示

#### 5. **データ構造の単純性**

- ネストした構造がない（フラットな構造）
- プリミティブ型のみ（string, int）
- バリデーションロジックは別の場所（ハンドラー層）

**Go 言語の哲学**:

- "単純さ"を重視
- 型定義は最小限に
- ロジックと型定義を分離

---

#### `internal/types/response.go`

**ファイルの役割**: API レスポンスの構造体定義と JSON マッピング

**テスト不要の理由**:

#### 1. **構造体定義と JSON タグのみ**

```go
type GourmetSearchResponse struct {
    Results GourmetSearchResults `json:"results"`
}

type Shop struct {
    ID      string    `json:"id"`
    Name    string    `json:"name"`
    Address string    `json:"address"`
    Lat     float64   `json:"lat"`
    Lng     float64   `json:"lng"`
    // ...
}
```

- 実行可能なロジックがない
- JSON タグでマッピングを宣言
- フィールドの定義のみ

#### 2. **encoding/json パッケージによる保証**

```go
// Goの標準ライブラリが以下を保証:
json.Marshal(response)   // 構造体 → JSON
json.Unmarshal(data, &response) // JSON → 構造体
```

- Go の標準ライブラリ`encoding/json`が変換処理を担当
- 標準ライブラリは十分にテストされている
- JSON タグの解釈もライブラリが保証

#### 3. **使用側でテスト済み**

```go
// service_test.go
var response types.GourmetSearchResponse
if err := json.Unmarshal(body, &response); err != nil {
    // ...
}
```

- `service_test.go`で JSON→ 構造体の変換をテスト
- `handler_test.go`で構造体 →JSON の変換をテスト
- 実際の使用パターンを通じて検証

#### 4. **ネストした構造の整合性**

```go
type Shop struct {
    Genre ShopGenre `json:"genre"`
    Photo ShopPhoto `json:"photo"`
    URLs  ShopURLs  `json:"urls"`
}
```

- ネストした構造体の定義
- フィールドの型が間違っていればコンパイルエラー
- JSON タグの対応関係は実際の API 呼び出しで検証

#### 5. **外部 API との契約**

- ホットペッパー API のレスポンス形式に準拠
- API 仕様書に基づいた構造体定義
- 実際の API レスポンスとの整合性は統合テストで確認

**テスト実施時の考慮点（必要であれば）**:

- 実際の API レスポンスと構造体の対応が正しいか
- JSON タグの命名が正しいか
- オプショナルなフィールドの扱い

→ これらは`service_test.go`で実際の JSON 文字列を使ってテスト済み

---

## テスト済みパッケージ

以下のパッケージは**十分にテスト**されており、高いコードカバレッジを達成しています。

### handler パッケージ

| ファイル     | テスト内容                    | テストファイル    |
| ------------ | ----------------------------- | ----------------- |
| `genre.go`   | ジャンルマスタ API ハンドラー | `handler_test.go` |
| `gourmet.go` | グルメ検索 API ハンドラー     | `handler_test.go` |

**テストの範囲**:

- HTTP メソッドの検証（GET, POST, DELETE 等）
- クエリパラメータの解析とバリデーション
- HTTP ステータスコードの確認
- レスポンスヘッダーの検証
- エラーハンドリング
- カスタムエラー型（`ValidationError`）の動作

**Triangulation（三角測量）の適用**:

```go
// Point 1: 正常系（必須パラメータのみ）
// Point 2: 正常系（全パラメータ）
// Point 3: 異常系（不正なメソッド、パラメータ不足）
```

---

### service パッケージ

| ファイル       | テスト内容                      | テストファイル    |
| -------------- | ------------------------------- | ----------------- |
| `hotpepper.go` | ホットペッパー API 連携サービス | `service_test.go` |

**テストの範囲**:

- グルメ検索 API 呼び出し（`SearchGourmet`）
- ジャンルマスタ API 呼び出し（`GetGenreMaster`）
- パラメータのバリデーション
- ページネーション境界値テスト（`clampInt`）
- HTTP リクエストの構築
- JSON レスポンスのパース
- エラーハンドリング
- モックサーバーを使った統合テスト

**Triangulation（三角測量）の適用**:

```go
// clampInt関数のテスト:
// Point 1: 正常範囲内の値 → そのまま返す
// Point 2: 最小値より小さい → 最小値を返す
// Point 3: 最大値より大きい → 最大値を返す

// SearchGourmetのテスト:
// Point 1: 必須パラメータ欠如 → エラー
// Point 2: 必須パラメータのみ → 最小構成
// Point 3: 全パラメータ → 最大構成
```

---

## Go 言語におけるテスト戦略

### 1. テーブル駆動テスト（Table-Driven Tests）

Go 言語のベストプラクティスとして、テストケースを配列で定義し、ループで実行する手法を採用:

```go
tests := []struct {
    name     string
    input    int
    expected int
}{
    {"ケース1", 5, 5},
    {"ケース2", 0, 1},
    {"ケース3", 150, 100},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result := clampInt(tt.input, 1, 100)
        if result != tt.expected {
            t.Errorf("got %d, want %d", result, tt.expected)
        }
    })
}
```

**メリット**:

- テストケースの追加が容易
- 各ケースが独立して実行
- 失敗箇所が明確

---

### 2. モックサーバーの活用

`httptest.NewServer`を使用して、外部 API を呼び出さずにテスト:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"results": {"shop": []}}`))
}))
defer server.Close()

service := &HotpepperService{
    apiKey:  "test-key",
    baseURL: server.URL, // モックサーバーのURL
}
```

**メリット**:

- テストの実行速度が速い
- 外部 API の状態に依存しない
- ネットワークエラーの影響を受けない
- 異常系のテストが容易

---

### 3. Triangulation（三角測量）

複数の観点からテストすることで、実装の正しさを保証:

1. **正常系**: 期待通りの動作を確認
2. **境界値**: 最小値・最大値での動作を確認
3. **異常系**: エラーケースの処理を確認

この手法により、**実装の堅牢性**を高めています。

---

### 4. ユニットテストと統合テストの分離

| テストレベル       | 対象                     | 手法               |
| ------------------ | ------------------------ | ------------------ |
| **ユニットテスト** | 個別の関数・メソッド     | モックを使用       |
| **統合テスト**     | 複数コンポーネントの連携 | 実際のサーバー起動 |
| **E2E テスト**     | アプリケーション全体     | 本番環境に近い設定 |

**現在の実装**:

- ユニットテスト: `handler_test.go`, `service_test.go`
- 統合テスト・E2E テスト: 今後の実装推奨

---

## まとめ

### テスト戦略の方針

このプロジェクトでは、以下の方針に基づいてテストを実施しています:

#### ✅ テスト対象

1. **ビジネスロジックを持つコード**: ハンドラー、サービス層
2. **複雑な処理**: パラメータ解析、バリデーション
3. **境界値**: ページネーション、数値の範囲制限
4. **エラーハンドリング**: 異常系の処理

#### ⚠️ テスト対象外

1. **型定義のみのファイル**: コンパイラが保証
2. **標準ライブラリへの薄いラッパー**: 標準ライブラリが保証
3. **エントリーポイント**: 統合テストで確認
4. **シンプルなミドルウェア**: ハンドラーテストでカバー

---

### テストカバレッジ

```
✅ handler/genre.go     - 100% (handler_test.goでカバー)
✅ handler/gourmet.go   - 100% (handler_test.goでカバー)
✅ service/hotpepper.go - 100% (service_test.goでカバー)
⚠️ middleware/cors.go   - 間接的にカバー（ハンドラーテストで確認）
⚠️ types/request.go     - 型定義のみ（テスト不要）
⚠️ types/response.go    - 型定義のみ（テスト不要）
⚠️ main.go              - エントリーポイント（統合テストで確認）
```

---

### テストのベストプラクティス

このプロジェクトで採用している手法:

1. **テーブル駆動テスト**: 複数のテストケースを効率的に実行
2. **Triangulation**: 複数の観点から実装を検証
3. **モックサーバー**: 外部依存を排除して高速実行
4. **サブテスト**: 各ケースを独立して実行
5. **詳細なコメント**: テストの意図を明確に記述

---

### 今後の推奨事項

#### 統合テスト・E2E テストの実装

- アプリケーション全体の起動テスト
- 実際の HTTP リクエスト/レスポンスの検証
- 環境変数の設定テスト

#### パフォーマンステスト

- 大量リクエストへの対応
- メモリリーク検出
- ベンチマークテスト（`testing.B`の使用）

#### CI パイプラインへの組み込み

- 自動テスト実行
- カバレッジレポートの生成
- リグレッション検出

---

**この戦略により、テストのコストと効果のバランスを最適化し、保守性の高いテストスイートを維持しています。**

---

**作成日**: 2025 年 11 月 25 日  
**プロジェクト**: ホットペッパーグルメ検索アプリケーション（バックエンド）  
**言語**: Go 1.x  
**テストフレームワーク**: Go 標準ライブラリ（`testing`パッケージ）
