// Reactライブラリをインポート
import React from "react";

// LoadingSpinnerコンポーネント：ローディング中の表示
// React.FCはFunction Componentの型（Propsなし）
export const LoadingSpinner: React.FC = () => {
  return (
    // コンテナ：中央揃えで上下にパディング
    <div className="text-center py-12">
      {/* スピナー（回転するローディングアイコン） */}
      {/* inline-block: インライン要素として表示 */}
      {/* animate-spin: Tailwind CSSのアニメーション（無限回転） */}
      {/* rounded-full: 完全な円形 */}
      {/* border-b-2: 下側のみ太いボーダー（これが回転することでスピナーに見える） */}
      <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>

      {/* ローディングメッセージ */}
      {/* mt-4: 上マージン（スピナーとの間隔） */}
      <p className="mt-4 text-gray-600">検索中...</p>
    </div>
  );
};
