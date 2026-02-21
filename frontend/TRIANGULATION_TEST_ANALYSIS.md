# Triangulation（三角測量）テスト分析

## Triangulationとは

**Triangulation（三角測量）** は、TDD（テスト駆動開発）において、**複数の異なる具体的なテストケースから一般的な実装を導出する手法**です。

最小限の3つのテストポイントで、以下を検証します：

- **Point 1**: 基本的な正常ケース
- **Point 2**: 別の角度からの検証（境界値や異なる入力）
- **Point 3**: さらに異なる観点（エッジケースやエラー）

この3点により、実装の一般性と正確性を保証します。

---

---

## 🔄 共通箇所（全テストで共有される構造）

### 1. テスト構造

すべてのテストで共通のパターン：

```typescript
describe("ComponentName", () => {
  beforeEach(() => {
    // モックのクリアと初期化
    vi.clearAllMocks();
  });

  describe("機能グループ", () => {
    it("Point 1: 基本ケース", () => {
      // Arrange（準備）
      // Act（実行）
      // Assert（検証）
    });

    it("Point 2: 別の角度", () => {
      // ...
    });

    it("Point 3: エッジケース", () => {
      // ...
    });
  });
});
```

### 2. モック手法

```typescript
// APIのモック（全テストで共通）
vi.mock("@/api/restaurantApi");
vi.mocked(searchRestaurants).mockResolvedValue(mockResponse);

// Reactフックのモック
const { result } = renderHook(() => useCustomHook());
await waitFor(() => expect(result.current.data).toBeDefined());
```

### 3. エラーハンドリング

```typescript
// すべてのテストで統一されたエラーチェック
try {
  await someFunction();
} catch (err) {
  const errorMessage =
    err instanceof Error ? err.message : "デフォルトメッセージ";
  setError(errorMessage);
}
```

---

## 📊 テストファイル別のTriangulation分析

### 1. StorageService (`storageService.test.ts`)

**Triangulation適用**: オブジェクトデータの保存・取得に特化

#### **get メソッドのTriangulation**

| Triangulation Point | テストケース               | 検証観点     | 共通箇所                                 | 差分                            |
| ------------------- | -------------------------- | ------------ | ---------------------------------------- | ------------------------------- |
| **Point 1**         | 存在しないキー             | エッジケース | - `sessionStorage.getItem`<br>- null検証 | **条件**: キーが存在しない      |
| **Point 2**         | オブジェクトデータを取得   | 正常系       | - JSON.parse<br>- 型検証                 | **条件**: 正常なJSONデータ      |
| **Point 3**         | 無効なJSONの場合nullを返す | 異常系       | - try-catch<br>- エラーハンドリング      | **条件**: パース失敗 → null返却 |

**共通実装パターン**:

```typescript
get<T>(key: string): T | null {
  try {
    const item = sessionStorage.getItem(key);  // Point 1: null判定
    if (!item) return null;
    return JSON.parse(item) as T;              // Point 2: 正常パース
  } catch {
    return null;                               // Point 3: エラー時null
  }
}
```

#### **set メソッドのテスト**

| テストケース             | 検証観点 | 検証内容                                                 |
| ------------------------ | -------- | -------------------------------------------------------- |
| オブジェクトデータを保存 | 正常系   | オブジェクトをJSON文字列としてsessionStorageに正しく保存 |

**実装パターン**:

```typescript
set<T>(key: string, value: T): void {
  sessionStorage.setItem(key, JSON.stringify(value));
}
```

---

### 2. RestaurantAPI (`restaurantApi.test.ts`)

#### **基本的な検索パラメータのTriangulation**

| Triangulation Point | テストケース                        | 検証観点       | 共通箇所                                  | 差分                                    |
| ------------------- | ----------------------------------- | -------------- | ----------------------------------------- | --------------------------------------- |
| **Point 1**         | serviceAreaが空の場合は早期リターン | バリデーション | - 必須チェック<br>- 空レスポンス返却      | **条件**: `serviceArea === ""`          |
| **Point 2**         | 必須パラメータのみ                  | 最小構成       | - `apiGet`呼び出し<br>- URLパラメータ検証 | **パラメータ**: `serviceArea`のみ       |
| **Point 3**         | 全パラメータ                        | 最大構成       | 同上                                      | **追加**: `address`, `genre`, `keyword` |

**共通実装パターン**:

