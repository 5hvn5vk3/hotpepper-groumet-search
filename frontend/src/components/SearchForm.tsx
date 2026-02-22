// ReactライブラリとuseStateフックをインポート
import React, { useState } from "react";
// 型定義をインポート
import type { GourmetSearchParams } from "@/types";
// 都道府県コードの定数をインポート
import { SERVICE_AREAS } from "@/constants";
// ジャンルデータを取得するカスタムフックをインポート
import { useGenres } from "@/hooks/useGenres";

// SearchFormコンポーネントのProps（プロパティ）の型定義
interface SearchFormProps {
  onSearch: (params: GourmetSearchParams) => void; // 検索実行時のコールバック関数
  isLoading: boolean; // ローディング中かどうか
}

// SearchFormコンポーネント：レストラン検索フォームを表示
// React.FC<Props>はFunction Componentの型定義
// 分割代入でpropsから必要な値を取り出す
export const SearchForm: React.FC<SearchFormProps> = ({
  onSearch,
  isLoading,
}) => {
  // useState：各フォームフィールドの状態を管理
  // [状態変数, 更新関数] = useState(初期値)

  // 都道府県コードの状態（初期値：空文字列）
  const [serviceArea, setServiceArea] = useState("");
  // 住所キーワードの状態（初期値：空文字列）
  const [address, setAddress] = useState("");
  // 選択されたジャンルコードの状態（初期値：空文字列）
  const [selectedGenre, setSelectedGenre] = useState("");
  // フリーワード検索の状態（初期値：空文字列）
  const [keyword, setKeyword] = useState("");

  // useGenresフックでジャンル一覧を取得
  // 分割代入で必要なデータを取り出す
  const { genres, isLoading: genresLoading } = useGenres();

  // フォーム送信時のハンドラー関数
  const handleSubmit = (e: React.FormEvent) => {
    // e.preventDefault()でフォームのデフォルト動作（ページリロード）を防ぐ
    e.preventDefault();

    // 必須項目のバリデーション: 都道府県または緯度経度のいずれかを指定する
    if (!serviceArea && (lat === null || lng === null)) {
      alert(
        "都道府県か現在地のいずれかを指定してください（現在地ボタンで取得できます）",
      );
      return;
    }

    // onSearch関数を呼び出して検索を実行
    // 値が空の場合はundefinedを設定（||演算子でフォールバック）
    onSearch({
      serviceArea: serviceArea || undefined, // 都道府県コード（未指定可）
      address: address || undefined, // オプション：住所
      genre: selectedGenre || undefined, // オプション：ジャンル
      keyword: keyword || undefined, // オプション：キーワード
      lat: lat ?? undefined,
      lng: lng ?? undefined,
      range: range ?? undefined,
    });
  };

  // 現在地（Geolocation API）関連の状態
  const [lat, setLat] = useState<number | null>(null);
  const [lng, setLng] = useState<number | null>(null);
  const [range, setRange] = useState<number>(3); // デフォルトは3（約1000m）

  const handleUseCurrentLocation = () => {
    if (!navigator.geolocation) {
      alert("Geolocation API がサポートされていません");
      return;
    }
    navigator.geolocation.getCurrentPosition(
      (position) => {
        setLat(position.coords.latitude);
        setLng(position.coords.longitude);
      },
      (err) => {
        console.error(err);
        alert("現在地の取得に失敗しました");
      },
      { enableHighAccuracy: true },
    );
  };

  const clearLocation = () => {
    setLat(null);
    setLng(null);
  };

  // JSXを返す：フォームのHTML構造
  return (
    // form要素：onSubmitでフォーム送信を処理
    // classNameでTailwind CSSのクラスを適用（スタイリング）
    <form
      onSubmit={handleSubmit}
      className="bg-white shadow-md rounded-lg p-6 mb-6"
    >
      {/* フォームのタイトル */}
      <h2 className="text-2xl font-bold mb-4 text-gray-800">レストラン検索</h2>

      {/* グリッドレイアウト：モバイルで1列、タブレット以上で2列 */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        {/* 都道府県選択フィールド */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            都道府県 <span className="text-red-500">*</span> {/* 必須マーク */}
          </label>
          {/* select要素：プルダウンメニュー */}
          <select
            value={serviceArea} // 現在の値
            // onChange: 選択が変更された時に状態を更新
            // e.target.valueで選択された値を取得
            onChange={(e) => setServiceArea(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required // HTML5のバリデーション：必須項目
          >
            {/* デフォルトオプション */}
            <option value="">選択してください</option>
            {/* SERVICE_AREASオブジェクトをループして地域とその都道府県を表示 */}
            {/* Object.entries()でオブジェクトを[key, value]の配列に変換 */}
            {Object.entries(SERVICE_AREAS).map(([region, prefs]) => (
              // optgroup: 選択肢をグループ化
              <optgroup key={region} label={region}>
                {/* 各地域内の都道府県をループ */}
                {Object.entries(prefs).map(([pref, code]) => (
                  // option要素：選択肢
                  <option key={code} value={code}>
                    {pref}
                  </option>
                ))}
              </optgroup>
            ))}
          </select>
        </div>

        {/* 現在地取得と検索範囲 */}
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
            {lat && lng ? (
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
              <option value={1}>1: 約300m</option>
              <option value={2}>2: 約500m</option>
              <option value={3}>3: 約1000m</option>
              <option value={4}>4: 約2000m</option>
              <option value={5}>5: 約3000m</option>
            </select>
          </div>
        </div>

        {/* 住所入力フィールド */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            住所（部分一致）
          </label>
          {/* input要素：テキスト入力 */}
          <input
            type="text"
            value={address}
            onChange={(e) => setAddress(e.target.value)}
            placeholder="例: 新宿" // プレースホルダーテキスト
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* ジャンル選択フィールド */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            ジャンル
          </label>
          <select
            value={selectedGenre}
            onChange={(e) => setSelectedGenre(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={genresLoading} // ジャンルロード中は無効化
          >
            <option value="">すべて</option>
            {/* genresをmapでループして選択肢を生成 */}
            {genres.map((genre) => (
              <option key={genre.code} value={genre.code}>
                {genre.name}
              </option>
            ))}
          </select>
        </div>

        {/* キーワード入力フィールド */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            キーワード
          </label>
          <input
            type="text"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder="例: 個室"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      {/* 検索ボタン */}
      <button
        type="submit" // type="submit"でフォーム送信トリガー
        disabled={isLoading} // ローディング中は無効化
        className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
      >
        {/* 三項演算子でボタンテキストを切り替え */}
        {isLoading ? "検索中..." : "検索"}
      </button>
    </form>
  );
};
