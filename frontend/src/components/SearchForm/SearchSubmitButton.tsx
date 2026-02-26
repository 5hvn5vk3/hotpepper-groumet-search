import React from "react";

interface SearchSubmitButtonProps {
  isLoading: boolean;
  hasStoredLocation: boolean;
}

export const SearchSubmitButton: React.FC<SearchSubmitButtonProps> = ({
  isLoading,
  hasStoredLocation,
}) => {
  return (
    <button
      type="submit"
      disabled={isLoading}
      className="w-full bg-red-600 text-white py-2 px-4 rounded-md hover:bg-red-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
    >
      {isLoading
        ? "検索中..."
        : hasStoredLocation
          ? "検索（距離順）"
          : "検索（おススメ順）"}
    </button>
  );
};