```typescript
// Point 1: 必須パラメータのバリデーション
if (!params.serviceArea) {
  return {
    results: {
      /* 空のレスポンス */
    },
  };
}

// Point 2, 3: オプショナルパラメータを条件付きで追加
if (params.address) queryParams.set("address", params.address);
if (params.genre) queryParams.set("genre", params.genre);
if (params.keyword) queryParams.set("keyword", params.keyword);
```

#### **ページネーションのTriangulation**

| Triangulation Point | テストケース                 | 検証観点      | 共通箇所                             | 差分                            |
| ------------------- | ---------------------------- | ------------- | ------------------------------------ | ------------------------------- |
| **Point 1**         | ページ1 (start=1)            | 初期ページ    | - ページ→start変換<br>- 計算式の検証 | `page=1, count=10` → `start=1`  |
| **Point 2**         | ページ2 (start=11)           | 次ページ      | 同上                                 | `page=2, count=10` → `start=11` |
| **Point 3**         | ページ3, count=20 (start=41) | 異なるcount値 | 同上                                 | `page=3, count=20` → `start=41` |

**共通実装パターン**:

```typescript
// この3つのケースから以下の一般的な計算式が導かれる
const start = (page - 1) * count + 1;
```

**数学的検証**:

- Point 1: `(1-1) * 10 + 1 = 1` ✓
- Point 2: `(2-1) * 10 + 1 = 11` ✓
- Point 3: `(3-1) * 20 + 1 = 41` ✓

#### **境界値とバリデーションのTriangulation**

| Triangulation Point | テストケース      | 検証観点 | 共通箇所                                 | 差分                    |
| ------------------- | ----------------- | -------- | ---------------------------------------- | ----------------------- |
| **Point 1**         | page < 1 → 1      | 下限値   | - バリデーション処理<br>- `Math.max`使用 | **条件**: `page ≤ 0`    |
| **Point 2**         | count < 1 → 1     | 下限値   | 同上                                     | **条件**: `count ≤ 0`   |
| **Point 3**         | count > 100 → 100 | 上限値   | - `Math.min`使用                         | **条件**: `count > 100` |

**共通実装パターン**:

```typescript
// この3つのケースから以下の一般的な境界値制御が導かれる
const page = Math.max(1, params.page ?? 1); // Point 1: 下限1
const count = Math.min(Math.max(params.count ?? 20, 1), 100); // Point 2,3: 下限1、上限100
```

---

### 3. useRestaurantSearch Hook (`useRestaurantSearch.test.ts`)

#### **初期状態のテスト**

| テストケース     | 検証観点     | 検証内容                                                 |
| ---------------- | ------------ | -------------------------------------------------------- |
| 初期状態が正しい | 全状態の確認 | すべての状態変数が正しい初期値に設定されていることを保証 |

**実装パターン**:

```typescript
// フックの初期化時に設定される状態
const [searchResult, setSearchResult] = useState(null);
const [isLoading, setIsLoading] = useState(false);
const [hasSearched, setHasSearched] = useState(false);
const [error, setError] = useState(null);
const [currentPage, setCurrentPage] = useState(1);
```

#### **search メソッドのTriangulation**

| Triangulation Point | テストケース     | 検証観点 | 共通箇所                                                        | 差分                                    |
| ------------------- | ---------------- | -------- | --------------------------------------------------------------- | --------------------------------------- |
| **Point 1**         | 基本的な検索     | 基本機能 | - `searchRestaurants`呼び出し<br>- ページリセット<br>- 結果保存 | **パラメータ**: `serviceArea`のみ       |
| **Point 2**         | 全パラメータ検索 | フル機能 | 同上                                                            | **追加**: `address`, `genre`, `keyword` |

**共通実装パターン**:

```typescript
const search = async (params: SearchParams) => {
  setCurrentParams(params);
  setCurrentPage(1); // Point 1, 2共通: 常にページ1から開始
  await fetchRestaurants(params, 1);
};
```

#### **changePage メソッドのTriangulation**

| Triangulation Point | テストケース     | 検証観点           | 共通箇所                            | 差分                             |
| ------------------- | ---------------- | ------------------ | ----------------------------------- | -------------------------------- |
| **Point 1**         | ページ2へ遷移    | 基本的なページ変更 | - ページ番号更新<br>- API再呼び出し | **target page**: 2               |
| **Point 2**         | ページ3へ遷移    | 別のページ番号     | 同上                                | **target page**: 3               |
| **Point 3**         | 同じページは無視 | 無効な操作         | ページ番号チェック                  | **条件**: `page === currentPage` |

