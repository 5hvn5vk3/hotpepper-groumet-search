# ホットペッパーグルメサーチアプリケーション

ホットペッパーのグルメサーチ API、ジャンルマスタ API を使用したレストラン検索アプリケーションです。

### フロントエンド（React + TypeScript）

```
frontend/
├── src/
│   ├── components/         # プレゼンテーション層（UIコンポーネント）
|   |   ├── index.ts
│   │   ├── SearchForm.tsx
│   │   ├── RestaurantList.tsx
│   │   ├── RestaurantDetail.tsx
│   │   ├── Pagination.tsx
│   │   ├── ErrorMessage.tsx
│   │   └── LoadingSpinner.tsx
│   │
│   ├── hooks/             # カスタムフック（ビジネスロジックとステート管理）
│   │   ├── useRestaurantSearch.ts
│   │   ├── useGenres.ts
│   │   └── useModal.ts
│   │
│   ├── api/               # API通信層（バックエンドとの通信）
│   │   ├── client.ts      # 汎用HTTPクライアント
│   │   ├── restaurantApi.ts
│   │   └── genreApi.ts
│   │
│   ├── services/          # サービス層（ビジネスロジックとユーティリティ）
│   │   └── storageService.ts
│   │
│   ├── types.ts           # 型定義
│   ├── constants.ts       # 定数
│   └── App.tsx            # ルートコンポーネント
```

#### レイヤーの責務

- **Components**: UI の表示のみを担当。ビジネスロジックを含まない
- **Hooks**: ステート管理とビジネスロジックを担当
- **API**: バックエンドとの HTTP 通信を担当
- **Services**: ビジネスロジック、データの永続化、ユーティリティを担当

### バックエンド（Go）

```
backend/
├── main.go               # アプリケーションエントリーポイント
├── go.mod               # Goモジュール定義
└── internal/            # 内部パッケージ（外部から非公開）
    ├── handler/         # HTTPハンドラー層
    │   ├── gourmet.go   # グルメ検索ハンドラー（ValidationError型も含む）
    │   ├── genre.go     # ジャンルマスターハンドラー
    │   └── handler_test.go
    │
    ├── service/         # サービス層（ビジネスロジック）
    │   ├── hotpepper.go # ホットペッパーAPI通信
    │   └── service_test.go
    │
    ├── middleware/      # HTTPミドルウェア
    │   └── cors.go      # CORS処理
    │
    └── types/           # 型定義
        ├── request.go   # リクエストパラメータの型定義
        ├── response.go  # レスポンスの型定義
        └── README.md    # 型定義のドキュメント
```

#### レイヤーの責務

