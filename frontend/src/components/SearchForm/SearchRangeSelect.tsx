import React from "react";

interface SearchRangeSelectProps {
  range: number;
  hasStoredLocation: boolean;
  onChange: (range: number) => void;
}

export const SearchRangeSelect: React.FC<SearchRangeSelectProps> = ({
  range,
  hasStoredLocation,
  onChange,
}) => {
  return (
    <div>
      <div className="mt-2">
        <label className="block text-sm text-gray-700 mb-1">
          検索範囲 （現在地から）
        </label>
        <select
          value={range}
          onChange={(e) => onChange(Number(e.target.value))}
          className={
            "w-full px-3 py-2 border border-gray-300 rounded-md" +
            (hasStoredLocation ? "" : " bg-gray-200")
          }
          disabled={!hasStoredLocation}
        >
          {hasStoredLocation ? (
            <>
              <option value={1}>300m</option>
              <option value={2}>500m</option>
              <option value={3}>1000m</option>
              <option value={4}>2000m</option>
              <option value={5}>3000m</option>
            </>
          ) : (
            <option>位置情報がありません</option>
          )}
        </select>
      </div>
    </div>
  );
};