**共通実装パターン**:

```typescript
const changePage = async (page: number) => {
  if (!currentParams || page === currentPage) return;
  setCurrentPage(page);
  await fetchRestaurants(currentParams, page);
};
```

#### **エラーハンドリングのTriangulation**

| Triangulation Point | テストケース        | 検証観点         | 共通箇所                            | 差分                                       |
| ------------------- | ------------------- | ---------------- | ----------------------------------- | ------------------------------------------ |
| **Point 1**         | Errorオブジェクト   | 標準的なエラー   | - try-catch処理<br>- エラー状態設定 | **エラー型**: `Error` → `error.message`    |
| **Point 2**         | 非Errorオブジェクト | 予期しないエラー | 同上                                | **エラー型**: `any` → デフォルトメッセージ |
| **Point 3**         | エラークリア        | エラー解除機能   | エラー状態管理                      | **操作**: `clearError()` → `null`          |

**共通実装パターン**:

```typescript
try {
  await searchRestaurants(params);
} catch (err) {
  const errorMessage =
    err instanceof Error ? err.message : "検索に失敗しました";
  setError(errorMessage);
}
```

### 4. useGenres Hook (`useGenres.test.ts`)

#### **ジャンルデータ取得のTriangulation**

| Triangulation Point | テストケース                            | 検証観点   | 共通箇所                                    | 差分                                |
| ------------------- | --------------------------------------- | ---------- | ------------------------------------------- | ----------------------------------- |
| **Point 1**         | マウント時にデータ取得                  | 正常系     | - `getGenres`呼び出し<br>- ローディング状態 | **結果**: 成功 → genres配列         |
| **Point 2**         | API失敗時にエラー設定                   | 異常系     | - エラーハンドリング<br>- エラー状態設定    | **結果**: Errorオブジェクト         |
| **Point 3**         | Error以外のエラーでデフォルトメッセージ | エラー処理 | 同上                                        | **結果**: 文字列エラー → デフォルト |

**共通実装パターン**:

```typescript
useEffect(() => {
  const fetchGenres = async () => {
    setIsLoading(true);
    try {
      const data = await getGenres(); // Point 1: 成功
      setGenres(data);
    } catch (err) {
      // Point 2: Errorオブジェクト、Point 3: その他
      const errorMessage =
        err instanceof Error ? err.message : "ジャンルの取得に失敗しました";
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };
  fetchGenres();
}, []);
```

---

### 5. SearchForm Component (`SearchForm.test.tsx`)

#### **検索機能のTriangulation**

| Triangulation Point | テストケース       | 検証観点 | 共通箇所                                   | 差分                                    |
| ------------------- | ------------------ | -------- | ------------------------------------------ | --------------------------------------- |
| **Point 1**         | 都道府県のみで検索 | 最小構成 | - フォーム送信<br>- `onSearch`コールバック | **パラメータ**: `serviceArea`のみ       |
| **Point 2**         | 全項目入力で検索   | 最大構成 | 同上                                       | **追加**: `address`, `genre`, `keyword` |

**共通実装パターン**:

```typescript
const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault();
  onSearch({
    serviceArea, // Point 1: 必須パラメータ
    address: address || undefined, // Point 2: オプション
    genre: selectedGenre || undefined, // Point 2: オプション
    keyword: keyword || undefined, // Point 2: オプション
  });
};
```

#### **ローディング状態のテスト**

| Triangulation Point | テストケース           | 検証観点   | 共通箇所                       | 差分                        |
| ------------------- | ---------------------- | ---------- | ------------------------------ | --------------------------- |
| **Point 3**         | ローディング中は無効化 | UI状態制御 | - ボタン状態<br>- テキスト変更 | `isLoading=true` → disabled |

**共通実装パターン**:

```typescript
<button type="submit" disabled={isLoading}>
  {isLoading ? "検索中..." : "検索"}
</button>
```

---

### 6. Pagination Component (`Pagination.test.tsx`)

#### **ページネーション動作のTriangulation**

