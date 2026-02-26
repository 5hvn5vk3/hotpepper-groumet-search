// 各種コンポーネントをインポート
// @/はプロジェクトのsrcディレクトリへのエイリアス（省略記法）
import {
  ErrorMessage, // エラーメッセージ表示コンポーネント
  LoadingSpinner, // ローディング表示コンポーネント
  Pagination, // ページネーションコンポーネント
  RestaurantDetail, // レストラン詳細モーダルコンポーネント
  RestaurantList, // レストラン一覧表示コンポーネント
  SearchForm, // 検索フォームコンポーネント
} from "@/components"; // コンポーネントのインデックスから一括インポート
// TypeScript型定義をインポート
import type { RestaurantDetailStatus, SearchParams, Shop } from "@/types";
// カスタムフックをインポート
import { useRestaurantSearch } from "@/hooks/useRestaurantSearch"; // レストラン検索ロジックを管理
import { useModal } from "@/hooks/useModal"; // モーダルの開閉を管理
import { fetchRestaurantDetail } from "@/api/restaurantApi";
import { useMemo, useRef, useState } from "react"; // 現在地状態を管理
// 定数をインポート
import { ITEMS_PER_PAGE } from "@/constants"; // 1ページあたりの表示件数

// Appコンポーネント：アプリケーション全体のメインコンポーネント
function App() {
  // useRestaurantSearchフック：レストラン検索の状態とロジックを管理
  // 分割代入で必要な値と関数を取り出す
  const {
    searchResult, // 検索結果データ
    isLoading, // ローディング中かどうかのフラグ
    hasSearched, // 検索が実行されたかどうかのフラグ
    error, // エラーメッセージ
    currentPage, // 現在のページ番号
    search, // 検索を実行する関数
    changePage, // ページを変更する関数
    clearError, // エラーをクリアする関数
  } = useRestaurantSearch(ITEMS_PER_PAGE);

  // useModalフック：モーダル（詳細画面）の開閉状態を管理
  // Shop型のデータを扱うモーダル
  const { isOpen, data: selectedRestaurant, open, close } = useModal<Shop>();

  // ユーザーの現在地（緯度・経度）の状態
  const [userLat, setUserLat] = useState<number | null>(null);
  const [userLng, setUserLng] = useState<number | null>(null);
  const [detailRestaurant, setDetailRestaurant] = useState<Shop | null>(null);
  const [detailStatus, setDetailStatus] =
    useState<RestaurantDetailStatus>("idle");
  const detailRequestIdRef = useRef(0);

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
      // credit_cardは初回のlite+credit_cardを保持する
      credit_card: selectedRestaurant.credit_card,
    };
  }, [selectedRestaurant, detailRestaurant]);

  // SearchFormから現在地が取得されたときのハンドラー
  const handleLocationChange = (lat: number, lng: number) => {
    setUserLat(lat);
    setUserLng(lng);
  };

  // 検索フォームから検索が実行されたときのハンドラー関数
  const handleSearch = (params: SearchParams) => {
    // searchフックの関数を呼び出して検索を実行
    search(params);
  };

  // ページ変更ボタンがクリックされたときのハンドラー関数
  const handlePageChange = (page: number) => {
    // changePageフックの関数を呼び出してページを変更
    changePage(page);
  };

  // レストランカードがクリックされたときのハンドラー関数
  const handleSelectRestaurant = async (restaurant: Shop) => {
    // モーダルを開き、選択されたレストランのデータを渡す
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
    } catch (detailError) {
      if (detailRequestIdRef.current === requestId) {
        setDetailStatus("failed");
      }
      console.error("Failed to fetch restaurant detail:", detailError);
    }
  };

  const handleCloseRestaurantDetail = () => {
    // 進行中のレスポンスを無効化
    detailRequestIdRef.current += 1;
    setDetailRestaurant(null);
    setDetailStatus("idle");
    close();
  };

  // JSXを返す：画面に表示するHTML風の構造
  return (
    // 最小高さを画面全体に、背景色をグレーに設定
    <div className="min-h-screen bg-gray-100">
      {/* ヘッダー部分：青い背景で影付き */}
      <header className="fixed top-0 left-0 right-0 z-40 bg-red-600 text-white py-6 shadow-lg">
        {/* コンテナ：最大幅を設定し、中央揃え */}
        <div className="container mx-auto px-4">
          {/* タイトル：大きな太字で表示 */}
          <h1 className="text-3xl font-bold">ホットペッパー レストラン検索</h1>
          <p>Powered by ホットペッパーグルメ Webサービス</p>
        </div>
      </header>

      {/* メインコンテンツエリア */}
      <main className="container mx-auto px-4 pt-36 pb-8">
        {/* 検索フォームコンポーネント
            onSearch: 検索実行時のコールバック関数
            isLoading: ローディング中は検索ボタンを無効化 */}
        <SearchForm
          onSearch={handleSearch}
          isLoading={isLoading}
          hasSearched={hasSearched}
          onLocationChange={handleLocationChange}
        />

        {/* 条件付きレンダリング：エラーがある場合のみ表示
            &&演算子は左側がtrueの場合に右側を評価・レンダリング */}
        {error && <ErrorMessage message={error} onClose={clearError} />}

        {/* ローディング中の場合、スピナーを表示 */}
        {isLoading && <LoadingSpinner />}

        {/* ローディング中でなく、かつ検索が実行済みの場合、結果を表示
            !isLoadingはisLoadingがfalseであることを意味 */}
        {!isLoading && hasSearched && (
          <RestaurantList
            // searchResult?.shop: オプショナルチェイニング（searchResultがnullの場合エラーにならない）
            // ?? []: Null合体演算子（左側がnullまたはundefinedの場合、右側の空配列を使用）
            restaurants={searchResult?.shop ?? []}
            onSelectRestaurant={handleSelectRestaurant}
            userLat={userLat ?? undefined}
            userLng={userLng ?? undefined}
          />
        )}

        {/* 検索結果があり、かつ結果件数が1件以上の場合、ページネーションを表示 */}
        {searchResult && searchResult.results_available > 0 && (
          <Pagination
            currentPage={currentPage} // 現在のページ番号
            totalCount={searchResult.results_available} // 総件数
            count={ITEMS_PER_PAGE} // 1ページあたりの件数
            onPageChange={handlePageChange} // ページ変更時の処理
            disabled={isLoading} // ローディング中は無効化
            start={searchResult.results_start} // 表示開始位置
            available={searchResult.results_available} // 総利用可能件数
          />
        )}
      </main>

      {/* モーダルが開いており、かつレストランが選択されている場合、詳細画面を表示 */}
      {isOpen && modalRestaurant && (
        <RestaurantDetail
          restaurant={modalRestaurant}
          detailStatus={detailStatus}
          onClose={handleCloseRestaurantDetail}
          userLat={userLat ?? undefined}
          userLng={userLng ?? undefined}
        />
      )}
    </div>
  );
}

// AppコンポーネントをエクスポートしてWの他のファイルから使用可能にする
// defaultエクスポートは1ファイル1つのみ
export default App;