- **main.go**: アプリケーションの初期化とルーティング設定、各パッケージの統合
- **handler/**: HTTP リクエスト/レスポンスの処理、バリデーション、パラメータ解析
  - `ValidationError`型もここに配置（ハンドラー層固有のエラー型）
- **service/**: ホットペッパー API との通信ロジック、ビジネスルール、型付きレスポンス処理
- **middleware/**: HTTP ミドルウェア（CORS、ロギング等の横断的関心事）
- **types/**: API のリクエスト・レスポンスの型定義（レイヤー間で共有されるデータ構造）
  - リクエストパラメータ: `GourmetSearchParams`
  - レスポンス: `GourmetSearchResponse`, `GenreMasterResponse`など

## 技術スタック

### フロントエンド

- React 18
- TypeScript
- Vite
- TailwindCSS
- Vitest（テストフレームワーク）

### バックエンド

- Go 1.x
- 標準ライブラリのみ使用
- 型安全な構造体によるリクエスト・レスポンス処理

## セットアップ

### 環境変数

バックエンドに`.env`ファイルを作成し、以下を設定：

```
HOTPEPPER_API_KEY=your_api_key_here
PORT=8080
```

### フロントエンド

```bash
cd frontend
npm install
npm run dev
```

### バックエンド

```bash
cd backend
go run main.go
```

### テストの実行

```bash
# フロントエンド
cd frontend
npm test

# フロントエンド（カバレッジ付き）
npm run test:coverage

# バックエンド
cd backend
go test ./...

# バックエンド（詳細表示）
go test ./... -v

# バックエンド（カバレッジ付き）
go test ./... -cover
```

## 開発の原則

### 1. 関心の分離（Separation of Concerns）

各モジュールは単一の責務を持ち、他のモジュールと疎結合に保つ。

### 2. 型安全性（Type Safety）

- フロントエンド: TypeScript による静的型チェック
- バックエンド: Go の構造体による型安全なリクエスト・レスポンス処理
- API 仕様（OpenAPI）と型定義の整合性を保つ

### 3. 依存性の注入（Dependency Injection）

- バックエンド: サービスをハンドラーに注入
- フロントエンド: サービスをフックから呼び出し

### 4. 再利用性

- 共通ロジックはサービス層に集約
- UI コンポーネントは純粋に保ち、再利用可能にする
- 型定義は `types` パッケージに集約

### 5. テスタビリティ

各レイヤーが独立しているため、単体テストが容易。

- フロントエンド: Vitest によるユニットテスト
- バックエンド: Go 標準の testing パッケージ、モックサーバーを使用

## ディレクトリ構造の意図

### フロントエンド

- **components/**: 表示に専念。props を受け取り、UI を返す
- **hooks/**: ビジネスロジックをカプセル化。コンポーネントから状態管理を分離
- **api/**: バックエンドとの HTTP 通信を抽象化
- **services/**: ビジネスロジック、データの永続化（LocalStorage 等）を担当

### バックエンド

- **handler/**: HTTP レイヤーの処理（リクエスト解析、レスポンス生成、バリデーション）
- **service/**: ビジネスロジック（外部 API 呼び出し、JSON デコード、エラーハンドリング）
- **middleware/**: 横断的関心事（CORS、認証、ロギング等）
- **types/**: API のリクエスト・レスポンスの型定義
  - リクエスト型: API へのパラメータを表現
  - レスポンス型: API からの JSON レスポンスを構造化

#### パッケージ設計の原則

1. **internal/パッケージ**: 外部から非公開にすることで、API の安定性を保つ
2. **単一責任の原則**: 各パッケージは明確な責務を持つ
3. **依存関係の方向**: handler → service → types（一方向の依存）
4. **型安全性**: 構造体による型定義で、コンパイル時にエラーを検出
5. **テスタビリティ**: 各パッケージが独立してテスト可能
6. **統合されたパッケージ**: リクエストとレスポンスの型定義を `types` パッケージに集約

## API

API の詳細仕様は `api-spec.yaml`（OpenAPI 3.0 形式）を参照してください。

### GET /api/gourmet

レストラン検索

**パラメータ:**

- `lat` (任意): 緯度（`lng` とセットで現在地検索）
- `lng` (任意): 経度（`lat` とセットで現在地検索）
- `range` (任意): 検索範囲（1〜5、`lat/lng` 指定時のみ有効）
- `address` (任意): 住所（部分一致検索）
- `genre` (任意): ジャンルコード（例: G001=居酒屋）
- `keyword` (任意): キーワード（店名、キャッチ等）
- `start` (任意): 開始位置（1〜1000、デフォルト: 1）
- `count` (任意): 取得件数（1〜100、デフォルト: 20）

**検索実行条件:**

- `lat/lng` の組み合わせ、`address`、`keyword` のいずれか 1 つ以上が必要

**レスポンス:**
型定義された構造体（`types.GourmetSearchResponse`）で返却。
店舗情報の配列を含む。

### GET /api/genre

ジャンル一覧取得

**レスポンス:**
型定義された構造体（`types.GenreMasterResponse`）で返却。
ジャンル情報の配列を含む。

## 型定義

### バックエンド

バックエンドの型定義の詳細は `backend/internal/types/README.md` を参照してください。

主な型定義:

- **リクエスト**: `GourmetSearchParams`
- **レスポンス**: `GourmetSearchResponse`, `GenreMasterResponse`, `Shop`, `Genre` など

### フロントエンド

フロントエンドの型定義は `frontend/src/types.ts` に定義されています。

主な型定義:

- **検索パラメータ**: `SearchParams`
- **レスポンス**: `GourmetSearchResponse`, `GenreMasterResponse`
- **店舗情報**: `Restaurant`
- **ジャンル**: `Genre`
