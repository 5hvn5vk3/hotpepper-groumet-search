import { useState, useCallback } from "react";
import type { GourmetSearchResponse, GourmetSearchParams } from "@/types";
import { searchRestaurants } from "@/api/restaurantApi";
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
export const useRestaurantSearch = (itemsPerPage: number): UseRestaurantSearchResult => {
    const [searchResult, setSearchResult] = useState<GourmetSearchResponse["results"] | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [hasSearched, setHasSearched] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [currentParams, setCurrentParams] = useState<GourmetSearchParams | null>(null);
    const [currentPage, setCurrentPage] = useState(1);
    const fetchRestaurants = useCallback(async (params: GourmetSearchParams, page: number) => {
        setIsLoading(true);
        setHasSearched(true);
        setError(null);
        try {
            const response = await searchRestaurants({
                ...params,
                page,
                count: itemsPerPage,
            });
            setSearchResult(response.results);
        }
        catch (err) {
            const errorMessage = err instanceof Error ? err.message : "検索に失敗しました";
            setError(errorMessage);
            console.error("Search failed:", err);
        }
        finally {
            setIsLoading(false);
        }
    }, [itemsPerPage]);
    const search = useCallback(async (params: GourmetSearchParams) => {
        setCurrentParams(params);
        setCurrentPage(1);
        await fetchRestaurants(params, 1);
    }, [fetchRestaurants]);
    const changePage = useCallback(async (page: number) => {
        if (!currentParams || page === currentPage) {
            return;
        }
        setCurrentPage(page);
        await fetchRestaurants(currentParams, page);
    }, [currentParams, currentPage, fetchRestaurants]);
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
