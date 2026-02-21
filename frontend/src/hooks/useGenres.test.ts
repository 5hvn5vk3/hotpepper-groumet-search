// ===================================================================
// useGenres カスタムフックのテスト
// ===================================================================
// 【目的】ジャンルマスタデータ取得フックの動作を検証
// 【Triangulation】最小の3点でフックの主要機能をテスト
// ===================================================================

import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { useGenres } from "./useGenres";
import { getGenres } from "@/api/genreApi";
import type { Genre } from "@/types";

// getGenres APIをモック化
vi.mock("@/api/genreApi");

describe("useGenres", () => {
  const mockGenres: Genre[] = [
    { code: "G001", name: "居酒屋" },
    { code: "G002", name: "イタリアン" },
    { code: "G003", name: "中華" },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  // ===================================================================
  // Triangulation Point 1: 成功時の基本動作
  // ===================================================================
  it("マウント時にジャンルデータを取得して状態を更新する", async () => {
    // Arrange（準備）: API成功時のモックを設定
    vi.mocked(getGenres).mockResolvedValue(mockGenres);

    // Act（実行）: フックをレンダリング
    const { result } = renderHook(() => useGenres());

    // Assert（検証）: 初期状態
    expect(result.current.isLoading).toBe(true);
    expect(result.current.genres).toEqual([]);
    expect(result.current.error).toBeNull();

    // データ取得完了を待つ
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.genres).toEqual(mockGenres);
      expect(result.current.error).toBeNull();
    });

    // APIが1回だけ呼ばれたことを確認
    expect(getGenres).toHaveBeenCalledTimes(1);
  });

  // ===================================================================
  // Triangulation Point 2: エラー時の動作
  // ===================================================================
  it("API呼び出し失敗時にエラーメッセージを設定する", async () => {
    // Arrange（準備）: API失敗時のモックを設定
    const errorMessage = "ネットワークエラー";
    vi.mocked(getGenres).mockRejectedValue(new Error(errorMessage));

    // Act（実行）: フックをレンダリング
    const { result } = renderHook(() => useGenres());

    // Assert（検証）: エラー状態になることを確認
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.genres).toEqual([]);
      expect(result.current.error).toBe(errorMessage);
    });
  });

  // ===================================================================
  // Triangulation Point 3: Error以外のエラー型
  // ===================================================================
  it("Errorオブジェクト以外のエラーの場合デフォルトメッセージを設定する", async () => {
    // Arrange（準備）: 文字列エラーを投げる
    vi.mocked(getGenres).mockRejectedValue("Unknown error");

    // Act（実行）: フックをレンダリング
    const { result } = renderHook(() => useGenres());

    // Assert（検証）: デフォルトメッセージが設定されることを確認
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBe("ジャンルの取得に失敗しました");
    });
  });
});
