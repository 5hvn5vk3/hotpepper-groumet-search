// Reactライブラリをインポート
import React from "react";

// ErrorMessageコンポーネントのPropsの型定義
interface ErrorMessageProps {
  message: string; // 表示するエラーメッセージ
  onClose: () => void; // 閉じるボタンがクリックされた時のコールバック関数
}

// ErrorMessageコンポーネント：エラーメッセージを表示するバナー
export const ErrorMessage: React.FC<ErrorMessageProps> = ({
  message,
  onClose,
}) => {
  return (
    // エラーメッセージのコンテナ
    // bg-red-50: 薄い赤色の背景
    // border-red-200: 赤いボーダー
    // text-red-800: 濃い赤色のテキスト
    <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg mb-6 flex items-center justify-between">
      {/* 左側：アイコンとメッセージ */}
      <div className="flex items-center">
        {/* エラーアイコン（SVG） */}
        {/* SVGはスケーラブル・ベクター・グラフィックス：拡大縮小しても綺麗な画像形式 */}
        <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
          {/* pathタグでアイコンの形状を定義 */}
          {/* fillRule, clipRule: SVGの塗りつぶしルール */}
          <path
            fillRule="evenodd"
            // dプロパティ: SVGパスの座標情報（×マークの形を描画）
            d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
            clipRule="evenodd"
          />
        </svg>
        {/* エラーメッセージテキスト */}
        <span>{message}</span>
      </div>

      {/* 右側：閉じるボタン */}
      <button
        onClick={onClose} // クリックでonClose関数を実行
        className="text-red-800 hover:text-red-900 font-bold"
        aria-label="閉じる" // アクセシビリティ用のラベル（スクリーンリーダーで読まれる）
      >
        × {/* × 記号（閉じるアイコン） */}
      </button>
    </div>
  );
};
