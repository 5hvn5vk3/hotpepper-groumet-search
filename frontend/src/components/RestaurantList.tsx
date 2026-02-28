import React from "react";
import type { ShopListItem } from "@/types";
import { calculateDistance, formatDistance } from "@/utils/distance";
interface RestaurantListProps {
    restaurants: ShopListItem[];
    onSelectRestaurant: (restaurant: ShopListItem) => void;
    userLat?: number;
    userLng?: number;
}
export const RestaurantList: React.FC<RestaurantListProps> = ({ restaurants, onSelectRestaurant, userLat, userLng, }) => {
    if (restaurants.length === 0) {
        return (<div className="bg-white shadow-md rounded-lg p-8 text-center">
        <p className="text-gray-500">検索結果がありません</p>
      </div>);
    }
    return (<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      
      {restaurants.map((restaurant) => {
            const creditCardNames = Array.from(new Set((restaurant.credit_card ?? [])
                .map((card) => card.name?.trim())
                .filter((name): name is string => Boolean(name))));
            return (<div key={restaurant.id} onClick={() => onSelectRestaurant(restaurant)} className="bg-white shadow-md rounded-lg overflow-hidden cursor-pointer hover:shadow-xl transition-shadow">
            
            <div className="flex">
              
              <img src={restaurant.photo.pc.l} alt={restaurant.name} className="w-30 h-30 object-cover flex-shrink-0"/>
              
              <div className="p-4 flex-1">
                
                <h3 className="text-lg font-bold text-gray-800 mb-2 truncate">
                  
                  {restaurant.name}
                </h3>
                <div className="text-sm text-gray-600 mb-2">
                  
                  <p className="text-sm text-red-600">
                    {restaurant.genre.name}
                  </p>
                  
                  <p className="text-xs text-gray-500 truncate">
                    {restaurant.genre.catch}
                  </p>
                  <p className="text-xs text-gray-500 truncate">
                    
                    {restaurant.catch}
                  </p>
                </div>
                
                <p className="text-sm text-gray-600 mb-2">
                  アクセス：{restaurant.access}
                </p>
                
                {userLat !== undefined && userLng !== undefined && (<p className="text-sm text-red-500 font-medium mt-1">
                    📍ここから
                    {formatDistance(calculateDistance(userLat, userLng, restaurant.lat, restaurant.lng))}
                  </p>)}
                
                <p className="text-xs text-gray-600 mt-2 line-clamp-2">
                  {creditCardNames.length > 0
                    ? creditCardNames.join(" / ")
                    : "カード情報なし"}
                </p>
              </div>
            </div>
            
          </div>);
        })}
    </div>);
};
