# Triangulation（三角測量）の観点：共通箇所と差分

## Triangulation とは

**Triangulation（三角測量）** は、TDD（テスト駆動開発）において、**複数の異なる具体的なテストケースから一般的な実装を導出する手法**です。

最小限の 3 つのテストポイントで、以下を検証します：

- **Point 1**: 基本的な正常ケース
- **Point 2**: 別の角度からの検証（境界値や異なる入力）
- **Point 3**: さらに異なる観点（エッジケースやエラー）

この 3 点により、実装の一般性と正確性を保証します。

---

## 🔄 共通箇所（全テストで共有される構造）

### 1. テスト構造

```go
// すべてのテストで共通のパターン
func TestXxx(t *testing.T) {
    // テーブル駆動テスト：複数のケースを配列で定義
    tests := []struct {
        name     string  // テストケース名
        // ... 入力パラメータ
        expected xxx     // 期待される結果
        point    string  // Triangulationのポイント識別
    }{
        // Point 1, 2, 3の3つのケースを定義
    }

    // 各ケースをループで実行
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実行
            result := targetFunction(tt.input)

            // 検証
            if result != tt.expected {
                t.Errorf("...")
            }
        })
    }
}
```

### 2. モック手法

```go
// HTTPサーバーのモック（サービス層テストで共通）
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // モックレスポンスの作成
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"results": {...}}`))
}))
defer server.Close()
```

### 3. エラーハンドリング

```go
// すべてのテストで統一されたエラーチェック
result, err := someFunction(params)
if err != nil {
    t.Fatalf("Function returned error: %v", err)
}
```

### 4. 検証パターン

```go
// 戻り値の検証
if result != expected {
    t.Errorf("got %v; want %v", result, expected)
}

// HTTPステータスの検証
if rec.Code != expectedStatus {
    t.Errorf("status = %d; want %d", rec.Code, expectedStatus)
}

// JSON構造の検証
if !strings.Contains(string(body), "expected_field") {
    t.Errorf("Response missing expected field")
}
```

---

## 📊 差分（各テストケース固有の検証ポイント）

### 1. Helper 関数のテスト

#### `clampInt` 関数

| Point   | 入力値                    | 期待値 | 検証観点   | 差分の本質                           |
| ------- | ------------------------- | ------ | ---------- | ------------------------------------ |
| Point 1 | `value=5, min=1, max=10`  | `5`    | 正常範囲内 | **条件**: `min <= value <= max`      |
| Point 2 | `value=0, min=1, max=10`  | `1`    | 下限境界   | **条件**: `value < min` → 最小値返却 |
| Point 3 | `value=15, min=1, max=10` | `10`   | 上限境界   | **条件**: `value > max` → 最大値返却 |

**共通実装の導出**:

```go
func clampInt(value, min, max int) int {
    if value < min { return min }  // Point 2から導出
    if value > max { return max }  // Point 3から導出
    return value                    // Point 1から導出
}
```

#### `parseIntOrDefault` 関数

| Point   | 入力値          | 期待値          | 検証観点   | 差分の本質                      |
| ------- | --------------- | --------------- | ---------- | ------------------------------- |
| Point 1 | `""` (空文字列) | `10` (fallback) | 空値処理   | **条件**: `value == ""`         |
| Point 2 | `"42"`          | `42`            | 正常変換   | **条件**: 変換成功 `err == nil` |
| Point 3 | `"abc"`         | `10` (fallback) | 不正値処理 | **条件**: 変換失敗 `err != nil` |

**共通実装の導出**:

```go
func parseIntOrDefault(value string, fallback int) int {
    if value == "" { return fallback }        // Point 1から導出
    if v, err := strconv.Atoi(value); err == nil {
        return v                               // Point 2から導出
    }
    return fallback                            // Point 3から導出
}
```

---

### 2. サービス層のテスト

#### `SearchGourmet` メソッド

| Point   | パラメータ           | 検証する URL                         | エラー | 差分の本質                                         |
| ------- | -------------------- | ------------------------------------ | ------ | -------------------------------------------------- |
| Point 1 | `ServiceArea=""`     | -                                    | あり   | **バリデーション**: 必須パラメータが含まれていない |
| Point 2 | `ServiceArea`のみ    | `service_area=SA11&start=1&count=20` | なし   | **最小構成**: 必須パラメータのみ                   |
| Point 3 | + `Genre`, `Keyword` | + `genre=G001&keyword=...`           | なし   | **最大構成**: 全パラメータ                         |

**共通実装の導出**:

```go
// 必須パラメータのバリデーション（Point 1から導出）
if params.ServiceArea == "" {
    return nil, fmt.Errorf("service_area is required")
}

