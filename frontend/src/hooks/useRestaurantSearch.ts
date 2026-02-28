import { useState, useCallback, useRef } from "react";
import type { GourmetSearchResponse, GourmetSearchParams } from "@/types";
import { searchRestaurants } from "@/api/restaurantApi";
import { ApiError } from "@/api/client";
interface UseRestaurantSearchResult {
  searchResult: GourmetSearchResponse["results"] | null;
  isLoading: boolean;
  hasSearched: boolean;
  error: string | null;
  currentParams: GourmetSearchParams | null;
  currentPage: number;
  search: (params: GourmetSearchParams) => Promise<void>;
  changePage: (page: number) => Promise<void>;
  clearError: () => void;
}
export const useRestaurantSearch = (
  itemsPerPage: number,
): UseRestaurantSearchResult => {
  const [searchResult, setSearchResult] = useState<
    GourmetSearchResponse["results"] | null
  >(null);
  const [isLoading, setIsLoading] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentParams, setCurrentParams] =
    useState<GourmetSearchParams | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const latestRequestIdRef = useRef(0);
  const fetchRestaurants = useCallback(
    async (
      params: GourmetSearchParams,
      page: number,
      requestId: number,
    ): Promise<GourmetSearchResponse["results"] | null> => {
      setIsLoading(true);
      setHasSearched(true);
      setError(null);
      try {
        const response = await searchRestaurants({
          ...params,
          page,
          count: itemsPerPage,
        });
        if (latestRequestIdRef.current !== requestId) {
          return null;
        }
        return response.results;
      } catch (err) {
        if (latestRequestIdRef.current !== requestId) {
          return null;
        }
        if (err instanceof ApiError && err.status === 400) {
          setError("検索条件を変更してください");
        } else {
          setError("サービスが一時的に利用できません");
        }
        console.error("Search failed:", err);
        return null;
      } finally {
        if (latestRequestIdRef.current === requestId) {
          setIsLoading(false);
        }
      }
    },
    [itemsPerPage],
  );
  const search = useCallback(
    async (params: GourmetSearchParams) => {
      const requestId = ++latestRequestIdRef.current;
      setCurrentParams(params);
      const nextResults = await fetchRestaurants(params, 1, requestId);
      if (!nextResults) {
        return;
      }
      setSearchResult(nextResults);
      setCurrentPage(1);
    },
    [fetchRestaurants],
  );
  const changePage = useCallback(
    async (page: number) => {
      if (!currentParams || page === currentPage) {
        return;
      }
      const requestId = ++latestRequestIdRef.current;
      const nextResults = await fetchRestaurants(
        currentParams,
        page,
        requestId,
      );
      if (!nextResults) {
        return;
      }
      setSearchResult(nextResults);
      setCurrentPage(page);
    },
    [currentParams, currentPage, fetchRestaurants],
  );
  const clearError = useCallback(() => {
    setError(null);
  }, []);
  return {
    searchResult,
    isLoading,
    hasSearched,
    error,
    currentParams,
    currentPage,
    search,
    changePage,
    clearError,
  };
};
