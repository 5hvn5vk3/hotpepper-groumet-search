import React from "react";
import { calculateDistance, formatDistance } from "@/utils/distance";

interface DistanceSectionProps {
  userLat: number;
  userLng: number;
  restaurantLat: number;
  restaurantLng: number;
}

export const DistanceSection: React.FC<DistanceSectionProps> = ({
  userLat,
  userLng,
  restaurantLat,
  restaurantLng,
}) => {
  return (
    <div>
      <h4 className="font-semibold text-gray-700 mb-1">現在地からの距離</h4>
      <p className="text-red-600 font-medium">
        {formatDistance(
          calculateDistance(userLat, userLng, restaurantLat, restaurantLng),
        )}
      </p>
    </div>
  );
};
