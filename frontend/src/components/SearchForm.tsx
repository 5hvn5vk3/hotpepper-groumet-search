// ReactライブラリとuseStateフックをインポート
import React, { useState } from "react";
// 型定義をインポート
import type { GourmetSearchParams } from "@/types";
// ジャンルデータを取得するカスタムフックをインポート
import { useGenres } from "@/hooks/useGenres";

// SearchFormコンポーネントのProps（プロパティ）の型定義
interface SearchFormProps {
  onSearch: (params: GourmetSearchParams) => void; // 検索実行時のコールバック関数
  isLoading: boolean; // ローディング中かどうか
}

const EMPTY_SEARCH_MESSAGE =
  "位置情報が取得できないため検索できません。住所またはキーワード検索をお試しください";

// SearchFormコンポーネント：レストラン検索フォームを表示
export const SearchForm: React.FC<SearchFormProps> = ({
  onSearch,
  isLoading,
}) => {
  // 住所キーワードの状態（初期値：空文字列）
  const [address, setAddress] = useState("");
  // 選択されたジャンルコードの状態（初期値：空文字列）
  const [selectedGenre, setSelectedGenre] = useState("");
  // フリーワード検索の状態（初期値：空文字列）
  const [keyword, setKeyword] = useState("");
  // バリデーションメッセージの状態
  const [validationMessage, setValidationMessage] = useState("");

  // useGenresフックでジャンル一覧を取得
  const { genres, isLoading: genresLoading } = useGenres();

  // 現在地（Geolocation API）関連の状態
  const [lat, setLat] = useState<number | null>(null);
  const [lng, setLng] = useState<number | null>(null);
  const [range, setRange] = useState<number>(3); // デフォルトは3（1000m）

  // フォーム送信時のハンドラー関数
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const trimmedAddress = address.trim();
    const trimmedKeyword = keyword.trim();
    const hasLocation = lat !== null && lng !== null;
    const hasAddress = trimmedAddress !== "";
    const hasKeyword = trimmedKeyword !== "";

    if (!hasLocation && !hasAddress && !hasKeyword) {
      setValidationMessage(EMPTY_SEARCH_MESSAGE);
      return;
    }

    setValidationMessage("");
    onSearch({
      address: hasAddress ? trimmedAddress : undefined,
      genre: selectedGenre || undefined,
      keyword: hasKeyword ? trimmedKeyword : undefined,
      lat: hasLocation ? lat : undefined,
      lng: hasLocation ? lng : undefined,
      range: hasLocation ? range : undefined,
    });
  };

  const handleUseCurrentLocation = () => {
    if (!navigator.geolocation) {
      alert("Geolocation API がサポートされていません");
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        setLat(position.coords.latitude);
        setLng(position.coords.longitude);
        setValidationMessage("");
      },
      (err) => {
        console.error(err);
        alert("現在地の取得に失敗しました");
      },
      { enableHighAccuracy: true }
    );
  };

  const clearLocation = () => {
    setLat(null);
    setLng(null);
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="bg-white shadow-md rounded-lg p-6 mb-6"
    >
      <h2 className="text-2xl font-bold mb-4 text-gray-800">レストラン検索</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            現在地で検索
          </label>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={handleUseCurrentLocation}
              className="px-3 py-2 bg-gray-200 rounded-md hover:bg-gray-300"
            >
              現在地を取得
            </button>
            <button
              type="button"
              onClick={clearLocation}
              className="px-3 py-2 bg-gray-100 rounded-md hover:bg-gray-200"
            >
              クリア
            </button>
          </div>
          <div className="text-xs text-gray-500 mt-2">
            {lat !== null && lng !== null ? (
              <span>
                緯度: {lat.toFixed(5)}, 経度: {lng.toFixed(5)}
              </span>
            ) : (
              <span>現在地未取得</span>
            )}
          </div>
          <div className="mt-2">
            <label className="block text-sm text-gray-700 mb-1">検索範囲</label>
            <select
              value={range}
              onChange={(e) => setRange(Number(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value={1}>300m</option>
              <option value={2}>500m</option>
              <option value={3}>1000m</option>
              <option value={4}>2000m</option>
              <option value={5}>3000m</option>
            </select>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            住所（部分一致）
          </label>
          <input
            type="text"
            value={address}
            onChange={(e) => {
              setAddress(e.target.value);
              if (validationMessage) {
                setValidationMessage("");
              }
            }}
            placeholder="例: 新宿"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            ジャンル
          </label>
          <select
            value={selectedGenre}
            onChange={(e) => setSelectedGenre(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={genresLoading}
          >
            <option value="">すべて</option>
            {genres.map((genre) => (
              <option key={genre.code} value={genre.code}>
                {genre.name}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            キーワード
          </label>
          <input
            type="text"
            value={keyword}
            onChange={(e) => {
              setKeyword(e.target.value);
              if (validationMessage) {
                setValidationMessage("");
              }
            }}
            placeholder="例: 個室"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      {validationMessage && (
        <p className="mb-4 text-sm text-red-600" role="alert" aria-live="polite">
          {validationMessage}
        </p>
      )}

      <button
        type="submit"
        disabled={isLoading}
        className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
      >
        {isLoading ? "検索中..." : "検索"}
      </button>
    </form>
  );
};
