// 型定義をインポート
import type { GourmetSearchResponse, GourmetSearchParams, Shop } from "@/types";
// API関数をインポート
import { apiGet } from "@/api/client";

// searchRestaurants関数：検索パラメータを元にレストランを検索
// params: 検索条件（位置情報、住所、ジャンル、キーワード、ページ番号など）
// 戻り値: 検索結果のレスポンスデータ
export const searchRestaurants = async (
  params: GourmetSearchParams,
): Promise<GourmetSearchResponse> => {
  // 位置情報（lat/lng）または住所またはキーワードのいずれかがあれば検索を行う
  const hasLatLng =
    typeof params.lat === "number" && typeof params.lng === "number";
  const address = params.address?.trim() ?? "";
  const keyword = params.keyword?.trim() ?? "";
  const hasAddress = address !== "";
  const hasKeyword = keyword !== "";

  if (!hasLatLng && !hasAddress && !hasKeyword) {
    // 検索条件が未指定の場合は空の結果を返す
    return {
      results: {
        api_version: "1.0",
        results_available: 0,
        results_returned: "0",
        results_start: 1,
        shop: [],
      },
    };
  }

  // ページ番号のバリデーション：最小値を1に制限
  // Math.max(1, x)は1とxのうち大きい方を返す
  // params.page ?? 1は、params.pageがnullまたはundefinedの場合に1を使用
  const page = Math.max(1, params.page ?? 1);

  // 1ページあたりの件数のバリデーション：1〜100の範囲に制限
  // Math.min(max, x)はmaxとxのうち小さい方を返す
  const count = Math.min(Math.max(params.count ?? 20, 1), 100);

  // API用の開始位置を計算
  // 例：2ページ目で20件/ページの場合 → (2-1) * 20 + 1 = 21
  const start = (page - 1) * count + 1;

  // URLSearchParamsオブジェクトでクエリパラメータを構築
  // URLSearchParamsはクエリ文字列を簡単に作成できるWeb標準API
  const queryParams = new URLSearchParams({
    start: start.toString(), // 開始位置を文字列に変換
    count: count.toString(), // 取得件数を文字列に変換
  });

  // オプショナルパラメータの追加（値がある場合のみ）
  if (hasAddress) {
    queryParams.set("address", address); // 住所キーワード
  }
  if (params.genre) {
    queryParams.set("genre", params.genre); // ジャンルコード
  }
  if (hasKeyword) {
    queryParams.set("keyword", keyword); // フリーワード
  }
  // 緯度経度が指定されていればクエリに追加
  if (hasLatLng) {
    queryParams.set("lat", params.lat!.toString());
    queryParams.set("lng", params.lng!.toString());
    if (params.range) {
      queryParams.set("range", params.range.toString());
    }
  }

  // apiGet関数を使用してAPIリクエストを送信
  // テンプレートリテラルでURLを構築：/api/gourmet?パラメータ
  // GourmetSearchResponse型でレスポンスの型を指定
  return apiGet<GourmetSearchResponse>(`/api/gourmet?${queryParams}`);
};

// fetchRestaurantDetail関数：店舗IDから詳細情報を取得（type指定なし）
export const fetchRestaurantDetail = async (id: string): Promise<Shop | null> => {
  const shopID = id.trim();
  if (shopID === "") {
    return null;
  }

  const queryParams = new URLSearchParams({ id: shopID });
  const response = await apiGet<GourmetSearchResponse>(`/api/gourmet?${queryParams}`);
  return response.results.shop[0] ?? null;
};