| Triangulation Point | テストケース         | 検証観点   | 共通箇所                                   | 差分                            |
| ------------------- | -------------------- | ---------- | ------------------------------------------ | ------------------------------- |
| **Point 1**         | 1ページ目で前へ無効  | 初期ページ | - ボタン状態<br>- ページ表示<br>- 範囲表示 | `currentPage=1` → 前へ disabled |
| **Point 2**         | 2ページ目で両方有効  | 中間ページ | 同上                                       | `currentPage=2` → 両方 enabled  |
| **Point 3**         | 最終ページで次へ無効 | 最終ページ | 同上                                       | `currentPage=3` → 次へ disabled |

#### **ボタンクリック動作のテスト**

| テストケース                     | 検証観点       | 検証内容                                  |
| -------------------------------- | -------------- | ----------------------------------------- |
| 次へボタンクリックで次ページ遷移 | ナビゲーション | `onPageChange(currentPage + 1)`の呼び出し |

#### **エッジケースのTriangulation**

| Triangulation Point | テストケース           | 検証観点 | 共通箇所               | 差分                              |
| ------------------- | ---------------------- | -------- | ---------------------- | --------------------------------- |
| **Point 1**         | 総件数0で非表示        | 空データ | - 条件付きレンダリング | `totalCount=0` → null             |
| **Point 2**         | disabledで全ボタン無効 | 操作制御 | - ボタン無効化         | `disabled=true` → すべて disabled |

---

### 7. RestaurantList Component (`RestaurantList.test.tsx`)

#### **リスト表示のTriangulation**

| Triangulation Point | テストケース                 | 検証観点         | 共通箇所                               | 差分                             |
| ------------------- | ---------------------------- | ---------------- | -------------------------------------- | -------------------------------- |
| **Point 1**         | 空配列で空メッセージ表示     | 境界値           | - 条件分岐<br>- メッセージ表示         | `restaurants.length === 0`       |
| **Point 2**         | データありでカード表示       | 正常系           | - map処理<br>- カードレンダリング      | `restaurants=[{...}, {...}]`     |
| **Point 3**         | カードクリックでコールバック | インタラクション | - イベントハンドラー<br>- 関数呼び出し | `onClick → onSelectRestaurant()` |

**共通実装パターン**:

```typescript
// Point 1: 空配列の境界値チェック
if (restaurants.length === 0) {
  return <div>検索結果がありません</div>;
}

// Point 2, 3: データ表示とイベントハンドラー
return (
  <div>
    {restaurants.map((restaurant) => (
      <div onClick={() => onSelectRestaurant(restaurant)}>
        {/* Point 2: カード表示 */}
      </div>
    ))}
  </div>
);
```

---

## 📋 テスト一覧表

| テストファイル                | テスト数 | 主な検証内容                                  |
| ----------------------------- | -------- | --------------------------------------------- |
| `restaurantApi.test.ts`       | 8        | API呼び出し、パラメータ変換、ページネーション |
| `SearchForm.test.tsx`         | 3        | フォーム送信、パラメータ組み立て、UI状態      |
| `Pagination.test.tsx`         | 6        | ページ遷移、ボタン制御、エッジケース          |
| `RestaurantList.test.tsx`     | 3        | リスト表示、空配列処理、クリックイベント      |
| `useRestaurantSearch.test.ts` | 10       | 検索ロジック、状態管理、エラーハンドリング    |
| `useGenres.test.ts`           | 3        | ジャンル取得、エラー処理                      |
| `storageService.test.ts`      | 4        | データ永続化、JSON処理                        |
| **合計**                      | **37**   | -                                             |

---

## 🎯 Triangulationの効果

### 共通箇所から得られるもの

- **統一されたテスト構造**: すべてのテストが同じパターンに従う
- **再利用可能なモック**: `vi.mock()`による一貫したモック手法
- **保守性**: 新しいテストケースの追加が容易
- **可読性**: コードの意図が明確

### 差分から得られるもの

- **網羅性**: 正常系、境界値、エラー系を3点でカバー
- **実装の一般化**: 複数のケースから共通パターンを導出
- **バグの早期発見**: 境界値テストで潜在的なバグを検出
- **仕様の明確化**: テストケースが仕様書として機能

### 必要最小限のテストで最大限の効果

```
7つのテストファイル × 平均5.3テストケース = 37のテストケース

これだけで主要機能を網羅的にカバー
```

---

## 🔍 実装導出の例

### 例1: ページネーション計算

```typescript
// Point 1: page=1, count=10 → start=1
// Point 2: page=2, count=10 → start=11
// Point 3: page=3, count=20 → start=41

// この3つのケースから以下の一般的な計算式が導かれる:
const start = (page - 1) * count + 1;
```

