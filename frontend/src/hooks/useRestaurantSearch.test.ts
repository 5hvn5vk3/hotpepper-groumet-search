// ===================================================================
// useRestaurantSearch カスタムフックのテスト
// ===================================================================
// このファイルは、レストラン検索のロジックをカプセル化した
// Reactカスタムフック（useRestaurantSearch）の動作を検証するテストコードです。
//
// 【テストの目的】
// - 初期状態が正しく設定されるか
// - 検索機能（search）が正しく動作するか
// - ページ変更（changePage）が正しく動作するか
// - エラーハンドリングが適切に行われるか
// - 状態管理（useState）が正しく機能するか
// ===================================================================

// Vitestのテスト関数をインポート
import { describe, it, expect, vi, beforeEach } from "vitest";

// React Testing Library: Reactコンポーネントやフックをテストするためのライブラリ
// - renderHook: カスタムフックをテスト環境でレンダリングして実行
// - waitFor: 非同期の状態更新を待つためのユーティリティ
import { renderHook, waitFor } from "@testing-library/react";

// テスト対象のカスタムフック
import { useRestaurantSearch } from "./useRestaurantSearch";

// モック化するAPI関数
import { searchRestaurants } from "@/api/restaurantApi";

// 型定義
import type { GourmetSearchResponse } from "@/types";

// vi.mock(): searchRestaurants関数をモック（偽物）に置き換え
// 実際のAPIを呼ばず、テスト用の応答を返すようにする
vi.mock("@/api/restaurantApi");

