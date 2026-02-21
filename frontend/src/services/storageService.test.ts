// ===================================================================
// StorageServiceのテスト
// ===================================================================
// このファイルは、ブラウザのsessionStorageを管理するStorageServiceクラスの
// 動作を検証するテストコードです。
//
// 【テストの目的】
// - オブジェクトデータの保存（set）が正しく動作するか
// - オブジェクトデータの取得（get）が正しく動作するか
// - エラーケース（存在しないキー、無効なJSON）を適切に処理できるか
// ===================================================================

// Vitestのテスト関数をインポート
// - describe: テストをグループ化する関数（テストスイート）
// - it: 個別のテストケースを定義する関数
// - expect: 値を検証するためのアサーション関数
// - beforeEach: 各テストの前に実行される準備処理
// - vi: モック（偽物の関数）を作成するためのユーティリティ
import { describe, it, expect, beforeEach, vi } from "vitest";

// テスト対象のStorageServiceをインポート
import { storageService } from "./storageService";

// describe: "StorageService"という名前のテストグループを作成
// この中に複数のテストケースをまとめて記述します
describe("StorageService", () => {
  // beforeEach: 各テストケースの前に毎回実行される関数
  // テスト間で状態が混ざらないよう、毎回クリーンな状態にリセットします
  beforeEach(() => {
    // sessionStorage.clear(): ブラウザのsessionStorageを空にする
    // 前のテストで保存したデータが残らないようにする
    sessionStorage.clear();

    // vi.clearAllMocks(): モック関数の呼び出し履歴をクリア
    // 前のテストの影響を受けないようにする
    vi.clearAllMocks();
  });

  // ===================================================================
  // getメソッドのテスト: データ取得機能の検証
  // ===================================================================
  describe("get", () => {
    // 【目的】存在しないキーを指定した場合、nullが返されることを確認
    // 【重要性】データが存在しない場合の安全な処理を保証
    it("存在しないキーの場合nullを返す", () => {
      // Act（実行）: 存在しないキー"nonExistent"でデータを取得を試みる
      const result = storageService.get("nonExistent");

      // Assert（検証）: 結果がnullであることを確認
      // toBeNull(): 値が厳密にnullであることを検証
      expect(result).toBeNull();
    });

    // 【目的】オブジェクトデータの保存と取得が正しく動作することを確認
    // 【重要性】アプリケーションで使用する複合データ型の動作を保証
    it("オブジェクトデータを正しく取得できる", () => {
      // Arrange（準備）: テスト用のオブジェクトを作成
      const testData = { name: "test", value: 123 };

      // sessionStorageにオブジェクトをJSON文字列として保存
      sessionStorage.setItem("testKey", JSON.stringify(testData));

      // Act（実行）: storageServiceを使ってオブジェクトを取得
      // typeof testData: testDataの型を自動的に推論して使用
      const result = storageService.get<typeof testData>("testKey");

      // Assert（検証）: 取得したオブジェクトが元のオブジェクトと同じ構造・値を持つことを確認
      // toEqual(): オブジェクトの内容が等しいことを検証（深い比較）
      expect(result).toEqual(testData);
    });

    // エラーケース: 異常系のテスト
    // 【目的】不正なデータ（パースできないJSON）に対して安全に動作することを確認
    // 【重要性】予期しないデータがあってもアプリがクラッシュしないことを保証
    it("無効なJSONの場合nullを返す", () => {
      // Arrange（準備）: JSON形式でない不正な文字列を保存
      // "invalid json"はJSON.parse()でエラーになる
      sessionStorage.setItem("testKey", "invalid json");

      // Act（実行）: 不正なデータの取得を試みる
      const result = storageService.get("testKey");

      // Assert（検証）: エラーが発生してもnullが返されることを確認
      // エラーをキャッチしてnullを返す実装になっている
      expect(result).toBeNull();
    });
  });

  // ===================================================================
  // setメソッドのテスト: データ保存機能の検証
  // ===================================================================
  describe("set", () => {
    // 【目的】オブジェクトデータの保存が正しく動作することを確認
    // 【重要性】アプリケーションで使用する複合データ型の保存動作を保証
    it("オブジェクトデータを保存できる", () => {
      // Arrange（準備）: テスト用のオブジェクトを作成
      const testData = { name: "test", value: 123 };

      // Act（実行）: オブジェクトを保存
      storageService.set("testKey", testData);

      // Assert（検証）: オブジェクトがJSON文字列として正しく保存されているか確認
      const stored = sessionStorage.getItem("testKey");

      // JSON.stringify({ name: "test", value: 123 })
      // → '{"name":"test","value":123}'
      expect(stored).toBe(JSON.stringify(testData));
    });
  });
});
