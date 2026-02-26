// Reactライブラリをインポート
import React from "react";
// Shop型の定義をインポート
import type { Shop } from "@/types";
// 距離計算ユーティリティをインポート
import { calculateDistance, formatDistance } from "@/utils/distance";

// RestaurantDetailコンポーネントのPropsの型定義
interface RestaurantDetailProps {
  restaurant: Shop; // 表示するレストランのデータ
  onClose: () => void; // モーダルを閉じるためのコールバック関数
  userLat?: number; // ユーザーの緯度（現在地）
  userLng?: number; // ユーザーの経度（現在地）
}

// RestaurantDetailコンポーネント：レストランの詳細情報をモーダルで表示
export const RestaurantDetail: React.FC<RestaurantDetailProps> = ({
  restaurant,
  onClose,
  userLat,
  userLng,
}) => {
  const mapQuery = encodeURIComponent(
    `${restaurant.name} ${restaurant.address}`,
  );
  const googleMapsUrl = `https://www.google.com/maps/search/?api=1&query=${mapQuery}`;
  const creditCardNames = Array.from(
    new Set(
      (restaurant.credit_card ?? [])
        .map((card) => card.name?.trim())
        .filter((name): name is string => Boolean(name)),
    ),
  );

  return (
    // モーダルのオーバーレイ（背景）
    // fixed inset-0: 画面全体を覆う固定配置
    // bg-black bg-opacity-50: 半透明の黒背景
    // flex items-center justify-center: 中央揃え
    // z-50: 最前面に表示（z-indexが50）
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      {/* モーダルのコンテンツボックス */}
      {/* max-w-3xl: 最大幅を設定 */}
      {/* max-h-[90vh]: 最大高さを画面の90%に制限 */}
      {/* overflow-y-auto: 縦方向にスクロール可能 */}
      <div className="bg-white rounded-lg max-w-3xl w-full max-h-[90vh] overflow-y-auto">
        {/* ヘッダー部分（固定表示） */}
        {/* sticky top-0: スクロール時も上部に固定 */}
        <div className="sticky top-0 bg-white border-b px-6 py-4 flex justify-between items-center">
          {/* タイトル */}
          <h2 className="text-2xl font-bold text-gray-800">店舗詳細</h2>
          {/* 閉じるボタン */}
          <button
            onClick={onClose} // クリックでモーダルを閉じる
            className="text-gray-500 hover:text-gray-700 text-2xl"
          >
            × {/* × 記号（閉じるアイコン） */}
          </button>
        </div>

        {/* モーダルのメインコンテンツ */}
        <div className="p-6">
          {/* レストランの大きな画像 */}
          <img
            src={restaurant.photo.pc.l} // 大サイズの画像URL
            alt={restaurant.name}
            className="w-full h-64 object-cover rounded-lg mb-6"
          />

          {/* 詳細情報のセクション */}
          {/* space-y-4: 子要素間に縦方向のスペースを追加 */}
          <div className="space-y-4">
            {/* レストラン名とジャンル情報 */}
            <div>
              <h3 className="text-2xl font-bold text-gray-800 mb-2">
                {restaurant.name}
              </h3>
              <p className="text-red-600 font-medium">
                {restaurant.genre.name}
              </p>
              <p className="text-sm text-gray-600 mt-1">
                {restaurant.genre.catch}
              </p>
              <p className="text-sm text-gray-600 mt-1">{restaurant.catch}</p>
            </div>

            {/* 住所セクション */}
            <div>
              <h4 className="font-semibold text-gray-700 mb-1">住所</h4>
              <p className="text-gray-600">{restaurant.address}</p>
            </div>

            {/* 現在地からの距離セクション */}
            {userLat !== undefined && userLng !== undefined && (
              <div>
                <h4 className="font-semibold text-gray-700 mb-1">
                  現在地からの距離
                </h4>
                <p className="text-red-600 font-medium">
                  {formatDistance(
                    calculateDistance(
                      userLat,
                      userLng,
                      restaurant.lat,
                      restaurant.lng,
                    ),
                  )}
                </p>
              </div>
            )}

            {/* 営業時間 */}
            {restaurant.open && (
              <div>
                <h4 className="font-semibold text-gray-700 mb-1">営業時間</h4>
                <p className="text-gray-600">{restaurant.open}</p>
              </div>
            )}

            {/* 利用可能クレジットカード */}

            <div>
              <h4 className="font-semibold text-gray-700 mb-1">
                利用可能クレジットカード
              </h4>
              <p className="text-gray-600">
                {creditCardNames.length > 0
                  ? creditCardNames.join(" / ")
                  : "カード情報なし"}
              </p>
            </div>

            {/* アクセス情報セクション */}
            <div>
              <h4 className="font-semibold text-gray-700 mb-1">アクセス</h4>
              <p className="text-gray-600">{restaurant.access}</p>
            </div>

            {/* 外部リンクボタン */}
            <div className="pt-4">
              {/* a要素: Googleマップの詳細ページへのリンク */}
              <a
                href={googleMapsUrl} // Googleマップの詳細ページURL
                target="_blank" // 新しいタブで開く
                rel="noopener noreferrer" // セキュリティ対策（target="_blank"使用時の推奨設定）
                className="inline-block bg-green-600 text-white px-6 py-3 rounded-md hover:bg-yellow-700 transition-colors"
              >
                Googleマップで場所を見る
              </a>
            </div>

            <div className="pt-4">
              {/* a要素: ホットペッパーの詳細ページへのリンク */}
              <a
                href={restaurant.urls.pc} // レストランの詳細ページURL
                target="_blank" // 新しいタブで開く
                rel="noopener noreferrer" // セキュリティ対策（target="_blank"使用時の推奨設定）
                className="inline-block bg-red-600 text-white px-6 py-3 rounded-md hover:bg-red-700 transition-colors"
              >
                ホットペッパーで詳細を見る
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
