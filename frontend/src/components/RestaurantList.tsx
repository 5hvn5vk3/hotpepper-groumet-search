// Reactライブラリをインポート
import React from "react";
// Shop型の定義をインポート
import type { Shop } from "@/types";
// 距離計算ユーティリティをインポート
import { calculateDistance, formatDistance } from "@/utils/distance";

// RestaurantListコンポーネントのPropsの型定義
interface RestaurantListProps {
  restaurants: Shop[]; // 表示するレストランの配列
  onSelectRestaurant: (restaurant: Shop) => void; // レストランがクリックされた時のコールバック
  userLat?: number; // ユーザーの緯度（現在地）
  userLng?: number; // ユーザーの経度（現在地）
}

// RestaurantListコンポーネント：レストラン一覧をカード形式で表示
export const RestaurantList: React.FC<RestaurantListProps> = ({
  restaurants,
  onSelectRestaurant,
  userLat,
  userLng,
}) => {
  // レストランが0件の場合、「結果がありません」メッセージを表示
  if (restaurants.length === 0) {
    return (
      <div className="bg-white shadow-md rounded-lg p-8 text-center">
        <p className="text-gray-500">検索結果がありません</p>
      </div>
    );
  }

  // レストランが1件以上ある場合、カード形式で一覧表示
  return (
    // グリッドレイアウト：画面サイズに応じて列数を調整
    // モバイル: 1列、タブレット(md): 2列、デスクトップ(lg): 3列
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {/* restaurants配列をmapでループして各レストランのカードを生成 */}
      {restaurants.map((restaurant) => (
        // レストランカード
        <div
          key={restaurant.id} // Reactのリストレンダリングに必須のkey属性（一意のID）
          // onClick: カードがクリックされた時にonSelectRestaurant関数を呼び出す
          onClick={() => onSelectRestaurant(restaurant)}
          // className: Tailwind CSSでスタイリング
          // cursor-pointer: マウスカーソルをポインターに変更
          // hover:shadow-xl: ホバー時に影を強調
          className="bg-white shadow-md rounded-lg overflow-hidden cursor-pointer hover:shadow-xl transition-shadow"
        >
          {/* レストランの画像 */}
          <img
            src={restaurant.photo.pc.m} // 中サイズの画像URL
            alt={restaurant.name} // 代替テキスト（アクセシビリティ）
            // object-cover: 画像をコンテナに合わせてトリミング
            className="w-full h-48 object-cover"
          />
          {/* カードのコンテンツ部分 */}
          <div className="p-4">
            {/* レストラン名 */}
            <h3 className="text-lg font-bold text-gray-800 mb-2 truncate">
              {/* truncate: テキストが長い場合は省略記号(...)を表示 */}
              {restaurant.name}
            </h3>
            {/* ジャンル名 */}
            <p className="text-sm text-blue-600 mb-2">
              {restaurant.genre.name}
            </p>
            {/* キャッチコピー */}
            <p className="text-sm text-gray-600 mb-2 line-clamp-2">
              {/* line-clamp-2: テキストを2行までに制限し、それ以上は省略 */}
              {restaurant.catch}
            </p>
            {/* アクセス */}
            <p className="text-xs text-gray-500 truncate">
              アクセス：{restaurant.access}
            </p>
            {/* 現在地からの距離 */}
            {userLat !== undefined && userLng !== undefined && (
              <p className="text-xs text-green-600 font-medium mt-1">
                📍ここから
                {formatDistance(
                  calculateDistance(
                    userLat,
                    userLng,
                    restaurant.lat,
                    restaurant.lng,
                  ),
                )}
              </p>
            )}
          </div>
        </div>
      ))}
    </div>
  );
};
