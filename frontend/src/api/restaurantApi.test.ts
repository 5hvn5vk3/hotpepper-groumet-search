// ===================================================================
// RestaurantAPIのテスト
// ===================================================================
// このファイルは、レストラン検索API（searchRestaurants関数）の
// 動作を検証するテストコードです。
//
// 【テストの目的】
// - 検索パラメータが正しくURLクエリに変換されるか
// - ページネーション（ページ送り）の計算が正確か
// - 境界値（0以下、100超など）が適切にバリデーションされるか
// - APIの呼び出しが正しい形式で行われるか
// ===================================================================

// Vitestのテスト関数をインポート
import { describe, it, expect, vi, beforeEach } from "vitest";

// テスト対象の関数をインポート
import { searchRestaurants } from "./restaurantApi";

// モック（偽物）にするAPI呼び出し関数をインポート
import { apiGet } from "./client";

// 型定義をインポート
import type { GourmetSearchResponse, GourmetSearchParams } from "@/types";

// vi.mock(): 指定したモジュールをモック（偽物）に置き換える
// 実際のAPI呼び出しを行わず、テスト用の偽の応答を返すようにする
// 【理由】
// - 実際のAPIを呼ぶとテストが遅くなる
// - 外部サービスに依存せず、安定したテストができる
// - API呼び出しの形式だけを検証したい
vi.mock("./client");

