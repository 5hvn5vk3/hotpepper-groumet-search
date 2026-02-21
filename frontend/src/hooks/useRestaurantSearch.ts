// Reactのフックをインポート
// useState: 状態（データ）を管理するフック
// useCallback: 関数をメモ化（キャッシュ）して不要な再生成を防ぐフック
import { useState, useCallback } from "react";
// 型定義をインポート
import type { GourmetSearchResponse, GourmetSearchParams } from "@/types";
// レストラン検索API関数をインポート
import { searchRestaurants } from "@/api/restaurantApi";

// このフックが返す値の型定義
// インターフェースで構造を定義することで、型安全性を確保
interface UseRestaurantSearchResult {
  searchResult: GourmetSearchResponse["results"] | null; // 検索結果データ（なければnull）
  isLoading: boolean; // ローディング状態
  hasSearched: boolean; // 検索が実行されたか
  error: string | null; // エラーメッセージ（なければnull）
  currentParams: GourmetSearchParams | null; // 現在の検索パラメータ
  currentPage: number; // 現在のページ番号
  search: (params: GourmetSearchParams) => Promise<void>; // 検索実行関数
  changePage: (page: number) => Promise<void>; // ページ変更関数
  clearError: () => void; // エラークリア関数
}

// useRestaurantSearchフック：レストラン検索のロジックをカプセル化
// itemsPerPage: 1ページあたりの表示件数
// 戻り値: UseRestaurantSearchResult型のオブジェクト
export const useRestaurantSearch = (
  itemsPerPage: number
): UseRestaurantSearchResult => {
  // useState: 状態変数を宣言
  // [状態変数, 更新関数] = useState(初期値)の形式

  // 検索結果データの状態（初期値はnull）
  const [searchResult, setSearchResult] = useState<
    GourmetSearchResponse["results"] | null
  >(null);

  // ローディング状態（初期値はfalse）
  const [isLoading, setIsLoading] = useState(false);

  // 検索実行済みフラグ（初期値はfalse）
  const [hasSearched, setHasSearched] = useState(false);

  // エラーメッセージ（初期値はnull）
  const [error, setError] = useState<string | null>(null);

  // 現在の検索パラメータ（初期値はnull）
  const [currentParams, setCurrentParams] =
    useState<GourmetSearchParams | null>(null);

  // 現在のページ番号（初期値は1）
  const [currentPage, setCurrentPage] = useState(1);

  // fetchRestaurants: 実際にAPIを呼び出してレストランデータを取得する関数
  // useCallbackでメモ化
  const fetchRestaurants = useCallback(
    async (params: GourmetSearchParams, page: number) => {
      // ローディング開始
      setIsLoading(true);
      // 検索実行済みフラグをtrueに設定
      setHasSearched(true);
      // エラーをクリア
      setError(null);

      // try-catchでエラーハンドリング
      try {
        // searchRestaurants関数を呼び出してAPI通信
        // スプレッド構文（...params）で既存のパラメータを展開し、pageとcountを追加
        const response = await searchRestaurants({
          ...params, // 元のパラメータ（serviceArea, address, genre, keywordなど）
          page, // ページ番号を追加
          count: itemsPerPage, // 1ページの件数を追加
        });
        // 検索結果を状態にセット
        setSearchResult(response.results);
      } catch (err) {
        // エラーメッセージを作成
        // 三項演算子: 条件 ? 真の場合 : 偽の場合
        const errorMessage =
          err instanceof Error ? err.message : "検索に失敗しました";
        // エラーメッセージを状態にセット
        setError(errorMessage);
        // コンソールにエラーログを出力（開発者ツールで確認可能）
        console.error("Search failed:", err);
      } finally {
        // tryブロックとcatchブロックの後、必ず実行される
        // ローディング終了
        setIsLoading(false);
      }
    },
    []
  );

  // search: 新しい検索を実行する関数
  const search = useCallback(async (params: GourmetSearchParams) => {
    // 検索パラメータを状態に保存
    setCurrentParams(params);
    // ページ番号を1にリセット
    setCurrentPage(1);
    // 1ページ目のデータを取得
    await fetchRestaurants(params, 1);
  }, []);

  // changePage: ページを変更する関数
  const changePage = useCallback(
    async (page: number) => {
      // 検索パラメータがない場合、または同じページの場合は何もしない
      // ||は論理OR演算子：左右どちらかがtrueなら全体がtrue
      if (!currentParams || page === currentPage) {
        return;
      }
      // ページ番号を更新
      setCurrentPage(page);
      // 指定されたページのデータを取得
      await fetchRestaurants(currentParams, page);
    },
    [currentParams, currentPage] // これらの値が変わった場合のみ再生成
  );

  // clearError: エラーをクリアする関数
  const clearError = useCallback(() => {
    // エラーをnullに設定してクリア
    setError(null);
  }, []); // 依存配列が空：コンポーネントのマウント時に一度だけ作成

  // フックが返すオブジェクト：状態と関数をまとめて返す
  return {
    searchResult, // 検索結果データ
    isLoading, // ローディング状態
    hasSearched, // 検索実行済みフラグ
    error, // エラーメッセージ
    currentParams, // 現在の検索パラメータ
    currentPage, // 現在のページ番号
    search, // 検索実行関数
    changePage, // ページ変更関数
    clearError, // エラークリア関数
  };
};
