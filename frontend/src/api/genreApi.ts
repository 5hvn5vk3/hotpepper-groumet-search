import type { GenreMasterResponse, Genre } from "@/types";
import { apiGet } from "@/api/client";
import { storageService } from "@/services/storageService";
const GENRES_CACHE_KEY = "genres";
export const getGenres = async (): Promise<Genre[]> => {
    const cached = storageService.get<Genre[]>(GENRES_CACHE_KEY);
    if (cached) {
        return cached;
    }
    const response = await apiGet<GenreMasterResponse>("/api/genre");
    const genres = response.results.genre;
    storageService.set(GENRES_CACHE_KEY, genres);
    return genres;
};