// "searchRestaurants"という名前のテストグループ
describe("searchRestaurants", () => {
  // モックレスポンス: APIが返す偽のデータを定義
  // 実際のAPIレスポンスと同じ形式のデータを用意
  const mockResponse: GourmetSearchResponse = {
    results: {
      api_version: "1.0", // APIのバージョン
      results_available: 100, // 検索結果の総数
      results_returned: "20", // 今回返された結果の数（文字列型）
      results_start: 1, // 結果の開始位置
      shop: [], // レストランの配列（このテストでは空でOK）
    },
  };

  // beforeEach: 各テストの前に実行される準備処理
  beforeEach(() => {
    // モック関数の呼び出し履歴をクリア
    vi.clearAllMocks();

    // apiGet関数をモック化し、常にmockResponseを返すように設定
    // mockResolvedValue(): Promiseが成功した時の戻り値を設定
    // これにより、実際のHTTPリクエストを送らずにテストできる
    vi.mocked(apiGet).mockResolvedValue(mockResponse);
  });

  // ===================================================================
  // 基本的な検索パラメータのテスト
  // ===================================================================
  describe("基本的な検索パラメータ", () => {
    // Triangulation Point 1: serviceAreaが空の場合
    // 【目的】必須パラメータが未入力の場合の動作を確認
    // 【重要性】不正な入力に対する防御的プログラミング、早期リターン
    it("serviceAreaが空の場合はAPI呼び出しをスキップして空の結果を返す", async () => {
      // Arrange（準備）: serviceAreaが空文字列のパラメータ
      const params: GourmetSearchParams = {
        serviceArea: "",
      };

      // Act（実行）: 検索関数を実行
      const result = await searchRestaurants(params);

      // Assert（検証）: apiGetが呼ばれていないことを確認
      expect(apiGet).not.toHaveBeenCalled();

      // 空の結果が返されることを確認
      expect(result.results.results_available).toBe(0);
      expect(result.results.results_returned).toBe("0");
      expect(result.results.shop).toEqual([]);
    });

    // Triangulation Point 2: 必須パラメータのみで検索できるケース
    // 【目的】必須パラメータだけで検索が動作することを確認
    // 【重要性】最小限の正常構成での動作保証、デフォルト値が正しく適用されるか検証
    it("必須パラメータのみで検索できる", async () => {
      // Arrange（準備）: 必須パラメータのみのオブジェクトを作成
      // serviceArea: 都道府県コード（例: SS10は東京都）
      const params: GourmetSearchParams = {
        serviceArea: "SS10",
      };

      // Act（実行）
      await searchRestaurants(params);

      // Assert（検証）: 都道府県コードがURLに含まれているか
      expect(apiGet).toHaveBeenCalledWith(
        expect.stringContaining("service_area=SS10")
      );

      // デフォルト値: 開始位置が1（1ページ目の最初）
      expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("start=1"));

      // デフォルト値: 取得件数が20件
      expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("count=20"));
    });

    // Triangulation Point 2: 全てのパラメータを含むケース
    // 【目的】すべてのオプションパラメータが同時に使えることを確認
    // 【重要性】最大構成での動作保証、複数パラメータの組み合わせ検証
    it("全パラメータを含めて検索できる", async () => {
      // Arrange（準備）: すべてのパラメータを含むオブジェクト
      const params: GourmetSearchParams = {
        serviceArea: "SS10", // 必須: 都道府県
        address: "新宿", // オプション: 住所キーワード
        genre: "G001", // オプション: ジャンルコード
        keyword: "居酒屋", // オプション: フリーワード検索
      };

      // Act（実行）
      await searchRestaurants(params);

      // Assert（検証）: すべてのパラメータがURLに含まれているか確認
      expect(apiGet).toHaveBeenCalledWith(
        expect.stringContaining("service_area=SS10")
      );
      // "新宿" のURLエンコード
      expect(apiGet).toHaveBeenCalledWith(
        expect.stringContaining("address=%E6%96%B0%E5%AE%BF")
      );
      expect(apiGet).toHaveBeenCalledWith(
        expect.stringContaining("genre=G001")
      );
      // "居酒屋" のURLエンコード
      expect(apiGet).toHaveBeenCalledWith(
        expect.stringContaining("keyword=%E5%B1%85%E9%85%92%E5%B1%8B")
      );
    });
  });

  // ===================================================================
  // ページネーションのテスト: ページ番号とstart位置の変換ロジック
  // ===================================================================
  // 【計算式】 start = (page - 1) * count + 1
  describe("ページネーション", () => {
    // Triangulation Point 1: ページ1（デフォルト）
    // 【目的】1ページ目のstart位置が正しく計算されることを確認
    // 【計算】(1 - 1) * 10 + 1 = 1
    it("ページ1の場合start=1になる", async () => {
      const params: GourmetSearchParams = {
        serviceArea: "SS10",
        page: 1, // 1ページ目を指定
        count: 10, // 1ページ10件
      };

      await searchRestaurants(params);

      // 1ページ目の開始位置は1
      expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("start=1"));
    });

    // Triangulation Point 2: ページ2
    // 【目的】2ページ目のstart位置が正しく計算されることを確認
    // 【計算】(2 - 1) * 10 + 1 = 11
    // 【意味】1ページに10件あるので、2ページ目は11件目から始まる
    it("ページ2の場合start=11になる", async () => {
      const params: GourmetSearchParams = {
        serviceArea: "SS10",
        page: 2, // 2ページ目を指定
        count: 10, // 1ページ10件
      };

      await searchRestaurants(params);

      // 2ページ目の開始位置は11（1ページ10件 + 1）
      expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("start=11"));
    });

    // Triangulation Point 3: ページ3、異なるcount
    // 【目的】異なる件数設定でも計算が正しいことを確認
    // 【計算】(3 - 1) * 20 + 1 = 41
    // 【意味】1ページに20件あるので、3ページ目は41件目から始まる
    it("ページ3、count=20の場合start=41になる", async () => {
      const params: GourmetSearchParams = {
        serviceArea: "SS10",
        page: 3, // 3ページ目を指定
        count: 20, // 1ページ20件（前のテストと異なる）
      };

      await searchRestaurants(params);

      // 3ページ目、20件/ページの場合の開始位置は41
      // （1ページ目: 1-20、2ページ目: 21-40、3ページ目: 41-60）
      expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("start=41"));
    });
  });

  // ===================================================================
  // 境界値とバリデーションのテスト: 異常値の処理
  // ===================================================================
  describe("境界値とバリデーション", () => {
    // Triangulation Point 1: page < 1の場合
    // 【目的】ページ番号が0以下の場合、1として扱われることを確認
    // 【重要性】不正な入力に対する防御的プログラミング
    it("ページが0以下の場合1として扱う", async () => {
      const params: GourmetSearchParams = {
        serviceArea: "SS10",
        page: 0, // 不正な値: ページ0（通常は1から始まる）
      };

      await searchRestaurants(params);

      // Math.max(1, 0) = 1 により、ページ1として処理される
      expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("start=1"));
    });

    // Triangulation Point 2: count < 1の場合
    // 【目的】取得件数が0以下の場合、1として扱われることを確認
    // 【重要性】最低1件は取得する保証
    it("countが0以下の場合1として扱う", async () => {
      const params: GourmetSearchParams = {
        serviceArea: "SS10",
        count: 0, // 不正な値: 0件取得（意味がない）
      };

      await searchRestaurants(params);

      // Math.max(0, 1) = 1 により、最低1件は取得
      expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("count=1"));
    });

    // Triangulation Point 3: count > 100の場合
    // 【目的】取得件数が上限を超える場合、100に制限されることを確認
    // 【重要性】APIの負荷制限、一度に大量のデータを取得させない
    it("countが100を超える場合100として扱う", async () => {
      const params: GourmetSearchParams = {
        serviceArea: "SS10",
        count: 150, // 上限超過: 150件（APIの上限は100件）
      };

      await searchRestaurants(params);

      // Math.min(150, 100) = 100 により、上限100件に制限
      expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("count=100"));
    });
  });
});