// 必須パラメータは常に追加（Point 2, 3で共通）
queryParams.Set("service_area", params.ServiceArea)

// オプションパラメータは条件付き追加（Point 3で有効）
if params.Address != "" {
    queryParams.Set("address", params.Address)
}
if params.Genre != "" {
    queryParams.Set("genre", params.Genre)
}
if params.Keyword != "" {
    queryParams.Set("keyword", params.Keyword)
}
```

#### `GetGenreMaster` メソッド

| Point   | HTTP ステータス    | レスポンス                      | エラー | 差分の本質                 |
| ------- | ------------------ | ------------------------------- | ------ | -------------------------- |
| Point 1 | `200 OK`           | `{"results": {"genre": [...]}}` | なし   | **成功ケース**: 正常データ |
| Point 2 | `200 OK`           | `{"results": {"genre": []}}`    | なし   | **空データ**: エッジケース |
| Point 3 | `401 Unauthorized` | `{"error": "..."}`              | あり   | **エラーケース**: 異常系   |

**共通実装の導出**:

```go
if resp.StatusCode != http.StatusOK {  // Point 3から導出
    return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
}
// Point 1, 2は正常レスポンスとして処理
return body, nil
```

#### ページネーション境界値

| Point   | 入力        | 期待値      | 差分の本質                  |
| ------- | ----------- | ----------- | --------------------------- |
| Point 1 | `start=0`   | `start=1`   | **下限補正**: `start < 1`   |
| Point 2 | `count=0`   | `count=1`   | **下限補正**: `count < 1`   |
| Point 3 | `count=150` | `count=100` | **上限補正**: `count > 100` |

**共通実装の導出**:

```go
// clampInt関数による統一的な境界値制御
start := clampInt(params.Start, 1, 1000)  // 全Pointでclampを適用
count := clampInt(params.Count, 1, 100)
```

---

### 3. ハンドラー層のテスト

#### `GourmetHandler.Handle` メソッド

| Point   | HTTP メソッド | クエリ              | 期待ステータス           | 差分の本質                       |
| ------- | ------------- | ------------------- | ------------------------ | -------------------------------- |
| Point 1 | `GET`         | `service_area=SA11` | `200 OK`                 | **正常ケース（最小）**: 基本機能 |
| Point 2 | `GET`         | 全パラメータ        | `200 OK`                 | **正常ケース（最大）**: 拡張機能 |
| Point 3 | `POST`        | -                   | `405 Method Not Allowed` | **エラーケース**: メソッド検証   |

**共通実装の導出**:

```go
// メソッド検証（Point 3から導出）
if r.Method != http.MethodGet {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
}
// Point 1, 2は正常処理へ進む
```

#### `parseParams` メソッド

| Point   | クエリパラメータ    | 解析結果         | エラー | 差分の本質                       |
| ------- | ------------------- | ---------------- | ------ | -------------------------------- |
| Point 1 | `service_area=SA11` | デフォルト値適用 | なし   | **必須のみ**: 最小入力           |
| Point 2 | 全パラメータ指定    | 全フィールド設定 | なし   | **全指定**: 最大入力             |
| Point 3 | `service_area`なし  | 空の構造体       | あり   | **バリデーション**: 必須チェック |

**共通実装の導出**:

```go
// 必須パラメータのバリデーション（Point 3から導出）
serviceArea := query.Get("service_area")
if serviceArea == "" {
    return params.GourmetSearchParams{}, &ValidationError{...}
}

