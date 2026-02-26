import React from "react";
import type { RestaurantDetailStatus, Shop } from "@/types";
import { ActionLinks } from "./ActionLinks";
import { DeferredSection } from "./DeferredSection";
import { DetailSection } from "./DetailSection";
import { DistanceSection } from "./DistanceSection";
import { FacilitiesSection } from "./FacilitiesSection";
interface RestaurantDetailProps {
  restaurant: Shop;
  detailStatus: RestaurantDetailStatus;
  onClose: () => void;
  userLat?: number;
  userLng?: number;
}
export const RestaurantDetail: React.FC<RestaurantDetailProps> = ({
  restaurant,
  detailStatus,
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
  const facilityItems = [
    { label: "Wi-Fi", value: restaurant.wifi },
    { label: "個室", value: restaurant.private_room },
    { label: "禁煙席", value: restaurant.non_smoking },
    { label: "駐車場", value: restaurant.parking },
    { label: "ランチ", value: restaurant.lunch },
    { label: "深夜営業", value: restaurant.midnight },
  ];
  const renderDeferredValue = (value?: string) => {
    if (detailStatus === "loading") {
      return "通信中";
    }
    if (detailStatus === "failed") {
      return "取得失敗";
    }
    return value?.trim() ? value : "情報なし";
  };
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg max-w-3xl w-full max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-white border-b px-6 py-4 flex justify-between items-center">
          <h2 className="text-2xl font-bold text-gray-800">店舗詳細</h2>

          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-2xl"
          >
            ×
          </button>
        </div>

        <div className="p-6">
          <img
            src={restaurant.photo.pc.l}
            alt={restaurant.name}
            className="w-full h-64 object-cover rounded-lg mb-6"
          />

          <div className="space-y-4">
            {detailStatus === "loading" && (
              <p className="text-sm text-blue-600">
                詳細情報を取得中です。未取得項目は「通信中」と表示しています。
              </p>
            )}

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

            {userLat !== undefined && userLng !== undefined && (
              <div>
                <DistanceSection
                  userLat={userLat}
                  userLng={userLng}
                  restaurantLat={restaurant.lat}
                  restaurantLng={restaurant.lng}
                />
              </div>
            )}

            <ActionLinks
              googleMapsUrl={googleMapsUrl}
              hotpepperUrl={restaurant.urls.pc}
            />

            <DetailSection title="アクセス" value={restaurant.access} />

            <DetailSection title="住所" value={restaurant.address} />

            <DeferredSection
              title="営業時間"
              value={restaurant.open}
              renderDeferredValue={renderDeferredValue}
            />

            <DeferredSection
              title="定休日"
              value={restaurant.close}
              renderDeferredValue={renderDeferredValue}
            />

            {restaurant.budget_memo && (
              <DetailSection title="料金備考" value={restaurant.budget_memo} />
            )}

            <DetailSection
              title="利用可能クレジットカード"
              value={
                creditCardNames.length > 0
                  ? creditCardNames.join(" / ")
                  : "カード情報なし"
              }
            />

            <FacilitiesSection
              items={facilityItems}
              renderDeferredValue={renderDeferredValue}
            />

            {restaurant.shop_detail_memo && (
              <DetailSection title="備考" value={restaurant.shop_detail_memo} />
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
