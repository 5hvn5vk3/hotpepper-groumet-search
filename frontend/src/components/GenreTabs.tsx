import React from "react";
import type { Genre } from "@/types";
interface GenreTabsProps {
  genres: Genre[];
  selectedGenre: string;
  canUseGenreTabs: boolean;
  onGenreTabClick: (genreCode: string) => void;
}
const splitGenreLabel = (label: string): string => label.replace(/・/g, "\n・");
const getTabClassName = (
  isSelected: boolean,
  canUseGenreTabs: boolean,
  index: number,
): string => {
  const activeColor =
    index % 2 === 0 ? "bg-red-600 text-white" : "bg-red-800 text-white";
  return (
    "shrink-0 md:shrink flex items-center justify-center rounded-md border px-3 py-2 text-[11px] font-bold leading-tight whitespace-pre-line text-left transition-colors " +
    (isSelected ? "bg-white text-red-600 border-red-600" : activeColor) +
    (!canUseGenreTabs ? " opacity-50 cursor-not-allowed" : "")
  );
};
export const GenreTabs: React.FC<GenreTabsProps> = ({
  genres,
  selectedGenre,
  canUseGenreTabs,
  onGenreTabClick,
}) => {
  return (
    <div>
      <div
        role="tablist"
        aria-label="ジャンルタブ"
        className="flex items-stretch overflow-x-auto pb-2 md:overflow-visible"
      >
        <button
          type="button"
          role="tab"
          aria-selected={selectedGenre === ""}
          aria-label="すべて"
          disabled={!canUseGenreTabs}
          onClick={() => onGenreTabClick("")}
          className={getTabClassName(selectedGenre === "", canUseGenreTabs, 0)}
        >
          すべて
        </button>

        {genres.map((genre, index) => (
          <button
            key={genre.code}
            type="button"
            role="tab"
            aria-selected={selectedGenre === genre.code}
            aria-label={genre.name}
            disabled={!canUseGenreTabs}
            onClick={() => onGenreTabClick(genre.code)}
            className={getTabClassName(
              selectedGenre === genre.code,
              canUseGenreTabs,
              index + 1,
            )}
          >
            {splitGenreLabel(genre.name)}
          </button>
        ))}
      </div>
    </div>
  );
};