// オプションパラメータはデフォルト値付き（Point 1, 2から導出）
start := parseIntOrDefault(query.Get("start"), 1)
count := parseIntOrDefault(query.Get("count"), 20)
```

#### `GenreHandler.Handle` メソッド

| Point   | HTTP メソッド | 期待ステータス           | 差分の本質                               |
| ------- | ------------- | ------------------------ | ---------------------------------------- |
| Point 1 | `GET`         | `200 OK`                 | **正常**: 許可されたメソッド             |
| Point 2 | `POST`        | `405 Method Not Allowed` | **エラー 1**: 不正メソッド（書き込み系） |
| Point 3 | `DELETE`      | `405 Method Not Allowed` | **エラー 2**: 不正メソッド（削除系）     |

**共通実装の導出**:

```go
// GETメソッドのみ許可（全Pointから導出）
if r.Method != http.MethodGet {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
}
```

#### `ValidationError` 型

`ValidationError`は`handler/gourmet.go`に定義されているカスタムエラー型です（`gourmet`ハンドラー固有のため、`params`パッケージではなくハンドラー内に配置）。

| Point   | メッセージ                   | 戻り値 | 差分の本質                 |
| ------- | ---------------------------- | ------ | -------------------------- |
| Point 1 | `"service_area is required"` | 同じ   | **標準**: 英語メッセージ   |
| Point 2 | `""` (空)                    | `""`   | **エッジケース**: 空文字列 |
| Point 3 | `"サービスエリアは必須です"` | 同じ   | **国際化**: 多言語対応     |

**共通実装の導出**:

```go
// シンプルなerrorインターフェース実装（全Pointから導出）
// internal/handler/gourmet.go に配置
func (e *ValidationError) Error() string {
    return e.Message  // メッセージをそのまま返す
}
```

---

## 🎯 Triangulation の効果

### 1. 共通箇所から得られるもの

- **統一されたテスト構造**: すべてのテストが同じパターンに従う
- **再利用可能なモック**: `httptest`による一貫したモック手法
- **保守性**: 新しいテストケースの追加が容易
- **可読性**: コードの意図が明確

### 2. 差分から得られるもの

- **網羅性**: 正常系、境界値、エラー系を 3 点でカバー
- **実装の一般化**: 複数のケースから共通パターンを導出
- **バグの早期発見**: 境界値テストで潜在的なバグを検出
- **仕様の明確化**: テストケースが仕様書として機能

### 3. 最小限のテストで最大限の効果

```
9つのテスト関数 × 3つのポイント = 27のテストケース

これだけで主要機能の68.0%をカバー
```

---

## 📋 テスト一覧表

| テスト関数                      | Point 1            | Point 2        | Point 3              | 共通箇所                       | 差分の種類          |
| ------------------------------- | ------------------ | -------------- | -------------------- | ------------------------------ | ------------------- |
| `TestClampInt`                  | 範囲内             | 下限           | 上限                 | 関数呼び出し、戻り値検証       | 入力値の違い        |
| `TestParseIntOrDefault`         | 空文字列           | 正常変換       | 不正値               | 同上                           | 文字列の型          |
| `TestSearchGourmet`             | 必須パラメータ欠如 | 必須のみ       | 全パラメータ         | HTTP モック、URL 検証          | パラメータの有無    |
| `TestGetGenreMaster`            | 正常               | 空データ       | エラー               | HTTP モック、ステータス確認    | レスポンス種類      |
| `TestSearchGourmetPagination`   | start 下限         | count 下限     | count 上限           | HTTP モック、補正検証          | 境界値の種類        |
| `TestGourmetHandlerHandle`      | GET 正常（必須）   | GET 正常（全） | POST                 | ハンドラー実行、ステータス確認 | メソッド/パラメータ |
| `TestGourmetHandlerParseParams` | 必須のみ           | 全パラメータ   | バリデーションエラー | パラメータ解析、構造体検証     | パラメータの有無    |
| `TestGenreHandlerHandle`        | GET                | POST           | DELETE               | ハンドラー実行、ステータス確認 | HTTP メソッド       |
| `TestValidationError`           | 英語               | 空             | 日本語               | Error()呼び出し、文字列比較    | メッセージ内容      |

---

## 🔍 実装導出の例

### 例 1: clampInt 関数

```go
// Point 1: value=5 (範囲内) → 5を返す
// Point 2: value=0 (下限外) → 1を返す
// Point 3: value=15 (上限外) → 10を返す

