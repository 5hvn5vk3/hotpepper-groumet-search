// ===================================================================
// Paginationコンポーネントのテスト
// ===================================================================
// 【目的】Paginationコンポーネントの動作を検証
// 【Triangulation】最小の3点でページネーションの主要機能をテスト
// ===================================================================

// Vitestのテスト関数をインポート
import { describe, it, expect, vi, beforeEach } from "vitest";

// React Testing Library
import { render, screen, fireEvent } from "@testing-library/react";

// テスト対象のコンポーネント
import { Pagination } from "./Pagination";

// ===================================================================
// Triangulation Point の概要
// ===================================================================
// Point 1: 初期ページ（1ページ目） - 「前へ」ボタンが無効
// Point 2: 中間ページ（2ページ目） - 両方のボタンが有効
// Point 3: 最終ページ（3ページ目） - 「次へ」ボタンが無効
// ===================================================================

describe("Pagination", () => {
  // モック関数の準備
  let mockOnPageChange: ReturnType<typeof vi.fn>;

  // 各テストの前に実行される準備処理
  beforeEach(() => {
    // ページ変更コールバックのモック関数を作成
    mockOnPageChange = vi.fn() as any;
  });

  // ===================================================================
  // 基本的なページネーション動作のテスト（Triangulation）
  // ===================================================================

  describe("ページネーション動作", () => {
    // Triangulation Point 1: 初期ページ（1ページ目）
    // 【目的】1ページ目では「前へ」ボタンが無効化されることを確認
    // 【重要性】境界値（最小ページ）での動作保証
    // 【差分】currentPage=1 → 「前へ」無効、「次へ」有効
    it("1ページ目では前へボタンが無効化される", () => {
      // Arrange（準備）: 1ページ目でレンダリング
      // totalCount=60, count=20 → 全3ページ
      render(
        <Pagination
          currentPage={1}
          totalCount={60}
          count={20}
          onPageChange={mockOnPageChange as any}
          start={1}
          available={60}
        />
      );

      // Assert（検証）
      // 「前へ」ボタンが無効化されているか
      const prevButton = screen.getByRole("button", { name: "前へ" });
      expect(prevButton).toBeDisabled();

      // 「次へ」ボタンは有効か
      const nextButton = screen.getByRole("button", { name: "次へ" });
      expect(nextButton).not.toBeDisabled();

      // ページ表示が正しいか
      expect(screen.getByText("1 / 3ページ")).toBeInTheDocument();

      // 表示範囲が正しいか
      expect(screen.getByText(/全60件中 1〜20件を表示/)).toBeInTheDocument();
    });

    // Triangulation Point 2: 中間ページ（2ページ目）
    // 【目的】中間ページでは両方のボタンが有効であることを確認
    // 【重要性】通常のページ遷移の動作保証
    // 【差分】currentPage=2 → 両方有効
    it("2ページ目では両方のボタンが有効", () => {
      // Arrange（準備）: 2ページ目でレンダリング
      render(
        <Pagination
          currentPage={2}
          totalCount={60}
          count={20}
          onPageChange={mockOnPageChange as any}
          start={21}
          available={60}
        />
      );

      // Assert（検証）
      // 「前へ」ボタンが有効か
      const prevButton = screen.getByRole("button", { name: "前へ" });
      expect(prevButton).not.toBeDisabled();

      // 「次へ」ボタンが有効か
      const nextButton = screen.getByRole("button", { name: "次へ" });
      expect(nextButton).not.toBeDisabled();

      // ページ表示が正しいか
      expect(screen.getByText("2 / 3ページ")).toBeInTheDocument();

      // 表示範囲が正しいか
      expect(screen.getByText(/全60件中 21〜40件を表示/)).toBeInTheDocument();
    });

    // Triangulation Point 3: 最終ページ（3ページ目）
    // 【目的】最終ページでは「次へ」ボタンが無効化されることを確認
    // 【重要性】境界値（最大ページ）での動作保証
    // 【差分】currentPage=3 → 「前へ」有効、「次へ」無効
    it("最終ページでは次へボタンが無効化される", () => {
      // Arrange（準備）: 3ページ目（最終ページ）でレンダリング
      render(
        <Pagination
          currentPage={3}
          totalCount={60}
          count={20}
          onPageChange={mockOnPageChange as any}
          start={41}
          available={60}
        />
      );

      // Assert（検証）
      // 「前へ」ボタンは有効か
      const prevButton = screen.getByRole("button", { name: "前へ" });
      expect(prevButton).not.toBeDisabled();

      // 「次へ」ボタンが無効化されているか
      const nextButton = screen.getByRole("button", { name: "次へ" });
      expect(nextButton).toBeDisabled();

      // ページ表示が正しいか
      expect(screen.getByText("3 / 3ページ")).toBeInTheDocument();

      // 表示範囲が正しいか（最後は60件まで）
      expect(screen.getByText(/全60件中 41〜60件を表示/)).toBeInTheDocument();
    });
  });

  // ===================================================================
  // ボタンクリック時の動作テスト
  // ===================================================================

  describe("ボタンクリック動作", () => {
    // Triangulation Point 4: 「次へ」ボタンのクリック
    // 【目的】「次へ」ボタンで次のページに遷移することを確認
    // 【重要性】基本的なナビゲーション動作
    it("次へボタンをクリックすると次のページに遷移", () => {
      // Arrange（準備）: 1ページ目でレンダリング
      render(
        <Pagination
          currentPage={1}
          totalCount={60}
          count={20}
          onPageChange={mockOnPageChange as any}
          start={1}
          available={60}
        />
      );

      // Act（実行）: 「次へ」ボタンをクリック
      const nextButton = screen.getByRole("button", { name: "次へ" });
      fireEvent.click(nextButton);

      // Assert（検証）
      // onPageChangeが2で呼ばれたか
      expect(mockOnPageChange).toHaveBeenCalledWith(2);
      expect(mockOnPageChange).toHaveBeenCalledTimes(1);
    });
  });

  // ===================================================================
  // エッジケースのテスト
  // ===================================================================

  describe("エッジケース", () => {
    // Triangulation Point 6: 総件数が0の場合
    // 【目的】検索結果が0件の場合、ページネーションが表示されないことを確認
    // 【重要性】空の結果に対する適切な表示
    it("総件数が0の場合は何も表示しない", () => {
      // Arrange（準備）: totalCount=0でレンダリング
      const { container } = render(
        <Pagination
          currentPage={1}
          totalCount={0}
          count={20}
          onPageChange={mockOnPageChange as any}
          start={1}
          available={0}
        />
      );

      // Assert（検証）
      // 何も表示されていないことを確認（nullを返す）
      expect(container.firstChild).toBeNull();
    });

    // Triangulation Point 7: disabledプロパティ
    // 【目的】disabled=trueの場合、ボタンが無効化されることを確認
    // 【重要性】ローディング中などの操作防止
    it("disabledプロパティがtrueの場合、ボタンが無効化される", () => {
      // Arrange（準備）: disabled=trueでレンダリング
      render(
        <Pagination
          currentPage={2}
          totalCount={60}
          count={20}
          onPageChange={mockOnPageChange as any}
          disabled={true}
          start={21}
          available={60}
        />
      );

      // Assert（検証）
      // 両方のボタンが無効化されているか
      const prevButton = screen.getByRole("button", { name: "前へ" });
      const nextButton = screen.getByRole("button", { name: "次へ" });

      expect(prevButton).toBeDisabled();
      expect(nextButton).toBeDisabled();

      // Act（実行）: ボタンをクリックしてみる
      fireEvent.click(prevButton);
      fireEvent.click(nextButton);

      // Assert（検証）: onPageChangeが呼ばれていないことを確認
      expect(mockOnPageChange).not.toHaveBeenCalled();
    });
  });
});

// ===================================================================
// Triangulationのまとめ
// ===================================================================
// Point 1: 1ページ目 → 「前へ」無効
// Point 2: 2ページ目 → 両方有効
// Point 3: 3ページ目 → 「次へ」無効
// Point 4: 「次へ」クリック → ページ+1
// Point 5: 「前へ」クリック → ページ-1
// Point 6: 総件数0 → 非表示
// Point 7: disabled → 全ボタン無効
// ===================================================================
// この7つのテストで、Paginationの主要な動作が網羅的にテストされる
// ===================================================================
