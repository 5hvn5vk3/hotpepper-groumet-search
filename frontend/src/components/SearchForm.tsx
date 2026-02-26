// Reactライブラリと主要フックをインポート
import React, { useCallback, useEffect, useState } from "react";
// 型定義をインポート
import type { GourmetSearchParams } from "@/types";
import { storageService } from "@/services/storageService";
import { SearchTextField } from "./SearchTextField";

// SearchFormコンポーネントのProps（プロパティ）の型定義
interface SearchFormProps {
  onSearch: (params: GourmetSearchParams) => void; // 検索実行時のコールバック関数
  isLoading: boolean; // ローディング中かどうか
  selectedGenre?: string; // 現在選択中のジャンル
  onLocationChange?: (lat: number, lng: number) => void; // 現在地が変化したときのコールバック
  onSearchStateChange?: (state: SearchFormState) => void; // フォーム状態が変化したときのコールバック
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

// SearchFormコンポーネント：レストラン検索フォームを表示
export const SearchForm: React.FC<SearchFormProps> = ({
  onSearch,
  isLoading,
  selectedGenre = "",
  onLocationChange,
  onSearchStateChange,
}) => {
  // 住所キーワードの状態（初期値：空文字列）
  const [address, setAddress] = useState("");
  // フリーワード検索の状態（初期値：空文字列）
  const [keyword, setKeyword] = useState("");
  // バリデーションメッセージの状態
  const [validationMessage, setValidationMessage] = useState("");

  // 現在地（Geolocation API）関連の状態
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
  const [range, setRange] = useState<number>(3); // デフォルトは3（1000m）
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

  // フォーム送信時のハンドラー関数
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const locationFromStorage = getStoredLocation();

    if (!navigator.geolocation) {
      executeSearch({
        addressValue: address,
        keywordValue: keyword,
        locationOverride: locationFromStorage,
      });
      return;
    }

    const hasExpiredSavedLocation =
      locationFromStorage !== null &&
      Date.now() - locationFromStorage.timestamp >=
        LOCATION_REFRESH_INTERVAL_MS;
    const supportsPermissionQuery =
      typeof navigator.permissions?.query === "function";

    if (hasExpiredSavedLocation && supportsPermissionQuery) {
      void navigator.permissions
        .query({ name: "geolocation" })
        .then((permissionStatus) => {
          if (permissionStatus.state === "granted") {
            void requestCurrentLocation(true).then((nextLocation) => {
              executeSearch({
                addressValue: address,
                keywordValue: keyword,
                locationOverride: nextLocation ?? locationFromStorage,
              });
            });
            return;
          }
          executeSearch({
            addressValue: address,
            keywordValue: keyword,
            locationOverride: locationFromStorage,
          });
        })
        .catch((err) => {
          console.error(err);
          executeSearch({
            addressValue: address,
            keywordValue: keyword,
            locationOverride: locationFromStorage,
          });
        });
      return;
    }

    executeSearch({
      addressValue: address,
      keywordValue: keyword,
      locationOverride: locationFromStorage,
    });
  };

  // ページアクセス時に現在地取得を試行
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
      <h2 className="text-2xl font-bold mb-4 text-gray-800">レストラン検索</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <div>
          <div className="mt-2">
            <label className="block text-sm text-gray-700 mb-1">
              検索範囲 （現在地から）
            </label>
            <select
              value={range}
              onChange={(e) => setRange(Number(e.target.value))}
              className={
                "w-full px-3 py-2 border border-gray-300 rounded-md" +
                (hasStoredLocation ? "" : " bg-gray-200")
              }
              disabled={!hasStoredLocation}
            >
              {hasStoredLocation ? (
                <>
                  <option value={1}>300m</option>
                  <option value={2}>500m</option>
                  <option value={3}>1000m</option>
                  <option value={4}>2000m</option>
                  <option value={5}>3000m</option>
                </>
              ) : (
                <option>位置情報がありません</option>
              )}
            </select>
          </div>
        </div>

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
    </form>
  );
};
