import React from "react";
import type { RestaurantDetailStatus, Shop } from "@/types";
import { calculateDistance, formatDistance } from "@/utils/distance";
interface RestaurantDetailProps {
    restaurant: Shop;
    detailStatus: RestaurantDetailStatus;
    onClose: () => void;
    userLat?: number;
    userLng?: number;
}
export const RestaurantDetail: React.FC<RestaurantDetailProps> = ({ restaurant, detailStatus, onClose, userLat, userLng, }) => {
    const mapQuery = encodeURIComponent(`${restaurant.name} ${restaurant.address}`);
    const googleMapsUrl = `https://www.google.com/maps/search/?api=1&query=${mapQuery}`;
    const creditCardNames = Array.from(new Set((restaurant.credit_card ?? [])
        .map((card) => card.name?.trim())
        .filter((name): name is string => Boolean(name))));
    const renderDeferredValue = (value?: string) => {
        if (detailStatus === "loading") {
            return "通信中";
        }
        if (detailStatus === "failed") {
            return "取得失敗";
        }
        return value?.trim() ? value : "情報なし";
    };
    return (<div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      
      
      
      
      <div className="bg-white rounded-lg max-w-3xl w-full max-h-[90vh] overflow-y-auto">
        
        
        <div className="sticky top-0 bg-white border-b px-6 py-4 flex justify-between items-center">
          
          <h2 className="text-2xl font-bold text-gray-800">店舗詳細</h2>
          
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700 text-2xl">
            × 
          </button>
        </div>

        
        <div className="p-6">
          
          <img src={restaurant.photo.pc.l} alt={restaurant.name} className="w-full h-64 object-cover rounded-lg mb-6"/>

          
          
          <div className="space-y-4">
            {detailStatus === "loading" && (<p className="text-sm text-blue-600">
                詳細情報を取得中です。未取得項目は「通信中」と表示しています。
              </p>)}

            
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

            
            {userLat !== undefined && userLng !== undefined && (<div>
                <h4 className="font-semibold text-gray-700 mb-1">
                  現在地からの距離
                </h4>
                <p className="text-red-600 font-medium">
                  {formatDistance(calculateDistance(userLat, userLng, restaurant.lat, restaurant.lng))}
                </p>
              </div>)}

            
            <div className="flex flex-wrap gap-4 mt-4">
              <div className="pt-4">
                
                <a href={googleMapsUrl} target="_blank" rel="noopener noreferrer" className="inline-block bg-green-600 text-white px-6 py-3 rounded-md hover:bg-yellow-700 transition-colors">
                  Googleマップで場所を見る
                </a>
              </div>

              <div className="pt-4">
                
                <a href={restaurant.urls.pc} target="_blank" rel="noopener noreferrer" className="inline-block bg-red-600 text-white px-6 py-3 rounded-md hover:bg-red-700 transition-colors">
                  ホットペッパーで詳細を見る
                </a>
              </div>
            </div>

            
            <div>
              <h4 className="font-semibold text-gray-700 mb-1">アクセス</h4>
              <p className="text-gray-600">{restaurant.access}</p>
            </div>

            
            <div>
              <h4 className="font-semibold text-gray-700 mb-1">住所</h4>
              <p className="text-gray-600">{restaurant.address}</p>
            </div>

            
            <div>
              <h4 className="font-semibold text-gray-700 mb-1">営業時間</h4>
              <p className="text-gray-600">
                {renderDeferredValue(restaurant.open)}
              </p>
            </div>

            <div>
              <h4 className="font-semibold text-gray-700 mb-1">定休日</h4>
              <p className="text-gray-600">
                {renderDeferredValue(restaurant.close)}
              </p>
            </div>

            {restaurant.budget_memo && (<div>
                <h4 className="font-semibold text-gray-700 mb-1">料金備考</h4>
                <p className="text-gray-600">{restaurant.budget_memo}</p>
              </div>)}

            
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

            <div>
              <h4 className="font-semibold text-gray-700 mb-1">設備・条件</h4>
              <p className="text-gray-600">
                Wi-Fi: {renderDeferredValue(restaurant.wifi)}
              </p>
              <p className="text-gray-600">
                個室: {renderDeferredValue(restaurant.private_room)}
              </p>
              <p className="text-gray-600">
                禁煙席: {renderDeferredValue(restaurant.non_smoking)}
              </p>
              <p className="text-gray-600">
                駐車場: {renderDeferredValue(restaurant.parking)}
              </p>
              <p className="text-gray-600">
                ランチ: {renderDeferredValue(restaurant.lunch)}
              </p>
              <p className="text-gray-600">
                深夜営業: {renderDeferredValue(restaurant.midnight)}
              </p>
            </div>

            {restaurant.shop_detail_memo && (<div>
                <h4 className="font-semibold text-gray-700 mb-1">備考</h4>
                <p className="text-gray-600">{restaurant.shop_detail_memo}</p>
              </div>)}
          </div>
        </div>
      </div>
    </div>);
};