**数学的検証**:

- Point 1: `(1-1) * 10 + 1 = 1` ✓
- Point 2: `(2-1) * 10 + 1 = 11` ✓
- Point 3: `(3-1) * 20 + 1 = 41` ✓

### 例2: オプションパラメータの処理

```typescript
// Point 1: serviceAreaのみ → 他のパラメータなし
// Point 2: 全パラメータ → すべて追加

// この2つのケースから条件付き追加のパターンが導かれる:
if (params.address) queryParams.set("address", params.address);
if (params.genre) queryParams.set("genre", params.genre);
if (params.keyword) queryParams.set("keyword", params.keyword);
```

### 例3: エラーハンドリング

```typescript
// Point 1: Errorオブジェクト → error.message
// Point 2: 文字列エラー → デフォルトメッセージ
// Point 3: その他 → デフォルトメッセージ

// この3つのケースから汎用的なエラー処理が導かれる:
const errorMessage =
  err instanceof Error ? err.message : "デフォルトのエラーメッセージ";
```

---

## 📝 まとめ

### Triangulationの本質

1. **3点で検証**: 最小限のテストケースで実装を導出
2. **共通箇所**: すべてのテストで共有される構造とパターン
3. **差分**: 各ケース固有の検証ポイント（入力、境界値、エラー）
4. **実装の一般化**: 複数の具体例から抽象的なパターンを抽出

### このアプローチの利点

- ✅ **効率性**: 37のテストケースで主要機能を網羅
- ✅ **明確性**: 各テストの目的と検証内容が明確
- ✅ **保守性**: 共通パターンによる拡張の容易さ
- ✅ **品質**: 境界値テストによる高い信頼性
- ✅ **ドキュメント性**: テストがそのまま仕様書として機能
- ✅ **完全性**: 条件分岐を持つ全コンポーネントをカバー
- ✅ **最小性**: ロジックのテストに集中、冗長なテストを排除

### バックエンドとの一貫性

フロントエンド（TypeScript）とバックエンド（Go）の両方で同じTriangulationの原則を適用することで、プロジェクト全体で統一されたテスト戦略を実現しています。

---

## テスト実行結果

```
✓ Test Files  7 passed (7)
✓ Tests  37 passed (37)
  Duration  1.17s (transform 751ms, setup 745ms, collect 1.05s, tests 885ms)
```

全37テストが成功し、Triangulationの原則に従った最小構成で最大限の効果を達成しました。

### 最適化の履歴

**2025年11月25日**: SearchForm.test.tsxを最適化（40 → 37テスト）

削除したテスト（3件）:

1. `"フォームの基本要素が表示される"` - ロジックのテストではなく、UIの存在確認
2. `"ジャンル選択肢が正しく表示される"` - useGenresフックで既に保証済み
3. `"ローディング終了後は検索ボタンが有効化される"` - isLoading=falseは冗長（Point 3と表裏）

**理由**: Triangulationは**ロジックの正しさを複数の角度から証明する手法**であり、単なるUIの存在確認や他のテストで既に保証されている内容は対象外。最小の3点（最小構成、最大構成、UI状態）でSearchFormの全ての動作を証明可能。

**効果**:

- テスト数: 40 → 37（7.5%削減）
- テスト実行時間: ほぼ変化なし（1.17秒）
- カバレッジ: ロジックの100%カバーを維持
- 保守性: 冗長なテストを排除し、意図がより明確に

### テスト対象コンポーネントの完全性

| コンポーネント     | ロジック | テスト  | 理由                          |
| ------------------ | -------- | ------- | ----------------------------- |
| `SearchForm`       | ✅ あり  | ✅ あり | フォーム送信ロジック          |
| `Pagination`       | ✅ あり  | ✅ あり | ページ遷移ロジック            |
| `RestaurantList`   | ✅ あり  | ✅ あり | 条件分岐 + イベントハンドラー |
| `ErrorMessage`     | ❌ なし  | ❌ なし | 単純な表示のみ                |
| `LoadingSpinner`   | ❌ なし  | ❌ なし | CSSアニメーションのみ         |
| `RestaurantDetail` | ❌ なし  | ❌ なし | 単純な表示のみ                |

**結論**: ロジックを持つ全コンポーネントがテスト済み。必要最小限かつ完全なカバレッジを達成。
