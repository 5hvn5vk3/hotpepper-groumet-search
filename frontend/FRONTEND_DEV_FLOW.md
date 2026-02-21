# フロントエンド開発フロー

バックエンド（API）が完成している前提で、React + TypeScript + Vite を使用したフロントエンドをゼロから構築する際の標準的な手順です。
依存関係の少ない「データ・通信層」から「UI層」へと積み上げていくことで、手戻りを防ぎ効率的に開発を進めることができます。

## 1. プロジェクトの初期化と環境構築

まずはアプリケーションの土台を作成します。

1.  **プロジェクト作成**
    ```bash
    npm create vite@latest frontend -- --template react-ts
    ```
2.  **ライブラリの導入**
    - スタイリング: Tailwind CSS
    - テスト: Vitest, React Testing Library
    - その他: 必要に応じてルーティング (React Router) など
3.  **設定ファイルの整備**
    - `vite.config.ts`: エイリアス（`@/`）の設定など
    - `tsconfig.json`: TypeScriptのコンパイル設定

## 2. 型定義の作成 (Domain Layer)

バックエンドとの「契約」となるデータ構造を定義します。

- **作成ファイル**: `src/types.ts`
- **目的**: APIレスポンスやコンポーネント間のデータ受け渡しで型安全性を確保するため。
- **作業内容**:
  - API仕様書（Swagger/OpenAPIなど）を参照。
  - レスポンスのJSON構造に合わせてインターフェース（`Shop`, `GourmetSearchResponse`, `Genre` など）を定義。

## 3. APIクライアントの実装 (Infrastructure Layer)

定義した型を使用して、実際にバックエンドと通信する処理を実装します。

- **作成ファイル**:
  - `src/api/client.ts` (共通クライアント)
  - `src/api/genreApi.ts`, `src/api/restaurantApi.ts` (個別API)
- **作業内容**:
  - `fetch` または `axios` をラップした汎用クライアントを作成（ベースURL結合、エラーハンドリング共通化）。
  - 各エンドポイントに対応する関数を実装し、戻り値に型定義を適用。

## 4. カスタムフックの実装 (Logic Layer)

UIからデータ取得ロジックを分離し、状態管理（ローディング、エラー、データ保持）を行うフックを作成します。

- **作成ファイル**: `src/hooks/useGenres.ts`, `src/hooks/useRestaurantSearch.ts`
- **目的**: UIコンポーネントを純粋な表示に集中させ、ロジックの再利用性とテスト容易性を高める（Headless UIパターン）。
- **作業内容**:
  - `useState` で `data`, `isLoading`, `error` を管理。
  - `useEffect` やイベントハンドラ内でAPI関数を呼び出す。

## 5. UIコンポーネントの実装 (Presentation Layer)

ロジック部分ができたら、それを表示するUI部品をボトムアップで作成します。

- **作成ファイル**: `src/components/` 配下
- **順序**:
  1.  **共通部品**: `LoadingSpinner.tsx`, `ErrorMessage.tsx`
  2.  **機能部品**:
      - `SearchForm.tsx`: 検索条件入力（`useGenres` を使用）
      - `RestaurantList.tsx`: 結果表示
      - `Pagination.tsx`: ページネーション

## 6. アプリケーションの統合 (Integration)

作成したフックとコンポーネントを組み合わせ、アプリケーションとして完成させます。

- **作成ファイル**: `src/App.tsx`
- **作業内容**:
  - `useRestaurantSearch` フックで検索状態を管理。
  - `SearchForm` に検索実行関数を渡す。
  - `RestaurantList` に検索結果データを渡す。
  - 全体のレイアウト調整。

## 💡 開発のヒント：機能単位で進める（Vertical Slice）

上記の手順は「レイヤー（層）」ごとの解説ですが、実際の開発では**1つの機能ごとにこのサイクルを回す**のがおすすめです。
これを「垂直スライス（Vertical Slice）」開発と呼びます。

### 垂直スライス開発の具体的な手順

#### サイクル1：ジャンル一覧機能（まずは小さな成功体験を）

1.  **型定義**: `Genre` 型を定義する。
2.  **API**: `getGenres` 関数を実装し、コンソールでデータが取れるか確認する。
3.  **フック**: `useGenres` を作り、データ取得状態を管理できるようにする。
4.  **UI**: `SearchForm` のプルダウン部分だけを作り、ジャンルが表示されることを確認する。
    - 👉 **ここで一度画面が動く！**

#### サイクル2：店舗検索機能（メイン機能の実装）

1.  **型定義**: `Shop`, `GourmetSearchResponse` 型を定義する。
2.  **API**: `searchRestaurants` 関数を実装する。
3.  **フック**: `useRestaurantSearch` を作り、検索ボタンを押したらデータが変わるようにする。
4.  **UI**: `RestaurantList` を作り、検索結果がリスト表示されるようにする。
    - 👉 **検索して結果が出る！**

#### サイクル3：詳細表示・ページネーション（機能の肉付け）

1.  **UI**: `Pagination` コンポーネントを追加し、ページ送りを実装する。
2.  **UI**: `RestaurantDetail` モーダルを追加し、詳細情報を見れるようにする。
3.  **統合**: 全体のレイアウトやエラー処理をブラッシュアップする。

こうすることで、**「動く画面」をこまめに見ることができ、モチベーションを維持しやすくなります。**
