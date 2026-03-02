import React, { useEffect, useState } from "react";
interface ScrollToTopButtonProps {
  showAfter?: number;
}
export const ScrollToTopButton: React.FC<ScrollToTopButtonProps> = ({
  showAfter = 300,
}) => {
  const [isVisible, setIsVisible] = useState(false);
  useEffect(() => {
    const handleScroll = () => {
      setIsVisible(window.scrollY > showAfter);
    };
    handleScroll();
    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, [showAfter]);
  const handleClick = () => {
    window.scrollTo({ top: 0, behavior: "smooth" });
  };
  if (!isVisible) {
    return null;
  }
  return (
    <div className="flex flex-col items-center">
      <button
        type="button"
        onClick={handleClick}
        aria-label="ページ上部へ戻る"
        className="h-12 w-12 rounded-full bg-gray-800 text-white shadow-lg transition hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
      >
        <span aria-hidden="true" className="text-xl leading-none">
          ↑
        </span>
      </button>
      <p className="text-xs text-gray-600 mt-1">Back to Top</p>
    </div>
  );
};
