// Reactのフックをインポート
// useState: 状態管理
// useEffect: 副作用（データ取得など）を扱うフック
import { useState, useEffect } from "react";
// 型定義をインポート
import type { Genre } from "@/types";
// ジャンルAPI関数をインポート
import { getGenres } from "@/api/genreApi";

// このフックが返す値の型定義
interface UseGenresResult {
  genres: Genre[]; // ジャンル一覧
  isLoading: boolean; // ローディング状態
  error: string | null; // エラーメッセージ
}

// useGenresフック：ジャンルマスタデータを取得・管理するフック
export const useGenres = (): UseGenresResult => {
  // ジャンル一覧の状態（初期値：空配列）
  const [genres, setGenres] = useState<Genre[]>([]);

  // ローディング状態（初期値：false）
  const [isLoading, setIsLoading] = useState(false);

  // エラーメッセージ（初期値：null）
  const [error, setError] = useState<string | null>(null);

  // useEffect：コンポーネントのマウント時にAPIからジャンルデータを取得
  // 第2引数の依存配列が空[]なので、マウント時に一度だけ実行される
  useEffect(() => {
    // 非同期関数を定義（useEffectのコールバック自体はasyncにできないため）
    const fetchGenres = async () => {
      // ローディング開始
      setIsLoading(true);
      // エラーをクリア
      setError(null);

      try {
        // getGenres関数を使ってジャンルデータを取得
        // await: Promiseの完了を待つ
        const fetchedGenres = await getGenres();
        // 取得したジャンルデータを状態にセット
        setGenres(fetchedGenres);
      } catch (err) {
        // エラーが発生した場合の処理
        // エラーメッセージを作成（Error型の場合はそのメッセージ、それ以外はデフォルトメッセージ）
        const errorMessage =
          err instanceof Error ? err.message : "ジャンルの取得に失敗しました";
        // エラーメッセージを状態にセット
        setError(errorMessage);
        // コンソールにエラーログを出力
        console.error("Failed to fetch genres:", err);
      } finally {
        // try/catchの後、必ず実行される処理
        // ローディング終了
        setIsLoading(false);
      }
    };

    // 定義した非同期関数を実行
    fetchGenres();
  }, []); // 依存配列が空：マウント時に一度だけ実行

  // フックの戻り値：取得したジャンルデータと状態
  return { genres, isLoading, error };
};
