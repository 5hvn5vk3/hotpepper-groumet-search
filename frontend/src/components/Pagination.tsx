import React from "react";
interface PaginationProps {
    currentPage: number;
    totalCount: number;
    count: number;
    onPageChange: (page: number) => void;
    disabled?: boolean;
    start: number;
    available: number;
}
export const Pagination: React.FC<PaginationProps> = ({ currentPage, totalCount, count, onPageChange, disabled = false, start, available, }) => {
    if (totalCount <= 0) {
        return null;
    }
    const totalPages = Math.max(1, Math.ceil(totalCount / count));
    const clampedCurrentPage = Math.min(Math.max(currentPage, 1), totalPages);
    const canPrev = clampedCurrentPage > 1;
    const canNext = clampedCurrentPage < totalPages;
    const rangeStart = Math.min(Math.max(1, start), totalCount);
    const safeAvailable = available > 0 ? available : 1;
    const rangeEnd = Math.min(totalCount, Math.min(safeAvailable, rangeStart + count - 1));
    const handlePrev = () => {
        if (canPrev && !disabled) {
            onPageChange(clampedCurrentPage - 1);
        }
    };
    const handleNext = () => {
        if (canNext && !disabled) {
            onPageChange(clampedCurrentPage + 1);
        }
    };
    return (<div className="bg-white shadow-md rounded-lg p-4 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      
      <p className="text-sm text-gray-600">
        
        全{totalCount.toLocaleString()}件中 {rangeStart.toLocaleString()}〜
        {rangeEnd.toLocaleString()}件を表示
      </p>

      
      <div className="flex items-center gap-2">
        
        <button type="button" onClick={handlePrev} disabled={!canPrev || disabled} className="px-4 py-2 rounded-md border border-gray-300 text-gray-700 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed">
          前へ
        </button>

        
        <span className="text-sm font-medium text-gray-700">
          {clampedCurrentPage} / {totalPages}ページ
        </span>

        
        <button type="button" onClick={handleNext} disabled={!canNext || disabled} className="px-4 py-2 rounded-md border border-gray-300 text-gray-700 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed">
          次へ
        </button>
      </div>
    </div>);
};