// この3つのケースから以下の一般的な実装が導かれる:
func clampInt(value, min, max int) int {
    if value < min {  // Point 2の条件
        return min    // Point 2の戻り値
    }
    if value > max {  // Point 3の条件
        return max    // Point 3の戻り値
    }
    return value      // Point 1の戻り値
}
```

### 例 2: オプションパラメータの処理

```go
// Point 1: Addressなし → URLにaddressパラメータなし
// Point 2: Address="さいたま" → URLに"address=..."追加
// Point 3: Address + Genre + Keyword → 全て追加

// この3つのケースから条件付き追加のパターンが導かれる:
if params.Address != "" {      // Point 2, 3で必要
    queryParams.Set("address", params.Address)
}
if params.Genre != "" {        // Point 3で必要
    queryParams.Set("genre", params.Genre)
}
if params.Keyword != "" {      // Point 3で必要
    queryParams.Set("keyword", params.Keyword)
}
```

### 例 3: HTTP メソッドのバリデーション

```go
// Point 1: GET → 200 OK
// Point 2: POST → 405 Method Not Allowed
// Point 3: DELETE → 405 Method Not Allowed

// この3つのケースからGETのみ許可するパターンが導かれる:
if r.Method != http.MethodGet {  // Point 2, 3の条件
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
}
// Point 1は正常処理へ
```

---

## � パッケージ構成

### `internal/params`

検索パラメータを定義するパッケージ。各種検索リクエストのパラメータ構造体を管理します。

- **`params.go`**: `GourmetSearchParams` - グルメ検索のパラメータ定義
- 旧名: `internal/model` → より明確な命名に変更

**設計思想**:

- ✅ 内容が明確（パラメータ専用）
- ✅ 拡張性あり（将来的に他の検索パラメータも追加可能）
- ✅ レイヤー間のデータ転送に使用

### `internal/handler`

HTTP リクエストを処理するハンドラー層。

- **`gourmet.go`**: グルメ検索 API、`ValidationError`型も含む
- **`genre.go`**: ジャンルマスター取得 API

**`ValidationError`の配置理由**:

- 現在`gourmet.go`でのみ使用されている
- YAGNI 原則（必要になってから実装する）に従い、使用箇所に近い場所に配置
- 将来、他のハンドラーでも使う必要が出たら`params`や専用パッケージに移動を検討

---

## �📝 まとめ

### Triangulation の本質

1. **3 点で検証**: 最小限のテストケースで実装を導出
2. **共通箇所**: すべてのテストで共有される構造とパターン
3. **差分**: 各ケース固有の検証ポイント（入力、境界値、エラー）
4. **実装の一般化**: 複数の具体例から抽象的なパターンを抽出

### このアプローチの利点

- ✅ **効率性**: 27 のテストケースで主要機能を網羅
- ✅ **明確性**: 各テストの目的と検証内容が明確
- ✅ **保守性**: テーブル駆動による拡張の容易さ
- ✅ **品質**: 境界値テストによる高い信頼性
- ✅ **ドキュメント性**: テストがそのまま仕様書として機能

### フロントエンドとの一貫性

バックエンド（Go）とフロントエンド（TypeScript）の両方で同じ Triangulation の原則を適用することで、プロジェクト全体で統一されたテスト戦略を実現しています。
