// Reactライブラリをインポート
import React from "react";

// PaginationコンポーネントのPropsの型定義
interface PaginationProps {
  currentPage: number; // 現在のページ番号
  totalCount: number; // 総件数
  count: number; // 1ページあたりの表示件数
  onPageChange: (page: number) => void; // ページ変更時のコールバック関数
  disabled?: boolean; // 無効化フラグ（オプション）
  start: number; // 表示開始位置
  available: number; // 総利用可能件数
}

// Paginationコンポーネント：ページネーション（ページ切り替え）UIを表示
export const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalCount,
  count,
  onPageChange,
  disabled = false, // デフォルト値: false
  start,
  available,
}) => {
  // 総件数が0以下の場合、ページネーションを表示しない
  // null: 何も表示しない
  if (totalCount <= 0) {
    return null;
  }

  // 総ページ数を計算
  // Math.ceil(): 切り上げ（例：20.1 → 21）
  // Math.max(1, ...): 最小値を1に設定（0ページにならないように）
  const totalPages = Math.max(1, Math.ceil(totalCount / count));

  // 現在のページ番号を有効な範囲内に制限
  // Math.min(x, max): xとmaxのうち小さい方
  // Math.max(x, min): xとminのうち大きい方
  const clampedCurrentPage = Math.min(Math.max(currentPage, 1), totalPages);

  // 前のページに戻れるかどうか（1ページ目より大きい場合true）
  const canPrev = clampedCurrentPage > 1;

  // 次のページに進めるかどうか（最終ページより小さい場合true）
  const canNext = clampedCurrentPage < totalPages;

  // 表示範囲の開始位置を計算（1以上、totalCount以下）
  const rangeStart = Math.min(Math.max(1, start), totalCount);

  // 利用可能件数を安全に取得（0以下の場合は1を使用）
  const safeAvailable = available > 0 ? available : 1;

  // 表示範囲の終了位置を計算
  // rangeStart + count - 1: 終了位置（1ページあたりの件数を使用）
  // Math.min(..., totalCount): 総件数を超えないように制限
  // Math.min(..., safeAvailable): 利用可能件数も超えないように制限
  const rangeEnd = Math.min(
    totalCount,
    Math.min(safeAvailable, rangeStart + count - 1)
  );

  // 「前へ」ボタンがクリックされたときのハンドラー
  const handlePrev = () => {
    // 前のページに戻れる場合、かつ無効化されていない場合のみ実行
    // &&演算子: 両方がtrueの場合のみ右側を評価
    if (canPrev && !disabled) {
      // 現在のページから1を引いた値でページ変更を通知
      onPageChange(clampedCurrentPage - 1);
    }
  };

  // 「次へ」ボタンがクリックされたときのハンドラー
  const handleNext = () => {
    // 次のページに進める場合、かつ無効化されていない場合のみ実行
    if (canNext && !disabled) {
      // 現在のページに1を足した値でページ変更を通知
      onPageChange(clampedCurrentPage + 1);
    }
  };

  // JSXを返す：ページネーションUIの構造
  return (
    // ページネーションコンテナ
    // flex-col: モバイルでは縦並び
    // md:flex-row: タブレット以上では横並び
    <div className="bg-white shadow-md rounded-lg p-4 mt-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      {/* 表示範囲の情報テキスト */}
      <p className="text-sm text-gray-600">
        {/* toLocaleString(): 数値を地域の形式でフォーマット（例：1000 → 1,000） */}
        全{totalCount.toLocaleString()}件中 {rangeStart.toLocaleString()}〜
        {rangeEnd.toLocaleString()}件を表示
      </p>

      {/* ページ切り替えボタンとページ表示 */}
      <div className="flex items-center gap-2">
        {/* 「前へ」ボタン */}
        <button
          type="button"
          onClick={handlePrev}
          // disabled属性: 前のページがない場合、または全体が無効化されている場合
          // !canPrev: canPrevがfalseの場合true（論理否定演算子）
          disabled={!canPrev || disabled}
          // disabled:で始まるクラス: ボタンが無効な場合のスタイル
          className="px-4 py-2 rounded-md border border-gray-300 text-gray-700 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          前へ
        </button>

        {/* 現在のページ番号 / 総ページ数 */}
        <span className="text-sm font-medium text-gray-700">
          {clampedCurrentPage} / {totalPages}ページ
        </span>

        {/* 「次へ」ボタン */}
        <button
          type="button"
          onClick={handleNext}
          // 次のページがない場合、または全体が無効化されている場合は無効
          disabled={!canNext || disabled}
          className="px-4 py-2 rounded-md border border-gray-300 text-gray-700 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          次へ
        </button>
      </div>
    </div>
  );
};