// "useRestaurantSearch"という名前のテストグループ
describe("useRestaurantSearch", () => {
  // モックデータ: APIが返す偽のレストラン検索結果
  // 実際のAPIレスポンスと同じ形式のデータを用意
  const mockSearchResult: GourmetSearchResponse["results"] = {
    api_version: "1.0",
    results_available: 100, // 全体で100件の検索結果がある
    results_returned: "20", // このページには20件返す（文字列型）
    results_start: 1, // 1件目から開始
    shop: [
      // レストラン情報の配列（1件のサンプルデータ）
      {
        id: "J001234567", // 店舗ID
        name: "テストレストラン", // 店舗名
        address: "東京都渋谷区", // 住所
        lat: 35.6595, // 緯度
        lng: 139.7004, // 経度
        genre: { name: "居酒屋", catch: "気軽に楽しめる" }, // ジャンル情報
        catch: "美味しい居酒屋", // キャッチコピー
        access: "渋谷駅徒歩5分", // アクセス情報
        urls: { pc: "https://example.com" }, // 店舗URL
        photo: {
          // 店舗写真
          pc: {
            l: "https://example.com/photo_l.jpg", // 大サイズ
            m: "https://example.com/photo_m.jpg", // 中サイズ
            s: "https://example.com/photo_s.jpg", // 小サイズ
          },
        },
      },
    ],
  };

  // beforeEach: 各テストの前に実行される準備処理
  beforeEach(() => {
    // モック関数の呼び出し履歴をクリア
    vi.clearAllMocks();

    // searchRestaurants関数が呼ばれた時、mockSearchResultを返すように設定
    // mockResolvedValue(): Promiseが成功した時の戻り値を設定
    vi.mocked(searchRestaurants).mockResolvedValue({
      results: mockSearchResult,
    });
  });

  // ===================================================================
  // 初期状態のテスト: フックが正しく初期化されるか
  // ===================================================================
  describe("初期状態", () => {
    // Triangulation Point 1: デフォルト状態の確認
    // 【目的】フックの初期状態がすべて正しい値になっているか確認
    // 【重要性】アプリ起動時の状態が正しいことを保証
    it("初期状態が正しく設定される", () => {
      // Act（実行）: renderHookでカスタムフックを実行
      // itemsPerPage=20 として初期化
      // result.current: フックが返す現在の値にアクセス
      const { result } = renderHook(() => useRestaurantSearch(20));

      // Assert（検証）: すべての初期状態を確認

      // 検索結果はまだないのでnull
      expect(result.current.searchResult).toBeNull();

      // ローディング中ではない
      expect(result.current.isLoading).toBe(false);

      // まだ検索を実行していない
      expect(result.current.hasSearched).toBe(false);

      // エラーは発生していない
      expect(result.current.error).toBeNull();

      // 検索パラメータはまだ設定されていない
      expect(result.current.currentParams).toBeNull();

      // 現在のページは1ページ目
      expect(result.current.currentPage).toBe(1);
    });
  });

  // ===================================================================
  // searchメソッドのテスト: 検索機能の動作確認
  // ===================================================================
  describe("search", () => {
    // Triangulation Point 1: 基本的な検索
    // 【目的】最小限のパラメータで検索が成功することを確認
    // 【重要性】基本的な検索フローの動作保証
    it("検索を実行して結果を取得できる", async () => {
      // Arrange（準備）: フックを初期化
      const { result } = renderHook(() => useRestaurantSearch(20));

      // Act（実行）: 検索を実行（必須パラメータのみ）
      await result.current.search({ serviceArea: "SS10" });

      // Assert（検証）: waitForで非同期の状態更新を待つ
      // Reactの状態更新は非同期なので、waitForで完了を待つ必要がある
      await waitFor(() => {
        // 検索結果が正しく設定されている
        expect(result.current.searchResult).toEqual(mockSearchResult);

        // 検索実行済みフラグがtrueになっている
        expect(result.current.hasSearched).toBe(true);

        // ページが1にリセットされている（新規検索は常に1ページ目から）
        expect(result.current.currentPage).toBe(1);
      });
    });

    // Triangulation Point 2: 複数パラメータでの検索
    // 【目的】すべてのパラメータを含む検索が動作することを確認
    // 【重要性】最大構成での動作保証
    it("全パラメータを含む検索ができる", async () => {
      const { result } = renderHook(() => useRestaurantSearch(20));

      // Act（実行）: すべてのパラメータを含めて検索
      await result.current.search({
        serviceArea: "SS10",
        address: "新宿",
        genre: "G001",
        keyword: "居酒屋",
      });

      // Assert（検証）: すべてのパラメータが正しく渡されているか
      await waitFor(() => {
        expect(searchRestaurants).toHaveBeenCalledWith(
          expect.objectContaining({
            serviceArea: "SS10",
            address: "新宿",
            genre: "G001",
            keyword: "居酒屋",
            page: 1,
            count: 20,
          })
        );
      });
    });
  });

  // ===================================================================
  // changePageメソッドのテスト: ページ変更機能の動作確認
  // ===================================================================
  describe("changePage", () => {
    // Triangulation Point 1: 2ページ目への遷移
    // 【目的】ページ変更機能が基本的に動作することを確認
    // 【重要性】ページネーション機能の基本動作保証
    it("2ページ目に遷移できる", async () => {
      const { result } = renderHook(() => useRestaurantSearch(20));

      // Arrange（準備）: まず検索を実行してパラメータを設定
      await result.current.search({ serviceArea: "SS10" });
      await waitFor(() => expect(result.current.hasSearched).toBe(true));

      // Act（実行）: 2ページ目に遷移
      await result.current.changePage(2);

      // Assert（検証）
      await waitFor(() => {
        // 現在のページが2になっている
        expect(result.current.currentPage).toBe(2);

        // searchRestaurantsがpage=2で呼ばれている
        expect(searchRestaurants).toHaveBeenCalledWith(
          expect.objectContaining({ page: 2 })
        );
      });
    });

    // Triangulation Point 2: 3ページ目への遷移
    // 【目的】異なるページ番号でも動作することを確認
    // 【重要性】任意のページへの遷移が可能であることを保証
    it("3ページ目に遷移できる", async () => {
      const { result } = renderHook(() => useRestaurantSearch(20));

      // Arrange（準備）: 検索実行
      await result.current.search({ serviceArea: "SS10" });
      await waitFor(() => expect(result.current.hasSearched).toBe(true));

      // Act（実行）: 3ページ目に遷移
      await result.current.changePage(3);

      // Assert（検証）: 3ページ目になっている
      await waitFor(() => {
        expect(result.current.currentPage).toBe(3);
      });
    });

    // Triangulation Point 3: 同じページへの遷移は無視
    // 【目的】無駄なAPI呼び出しを防ぐ最適化が動作することを確認
    // 【重要性】パフォーマンス最適化、不要な通信の削減
    it("同じページ番号の場合は何もしない", async () => {
      const { result } = renderHook(() => useRestaurantSearch(20));

      // Arrange（準備）: 検索実行（1ページ目）
      await result.current.search({ serviceArea: "SS10" });
      await waitFor(() => expect(result.current.hasSearched).toBe(true));

      // モック関数の履歴をクリア（これ以降の呼び出しだけを記録）
      vi.clearAllMocks();

      // Act（実行）: 同じ1ページ目への遷移を試みる
      await result.current.changePage(1);

      // Assert（検証）: searchRestaurantsが呼ばれていないことを確認
      // 同じページなので、無駄なAPI呼び出しをしない
      expect(searchRestaurants).not.toHaveBeenCalled();
    });
  });

  // ===================================================================
  // エラーハンドリングのテスト: エラー処理の動作確認
  // ===================================================================
  describe("エラーハンドリング", () => {
    // Triangulation Point 1: Errorオブジェクトの場合
    // 【目的】標準的なErrorオブジェクトのエラーメッセージを正しく処理できるか確認
    // 【重要性】一般的なエラーケースの処理保証
    it("Errorオブジェクトのエラーメッセージを設定する", async () => {
      const { result } = renderHook(() => useRestaurantSearch(20));

      // Arrange（準備）: searchRestaurantsがエラーを返すように設定
      // mockRejectedValue(): Promiseが失敗した時の値を設定
      vi.mocked(searchRestaurants).mockRejectedValue(
        new Error("Network error") // 標準的なErrorオブジェクト
      );

      // Act（実行）: エラーが発生する検索を実行
      await result.current.search({ serviceArea: "SS10" });

      // Assert（検証）
      await waitFor(() => {
        // error.messageの値がエラー状態に設定されている
        expect(result.current.error).toBe("Network error");

        // ローディングが終了している
        expect(result.current.isLoading).toBe(false);
      });
    });

    // Triangulation Point 2: 文字列エラーの場合
    // 【目的】Errorオブジェクトでない（文字列などの）エラーも処理できるか確認
    // 【重要性】予期しない形式のエラーに対する防御的プログラミング
    it("文字列エラーの場合デフォルトメッセージを設定する", async () => {
      const { result } = renderHook(() => useRestaurantSearch(20));

      // Arrange（準備）: Errorオブジェクトでない文字列エラーを設定
      vi.mocked(searchRestaurants).mockRejectedValue("Unknown error");

      // Act（実行）
      await result.current.search({ serviceArea: "SS10" });

      // Assert（検証）
      await waitFor(() => {
        // Errorオブジェクトでない場合、デフォルトメッセージが使われる
        // 実装: err instanceof Error ? err.message : "検索に失敗しました"
        expect(result.current.error).toBe("検索に失敗しました");
      });
    });

    // Triangulation Point 3: エラークリア機能
    // 【目的】エラー状態を手動でクリアできることを確認
    // 【重要性】エラーメッセージを閉じる機能の実装
    it("clearErrorでエラーをクリアできる", async () => {
      const { result } = renderHook(() => useRestaurantSearch(20));

      // Arrange（準備）: エラーを発生させる
      vi.mocked(searchRestaurants).mockRejectedValue(
        new Error("Network error")
      );

      await result.current.search({ serviceArea: "SS10" });

      // エラーが設定されるまで待つ
      await waitFor(() => expect(result.current.error).toBeTruthy());

      // Act（実行）: clearError関数を呼んでエラーをクリア
      await waitFor(() => {
        result.current.clearError();
      });

      // Assert（検証）: エラーがnullにリセットされている
      await waitFor(() => {
        expect(result.current.error).toBeNull();
      });
    });
  });
});
