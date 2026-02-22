# Types Package

このパッケージは、ホットペッパーグルメ API のリクエスト・レスポンスに対応する型定義を提供します。

## 概要

`types`パッケージには、以下の API 仕様に基づいた構造体が定義されています：

- グルメサーチ API (`/api/gourmet`)
  - リクエストパラメータ
  - レスポンス構造
- ジャンルマスタ API (`/api/genre`)
  - レスポンス構造

## リクエスト型定義

### GourmetSearchParams

グルメサーチ API のリクエストパラメータを表す構造体です。

```go
type GourmetSearchParams struct {
    Address string  // オプション：住所のキーワード
    Genre   string  // オプション：ジャンルコード
    Keyword string  // オプション：フリーワード検索
    Start   int     // ページング：検索開始位置（1〜1000）
    Count   int     // ページング：取得件数（1〜100）
    Lat     float64 // 緯度（lat/lng検索時）
    Lng     float64 // 経度（lat/lng検索時）
    Range   int     // 検索範囲（1〜5、lat/lng検索時）
}
```

**パラメータ詳細:**

- `Address`: 住所の部分一致検索に使用するキーワード
- `Genre`: ジャンルコード（G001=居酒屋など）
- `Keyword`: 店名やキャッチコピーなどのフリーワード検索
- `Start`: 検索結果の開始位置（1〜1000 の範囲、デフォルト 1）
- `Count`: 1 回のリクエストで取得する件数（1〜100 の範囲、デフォルト 20）
- `Lat` / `Lng`: 現在地検索で使用する緯度・経度（両方指定時に有効）
- `Range`: 現在地検索の範囲（1〜5）

**入力ルール:**

- `lat/lng` の組み合わせ、`address`、`keyword` のいずれか 1 つ以上が必要

## レスポンス型定義

### グルメサーチ API 関連

#### GourmetSearchResponse

グルメサーチ API のレスポンス全体を表す構造体です。

```go
type GourmetSearchResponse struct {
    Results GourmetSearchResults `json:"results"`
}
```

#### GourmetSearchResults

検索結果の詳細情報を表す構造体です。

```go
type GourmetSearchResults struct {
    APIVersion       string `json:"api_version"`       // APIバージョン
    ResultsAvailable int    `json:"results_available"` // 検索条件に該当する総件数
    ResultsReturned  string `json:"results_returned"`  // 返却された件数（APIから文字列で返される）
    ResultsStart     int    `json:"results_start"`     // 開始位置
    Shop             []Shop `json:"shop"`              // 店舗情報の配列
}
```

#### Shop

店舗情報を表す構造体です。

```go
type Shop struct {
    ID      string    `json:"id"`      // 店舗ID
    Name    string    `json:"name"`    // 店舗名
    Address string    `json:"address"` // 住所
    Lat     float64   `json:"lat"`     // 緯度
    Lng     float64   `json:"lng"`     // 経度
    Genre   ShopGenre `json:"genre"`   // ジャンル情報
    Catch   string    `json:"catch"`   // キャッチコピー
    Access  string    `json:"access"`  // アクセス情報
    URLs    ShopURLs  `json:"urls"`    // URL情報
    Photo   ShopPhoto `json:"photo"`   // 写真情報
}
```

#### その他の構造体

- `ShopGenre`: 店舗のジャンル情報
- `ShopURLs`: 店舗の URL 情報
- `ShopPhoto`: 店舗の写真情報
- `ShopPhotoPC`: PC 向け写真の各サイズの URL

### ジャンルマスタ API 関連

#### GenreMasterResponse

ジャンルマスタ API のレスポンス全体を表す構造体です。

```go
type GenreMasterResponse struct {
    Results GenreMasterResults `json:"results"`
}
```

#### GenreMasterResults

ジャンルマスタの詳細情報を表す構造体です。

```go
type GenreMasterResults struct {
    APIVersion       string  `json:"api_version"`       // APIバージョン
    ResultsAvailable int     `json:"results_available"` // 総件数
    ResultsReturned  string  `json:"results_returned"`  // 返却された件数（APIから文字列で返される）
    ResultsStart     int     `json:"results_start"`     // 開始位置
    Genre            []Genre `json:"genre"`             // ジャンル情報の配列
}
```

#### Genre

ジャンル情報を表す構造体です。

```go
type Genre struct {
    Code string `json:"code"` // ジャンルコード
    Name string `json:"name"` // ジャンル名
}
```

## 使用例

### リクエストパラメータの作成

```go
import "backend/internal/types"

// 検索パラメータを作成
params := types.GourmetSearchParams{
    Address: "新宿",
    Genre:   "G001",
    Keyword: "個室",
    Lat:     35.6895,
    Lng:     139.6917,
    Range:   3,
    Start:   1,
    Count:   20,
}

// サービス層に渡して検索実行
response, err := service.SearchGourmet(params)
```

### レスポンスの処理

```go
// 検索結果からデータを取得
if response != nil {
    fmt.Printf("総件数: %d\n", response.Results.ResultsAvailable)
    fmt.Printf("返却件数: %s\n", response.Results.ResultsReturned)

    for _, shop := range response.Results.Shop {
        fmt.Printf("店舗名: %s\n", shop.Name)
        fmt.Printf("住所: %s\n", shop.Address)
        fmt.Printf("ジャンル: %s\n", shop.Genre.Name)
    }
}
```

### ジャンルマスタの処理

```go
// ジャンル一覧を取得
genreResponse, err := service.GetGenreMaster()
if err != nil {
    log.Fatal(err)
}

// ジャンル情報を表示
for _, genre := range genreResponse.Results.Genre {
    fmt.Printf("%s: %s\n", genre.Code, genre.Name)
}
```

## パッケージ統合の経緯

以前は`params`パッケージと`types`パッケージが分離していましたが、以下の理由で統合しました：

### 統合のメリット

1. **一貫性**: リクエストとレスポンスの型定義が同じパッケージに集約
2. **可読性**: 関連する型定義が一箇所にまとまり、コードが読みやすい
3. **保守性**: パッケージ数が減り、依存関係がシンプルに
4. **論理的な構造**: API に関連する型定義が論理的にグループ化

### ファイル構成

```
internal/types/
├── request.go   # リクエストパラメータの型定義
├── response.go  # レスポンスの型定義
└── README.md    # このファイル
```

## 設計思想

### 1. 型安全性

構造体を使用することで、コンパイル時に型の不一致を検出できます。

### 2. JSON タグの使用

レスポンス構造体の各フィールドには適切な`json`タグが付与されており、JSON のエンコード/デコードがスムーズに行えます。

### 3. API 仕様との整合性

`api-spec.yaml`に定義された OpenAPI 仕様と完全に一致する構造になっています。

### 4. 明確なドキュメント

型定義に詳細なコメントを付与し、各フィールドの意味を明確にしています。

## メリット

1. **型安全性**: コンパイル時にエラーを検出
2. **IDE サポート**: 自動補完やリファクタリング機能の利用が可能
3. **ドキュメント化**: 型定義がそのままドキュメントとして機能
4. **保守性**: 構造の変更が明示的になり、影響範囲を把握しやすい
5. **テスト容易性**: モックデータの作成が容易

## 参考

- API 仕様書: `/api-spec.yaml`
- サービス層の実装: `/internal/service/hotpepper.go`
- ハンドラー層の実装: `/internal/handler/gourmet.go`, `/internal/handler/genre.go`
