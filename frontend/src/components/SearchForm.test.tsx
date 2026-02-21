// ===================================================================
// SearchFormコンポーネントのテスト
// ===================================================================
// 【目的】SearchFormコンポーネントの動作を検証
// 【Triangulation】最小の3点でフォームの主要機能をテスト
// ===================================================================

// Vitestのテスト関数をインポート
import { describe, it, expect, vi, beforeEach } from "vitest";

// React Testing Library: Reactコンポーネントをテストするためのライブラリ
// - render: コンポーネントを仮想DOMにレンダリング
// - screen: レンダリングされた要素にアクセス
// - fireEvent: ユーザー操作（クリック、入力など）をシミュレート
import { render, screen, fireEvent } from "@testing-library/react";

// テスト対象のコンポーネントとフックをインポート
import { SearchForm } from "./SearchForm";
import * as useGenresHook from "@/hooks/useGenres";

// ===================================================================
// モックの設定
// ===================================================================

// useGenresフックをモック化（偽物に置き換え）
// 【理由】外部APIに依存しない独立したテストを実現
vi.mock("@/hooks/useGenres");

// モックジャンルデータ
const mockGenres = [
  { code: "G001", name: "居酒屋" },
  { code: "G002", name: "イタリアン" },
  { code: "G003", name: "中華" },
];

// ===================================================================
// Triangulation Point の概要
// ===================================================================
// Point 1: 必須項目のみで検索 - 最小構成の成功ケース
// Point 2: 全項目を入力して検索 - 最大構成の成功ケース
// Point 3: ローディング中のUI状態 - ボタン無効化
// ===================================================================

describe("SearchForm", () => {
  // モック関数の準備
  let mockOnSearch: ReturnType<typeof vi.fn>;

  // 各テストの前に実行される準備処理
  beforeEach(() => {
    // モック関数を作成（onSearchコールバックの偽物）
    mockOnSearch = vi.fn() as any;

    // useGenresフックのモック実装を設定
    // 【重要】spyOnで元のモジュールを監視し、戻り値を偽装
    vi.spyOn(useGenresHook, "useGenres").mockReturnValue({
      genres: mockGenres,
      isLoading: false,
      error: null,
    });
  });

  // ===================================================================
  // 検索機能のテスト（Triangulation）
  // ===================================================================
  // 【目的】フォーム送信時の動作を検証
  // 【Triangulation戦略】異なる入力パターンで動作を確認
  // ===================================================================

  describe("検索機能", () => {
    // Triangulation Point 1: 必須項目のみで検索
    // 【目的】都道府県のみ選択して検索が実行できることを確認
    // 【重要性】最小構成での動作保証（必須パラメータのみ）
    // 【差分】serviceAreaのみ → オプションパラメータは未定義
    it("都道府県のみ選択して検索できる", () => {
      // Arrange（準備）: コンポーネントをレンダリング
      render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

      // Act（実行）
      // 1. 都道府県を選択
      const selects = screen.getAllByRole("combobox");
      const serviceAreaSelect = selects[0]; // 1つ目は都道府県
      // SS10は存在しないので、実際の値（SA11=東京）を使用
      fireEvent.change(serviceAreaSelect, { target: { value: "SA11" } });

      // 2. 検索ボタンをクリック
      const searchButton = screen.getByRole("button", { name: "検索" });
      fireEvent.click(searchButton);

      // Assert（検証）
      // onSearch関数が正しいパラメータで呼ばれたか
      expect(mockOnSearch).toHaveBeenCalledWith({
        serviceArea: "SA11", // 必須項目（東京）
        address: undefined, // オプション項目は未定義
        genre: undefined, // オプション項目は未定義
        keyword: undefined, // オプション項目は未定義
      });

      // 呼び出し回数は1回
      expect(mockOnSearch).toHaveBeenCalledTimes(1);
    });

    // Triangulation Point 2: 全項目を入力して検索
    // 【目的】すべてのフィールドに値を入力して検索できることを確認
    // 【重要性】最大構成での動作保証（全パラメータ）
    // 【差分】serviceArea + address + genre + keyword → すべて指定
    it("全ての項目を入力して検索できる", () => {
      // Arrange（準備）: コンポーネントをレンダリング
      render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

      // Act（実行）
      // 1. 都道府県を選択（大阪）
      const selects = screen.getAllByRole("combobox");
      const serviceAreaSelect = selects[0]; // 1つ目は都道府県
      fireEvent.change(serviceAreaSelect, { target: { value: "SA23" } }); // SA23 = 大阪

      // 2. 住所を入力
      const addressInput = screen.getByPlaceholderText("例: 新宿");
      fireEvent.change(addressInput, { target: { value: "梅田" } });

      // 3. ジャンルを選択
      const genreSelect = selects[1]; // 2つ目はジャンル
      fireEvent.change(genreSelect, { target: { value: "G001" } });

      // 4. キーワードを入力
      const keywordInput = screen.getByPlaceholderText("例: 個室");
      fireEvent.change(keywordInput, { target: { value: "飲み放題" } });

      // 5. 検索ボタンをクリック
      const searchButton = screen.getByRole("button", { name: "検索" });
      fireEvent.click(searchButton);

      // Assert（検証）
      // onSearch関数が正しいパラメータで呼ばれたか
      expect(mockOnSearch).toHaveBeenCalledWith({
        serviceArea: "SA23", // 大阪
        address: "梅田",
        genre: "G001", // 居酒屋
        keyword: "飲み放題",
      });

      // 呼び出し回数は1回
      expect(mockOnSearch).toHaveBeenCalledTimes(1);
    });
  });

  // ===================================================================
  // ローディング状態のテスト
  // ===================================================================
  // 【目的】ローディング中の動作を検証
  // 【重要性】非同期処理中のUI状態確認
  // ===================================================================

  describe("ローディング状態", () => {
    // Triangulation Point 3: ローディング中のUI状態
    // 【目的】検索中はボタンが無効化され、テキストが変わることを確認
    // 【重要性】二重送信防止とユーザーへの適切なフィードバック
    // 【差分】isLoading=true → disabled=true, text="検索中..."
    it("ローディング中は検索ボタンが無効化される", () => {
      // Arrange（準備）: isLoading=trueでレンダリング
      render(<SearchForm onSearch={mockOnSearch as any} isLoading={true} />);

      // Assert（検証）
      // 検索ボタンが無効化されているか
      const searchButton = screen.getByRole("button", { name: "検索中..." });
      expect(searchButton).toBeDisabled();

      // ボタンのテキストが「検索中...」に変わっているか
      expect(searchButton).toHaveTextContent("検索中...");
    });
  });
});

// ===================================================================
// Triangulationのまとめ
// ===================================================================
// Point 1: 最小パラメータ → serviceAreaのみで検索成功
// Point 2: 最大パラメータ → 全項目入力で検索成功
// Point 3: ローディング状態 → ボタン無効化とテキスト変更
// ===================================================================
// この3つのテストで、SearchFormの主要な動作が網羅的にテストされる
// （Triangulationの原則：最小の3点で実装の正しさを証明）
//
// 削除したテスト（理由）：
// - 初期表示テスト：ロジックのテストではなく、UIの存在確認
// - ジャンル表示テスト：useGenresフックで既に保証されている
// - ローディング終了テスト：isLoadingのfalse側は冗長（Point 3と表裏）
// ===================================================================
