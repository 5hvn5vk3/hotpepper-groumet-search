import React, { useCallback, useEffect, useRef, useState } from "react";
import type { GourmetSearchParams } from "@/types";
import { storageService } from "@/services/storageService";
import { SearchTextField } from "./SearchTextField";
import { SearchRangeSelect } from "./SearchRangeSelect";
import { SearchSubmitButton } from "./SearchSubmitButton";
interface SearchFormProps {
  onSearch: (params: GourmetSearchParams) => void;
  isLoading: boolean;
  selectedGenre?: string;
  onLocationChange?: (lat: number, lng: number) => void;
  onSearchStateChange?: (state: SearchFormState) => void;
}
const EMPTY_SEARCH_MESSAGE =
  "位置情報が取得できないため検索できません。住所またはキーワード検索をお試しください";
const LOCATION_STORAGE_KEY = "searchCurrentLocation";
const LOCATION_REFRESH_INTERVAL_MS = 3 * 60 * 1000;
interface StoredLocation {
  lat: number;
  lng: number;
  timestamp: number;
}
export interface SearchFormState {
  address: string;
  keyword: string;
  lat: number | null;
  lng: number | null;
  range: number;
  hasStoredLocation: boolean;
}
const isStoredLocation = (value: unknown): value is StoredLocation => {
  if (!value || typeof value !== "object") {
    return false;
  }
  const candidate = value as Partial<StoredLocation>;
  return (
    Number.isFinite(candidate.lat) &&
    Number.isFinite(candidate.lng) &&
    Number.isFinite(candidate.timestamp)
  );
};
const isReloadNavigation = (): boolean => {
  if (typeof performance === "undefined") {
    return false;
  }
  const navigationEntries = performance.getEntriesByType("navigation");
  const firstNavigationEntry = navigationEntries[0] as
    | PerformanceNavigationTiming
    | undefined;
  return firstNavigationEntry?.type === "reload";
};
export const SearchForm: React.FC<SearchFormProps> = ({
  onSearch,
  isLoading,
  selectedGenre = "",
  onLocationChange,
  onSearchStateChange,
}) => {
  const isLocationRefreshInProgressRef = useRef(false);
  const [address, setAddress] = useState("");
  const [keyword, setKeyword] = useState("");
  const [validationMessage, setValidationMessage] = useState("");
  const [lat, setLat] = useState<number | null>(() => {
    const savedLocation =
      storageService.get<StoredLocation>(LOCATION_STORAGE_KEY);
    return isStoredLocation(savedLocation) ? savedLocation.lat : null;
  });
  const [lng, setLng] = useState<number | null>(() => {
    const savedLocation =
      storageService.get<StoredLocation>(LOCATION_STORAGE_KEY);
    return isStoredLocation(savedLocation) ? savedLocation.lng : null;
  });
  const [range, setRange] = useState<number>(3);
  const [hasStoredLocation, setHasStoredLocation] = useState<boolean>(() => {
    const savedLocation =
      storageService.get<StoredLocation>(LOCATION_STORAGE_KEY);
    return isStoredLocation(savedLocation);
  });
  const getStoredLocation = useCallback((): StoredLocation | null => {
    const savedLocation =
      storageService.get<StoredLocation>(LOCATION_STORAGE_KEY);
    return isStoredLocation(savedLocation) ? savedLocation : null;
  }, []);
  const requestCurrentLocation = useCallback(
    (showErrorAlert: boolean): Promise<StoredLocation | null> => {
      if (!navigator.geolocation) {
        if (showErrorAlert) {
          alert("Geolocation API がサポートされていません");
        }
        return Promise.resolve(null);
      }
      return new Promise((resolve) => {
        navigator.geolocation.getCurrentPosition(
          (position) => {
            const nextLocation: StoredLocation = {
              lat: position.coords.latitude,
              lng: position.coords.longitude,
              timestamp: Date.now(),
            };
            setLat(nextLocation.lat);
            setLng(nextLocation.lng);
            storageService.set(LOCATION_STORAGE_KEY, nextLocation);
            setHasStoredLocation(true);
            setValidationMessage("");
            onLocationChange?.(nextLocation.lat, nextLocation.lng);
            resolve(nextLocation);
          },
          (err) => {
            console.error(err);
            if (showErrorAlert) {
              alert("現在地の取得に失敗しました");
            }
            resolve(null);
          },
          { enableHighAccuracy: true },
        );
      });
    },
    [onLocationChange],
  );
  interface ExecuteSearchOptions {
    addressValue: string;
    keywordValue: string;
    locationOverride?: StoredLocation | null;
  }
  const executeSearch = useCallback(
    ({
      addressValue,
      keywordValue,
      locationOverride,
    }: ExecuteSearchOptions): boolean => {
      const trimmedAddress = addressValue.trim();
      const trimmedKeyword = keywordValue.trim();
      const effectiveLat = locationOverride?.lat ?? lat;
      const effectiveLng = locationOverride?.lng ?? lng;
      const hasLocation = effectiveLat !== null && effectiveLng !== null;
      const hasAddress = trimmedAddress !== "";
      const hasKeyword = trimmedKeyword !== "";
      if (!hasLocation && !hasAddress && !hasKeyword) {
        setValidationMessage(EMPTY_SEARCH_MESSAGE);
        return false;
      }
      setValidationMessage("");
      onSearch({
        address: hasAddress ? trimmedAddress : undefined,
        genre: selectedGenre || undefined,
        keyword: hasKeyword ? trimmedKeyword : undefined,
        lat: hasLocation ? effectiveLat : undefined,
        lng: hasLocation ? effectiveLng : undefined,
        range: hasLocation ? range : undefined,
      });
      return true;
    },
    [lat, lng, onSearch, range, selectedGenre],
  );
  const refreshLocationIfExpired = useCallback(
    async (showErrorAlert: boolean): Promise<StoredLocation | null> => {
      const locationFromStorage = getStoredLocation();
      if (!navigator.geolocation) {
        return locationFromStorage;
      }
      const hasExpiredSavedLocation =
        locationFromStorage !== null &&
        Date.now() - locationFromStorage.timestamp >=
          LOCATION_REFRESH_INTERVAL_MS;
      const supportsPermissionQuery =
        typeof navigator.permissions?.query === "function";
      if (!hasExpiredSavedLocation || !supportsPermissionQuery) {
        return locationFromStorage;
      }
      if (isLocationRefreshInProgressRef.current) {
        return locationFromStorage;
      }
      isLocationRefreshInProgressRef.current = true;
      try {
        const permissionStatus = await navigator.permissions.query({
          name: "geolocation",
        });
        if (permissionStatus.state !== "granted") {
          return locationFromStorage;
        }
        const nextLocation = await requestCurrentLocation(showErrorAlert);
        return nextLocation ?? locationFromStorage;
      } catch (err) {
        console.error(err);
        return locationFromStorage;
      } finally {
        isLocationRefreshInProgressRef.current = false;
      }
    },
    [getStoredLocation, requestCurrentLocation],
  );
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    void refreshLocationIfExpired(true).then((locationForSearch) => {
      executeSearch({
        addressValue: address,
        keywordValue: keyword,
        locationOverride: locationForSearch,
      });
    });
  };
  useEffect(() => {
    const savedLocation =
      storageService.get<StoredLocation>(LOCATION_STORAGE_KEY);
    const hasSavedLocation = isStoredLocation(savedLocation);
    setHasStoredLocation(hasSavedLocation);
    if (hasSavedLocation) {
      setLat(savedLocation.lat);
      setLng(savedLocation.lng);
      onLocationChange?.(savedLocation.lat, savedLocation.lng);
    }
    const shouldRefreshLocation =
      !hasSavedLocation ||
      (isReloadNavigation() &&
        Date.now() - savedLocation.timestamp >= LOCATION_REFRESH_INTERVAL_MS);
    if (shouldRefreshLocation) {
      void requestCurrentLocation(false);
    }
  }, [onLocationChange, requestCurrentLocation]);
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (document.visibilityState !== "visible") {
        return;
      }
      void refreshLocationIfExpired(false);
    };

    document.addEventListener("visibilitychange", handleVisibilityChange);
    return () => {
      document.removeEventListener("visibilitychange", handleVisibilityChange);
    };
  }, [refreshLocationIfExpired]);
  useEffect(() => {
    onSearchStateChange?.({
      address,
      keyword,
      lat,
      lng,
      range,
      hasStoredLocation,
    });
  }, [
    address,
    keyword,
    lat,
    lng,
    range,
    hasStoredLocation,
    onSearchStateChange,
  ]);
  return (
    <form
      onSubmit={handleSubmit}
      className="bg-white shadow-md rounded-lg p-6 mb-6"
    >
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <SearchRangeSelect
          range={range}
          hasStoredLocation={hasStoredLocation}
          onChange={setRange}
        />

        <SearchTextField
          label="住所（部分一致）"
          placeholder="例: 新宿"
          value={address}
          onChange={(value) => {
            setAddress(value);
            if (validationMessage) {
              setValidationMessage("");
            }
          }}
        />

        <SearchTextField
          label="キーワード"
          placeholder="例: 個室"
          value={keyword}
          onChange={(value) => {
            setKeyword(value);
            if (validationMessage) {
              setValidationMessage("");
            }
          }}
        />
      </div>

      {validationMessage && (
        <p
          className="mb-4 text-sm text-red-600"
          role="alert"
          aria-live="polite"
        >
          {validationMessage}
        </p>
      )}

      <SearchSubmitButton
        isLoading={isLoading}
        hasStoredLocation={hasStoredLocation}
      />
    </form>
  );
};
