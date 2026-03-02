import { useState, useEffect, useCallback } from "react";
import type { Genre } from "@/types";
import { getGenres } from "@/api/genreApi";
import { ApiError } from "@/api/client";
interface UseGenresResult {
  genres: Genre[];
  isLoading: boolean;
  error: string | null;
  clearError: () => void;
}
export const useGenres = (): UseGenresResult => {
  const [genres, setGenres] = useState<Genre[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    const fetchGenres = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const fetchedGenres = await getGenres();
        setGenres(fetchedGenres);
      } catch (err) {
        if (err instanceof ApiError && err.status === 400) {
          setError("検索条件を変更してください");
        } else {
          setError("ジャンルの取得に失敗しました");
        }
        console.error("Failed to fetch genres:", err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchGenres();
  }, []);
  const clearError = useCallback(() => {
    setError(null);
  }, []);
  return { genres, isLoading, error, clearError };
};
