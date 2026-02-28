import { AppHeader, ErrorMessage, GenreTabs, LoadingSpinner, Pagination, RestaurantDetail, RestaurantList, SearchForm, ScrollToTopButton, } from "@/components";
import type { RestaurantDetailStatus, GourmetSearchParams, Shop } from "@/types";
import { useGenres } from "@/hooks/useGenres";
import { useRestaurantSearch } from "@/hooks/useRestaurantSearch";
import { useModal } from "@/hooks/useModal";
import { fetchRestaurantDetail } from "@/api/restaurantApi";
import type { SearchFormState } from "@/components/SearchForm";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { ITEMS_PER_PAGE } from "@/constants";
function App() {
    const { searchResult, isLoading, hasSearched, error, currentPage, search, changePage, clearError, } = useRestaurantSearch(ITEMS_PER_PAGE);
    const { isOpen, data: selectedRestaurant, open, close } = useModal<Shop>();
    const [userLat, setUserLat] = useState<number | null>(null);
    const [userLng, setUserLng] = useState<number | null>(null);
    const [selectedGenre, setSelectedGenre] = useState("");
    const [searchFormState, setSearchFormState] = useState<SearchFormState>({
        address: "",
        keyword: "",
        lat: null,
        lng: null,
        range: 3,
        hasStoredLocation: false,
    });
    const [detailRestaurant, setDetailRestaurant] = useState<Shop | null>(null);
    const [detailStatus, setDetailStatus] = useState<RestaurantDetailStatus>("idle");
    const [fixedControlsHeight, setFixedControlsHeight] = useState(0);
    const detailRequestIdRef = useRef(0);
    const fixedControlsRef = useRef<HTMLDivElement | null>(null);
    const { genres, isLoading: genresLoading, error: genresError, clearError: clearGenresError } = useGenres();
    const modalRestaurant = useMemo(() => {
        if (!selectedRestaurant) {
            return null;
        }
        if (!detailRestaurant) {
            return selectedRestaurant;
        }
        return {
            ...selectedRestaurant,
            ...detailRestaurant,
            credit_card: selectedRestaurant.credit_card,
        };
    }, [selectedRestaurant, detailRestaurant]);
    const handleLocationChange = useCallback((lat: number, lng: number) => {
        setUserLat(lat);
        setUserLng(lng);
    }, []);
    const handleSearch = (params: GourmetSearchParams) => {
        setSelectedGenre(params.genre ?? "");
        search(params);
    };
    const handleSearchFormStateChange = useCallback((nextState: SearchFormState) => {
        setSearchFormState((prevState) => {
            const isUnchanged = prevState.address === nextState.address &&
                prevState.keyword === nextState.keyword &&
                prevState.lat === nextState.lat &&
                prevState.lng === nextState.lng &&
                prevState.range === nextState.range &&
                prevState.hasStoredLocation === nextState.hasStoredLocation;
            return isUnchanged ? prevState : nextState;
        });
    }, []);
    const hasLocationContext = searchFormState.hasStoredLocation ||
        (searchFormState.lat !== null && searchFormState.lng !== null);
    const hasSearchText = searchFormState.address.trim() !== "" ||
        searchFormState.keyword.trim() !== "";
    const canUseGenreTabs = hasSearched &&
        !isLoading &&
        !genresLoading &&
        (hasLocationContext || hasSearchText);
    const handleGenreTabClick = (genreCode: string) => {
        if (!canUseGenreTabs || selectedGenre === genreCode) {
            return;
        }
        const trimmedAddress = searchFormState.address.trim();
        const trimmedKeyword = searchFormState.keyword.trim();
        const hasLocation = searchFormState.lat !== null && searchFormState.lng !== null;
        const hasAddress = trimmedAddress !== "";
        const hasKeyword = trimmedKeyword !== "";
        if (!hasLocation && !hasAddress && !hasKeyword) {
            return;
        }
        handleSearch({
            address: hasAddress ? trimmedAddress : undefined,
            genre: genreCode || undefined,
            keyword: hasKeyword ? trimmedKeyword : undefined,
            lat: hasLocation ? searchFormState.lat! : undefined,
            lng: hasLocation ? searchFormState.lng! : undefined,
            range: hasLocation ? searchFormState.range : undefined,
        });
    };
    const handlePageChange = (page: number) => {
        changePage(page);
    };
    const handleSelectRestaurant = async (restaurant: Shop) => {
        open(restaurant);
        setDetailRestaurant(null);
        setDetailStatus("loading");
        const requestId = ++detailRequestIdRef.current;
        try {
            const detail = await fetchRestaurantDetail(restaurant.id);
            if (detailRequestIdRef.current !== requestId) {
                return;
            }
            if (!detail) {
                setDetailStatus("failed");
                return;
            }
            setDetailRestaurant(detail);
            setDetailStatus("ready");
        }
        catch (detailError) {
            if (detailRequestIdRef.current === requestId) {
                setDetailStatus("failed");
            }
            console.error("Failed to fetch restaurant detail:", detailError);
        }
    };
    const handleCloseRestaurantDetail = () => {
        detailRequestIdRef.current += 1;
        setDetailRestaurant(null);
        setDetailStatus("idle");
        close();
    };
    const hasPagination = searchResult !== null && searchResult.results_available > 0;
    const reservedBottomSpace = hasPagination
        ? Math.max(fixedControlsHeight, 224)
        : 0;
    useEffect(() => {
        if (!hasPagination) {
            setFixedControlsHeight(0);
            return;
        }
        const controlsElement = fixedControlsRef.current;
        if (!controlsElement) {
            return;
        }
        const updateHeight = () => {
            setFixedControlsHeight(controlsElement.getBoundingClientRect().height);
        };
        updateHeight();
        if (typeof ResizeObserver === "undefined") {
            window.addEventListener("resize", updateHeight);
            return () => {
                window.removeEventListener("resize", updateHeight);
            };
        }
        const observer = new ResizeObserver(updateHeight);
        observer.observe(controlsElement);
        return () => {
            observer.disconnect();
        };
    }, [hasPagination]);
    return (<div className="min-h-screen bg-gray-100">
      
      <AppHeader />

      
      <main className="container mx-auto px-4 py-8" style={hasPagination
            ? { paddingBottom: `calc(2rem + ${reservedBottomSpace}px)` }
            : undefined}>
        
        <SearchForm onSearch={handleSearch} isLoading={isLoading} selectedGenre={selectedGenre} onLocationChange={handleLocationChange} onSearchStateChange={handleSearchFormStateChange}/>

        {hasSearched && (<div className="sticky top-0 z-20 mb-6 rounded-lg bg-gray-100/95 py-2 backdrop-blur-sm">
            <GenreTabs genres={genres} selectedGenre={selectedGenre} canUseGenreTabs={canUseGenreTabs} onGenreTabClick={handleGenreTabClick}/>
          </div>)}

        
        {error && <ErrorMessage message={error} onClose={clearError}/>}
        {genresError && <ErrorMessage message={genresError} onClose={clearGenresError}/>}

        
        {isLoading && <LoadingSpinner />}

        
        {!isLoading && hasSearched && (<RestaurantList restaurants={searchResult?.shop ?? []} onSelectRestaurant={handleSelectRestaurant} userLat={userLat ?? undefined} userLng={userLng ?? undefined}/>)}

        
      </main>

      
      {hasPagination && searchResult && (<div ref={fixedControlsRef} className="fixed inset-x-0 bottom-0 z-30 border-t border-gray-200 bg-white/95 px-4 pb-[calc(env(safe-area-inset-bottom)+0.5rem)] pt-3 backdrop-blur-sm">
          <div className="container mx-auto flex flex-col items-end gap-2">
            <ScrollToTopButton />
            <div className="w-full">
              <Pagination currentPage={currentPage} totalCount={searchResult.results_available} count={ITEMS_PER_PAGE} onPageChange={handlePageChange} disabled={isLoading} start={searchResult.results_start} available={searchResult.results_available}/>
            </div>
          </div>
        </div>)}

      
      {isOpen && modalRestaurant && (<RestaurantDetail restaurant={modalRestaurant} detailStatus={detailStatus} onClose={handleCloseRestaurantDetail} userLat={userLat ?? undefined} userLng={userLng ?? undefined}/>)}
    </div>);
}
export default App;
